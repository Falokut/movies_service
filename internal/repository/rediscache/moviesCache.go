package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/movies_service/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type moviesCache struct {
	rdb     *redis.Client
	logger  *logrus.Logger
	metrics Metrics
}

func (c *moviesCache) PingContext(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging movies cache: %w", err)
	}

	return nil
}

func (c *moviesCache) Shutdown() {
	c.rdb.Close()
}

func NewMoviesCache(logger *logrus.Logger, opt *redis.Options, metrics Metrics) (*moviesCache, error) {
	logger.Info("Creating movies cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Info("Pinging movies cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	return &moviesCache{rdb: rdb, logger: logger, metrics: metrics}, nil
}

func (c *moviesCache) GetMovie(ctx context.Context, movieId int32) (movie models.RepositoryMovie, err error) {
	defer c.updateMetrics(err, "GetMovie")
	defer handleError(ctx, &err)
	defer c.logError(err, "GetMovie")

	res, err := c.rdb.Get(ctx, fmt.Sprint(movieId)).Bytes()
	if err != nil {
		return
	}

	var cached cachedMovie
	if err = json.Unmarshal(res, &cached); err != nil {
		return
	}

	movie.ID = movieId
	return convertCacheMovieToMovie(cached), nil
}

type cacheMovieFilterKey struct {
	MoviesIDs    string `json:"1,omitempty"`
	GenresIDs    string `json:"2,omitempty"`
	CountriesIDs string `json:"4,omitempty"`
	Title        string `json:"5,omitempty"`
	Limit        uint32 `json:"6,omitempty"`
	Offset       uint32 `json:"7,omitempty"`
}

func buildFilterKey(filter models.MoviesFilter, limit, offset uint32) string {
	key := cacheMovieFilterKey{
		MoviesIDs:    filter.MoviesIDs,
		GenresIDs:    filter.GenresIDs,
		CountriesIDs: filter.CountriesIDs,
		Title:        filter.Title,
		Limit:        limit,
		Offset:       offset,
	}

	keyBytes, _ := json.Marshal(key)
	return string(keyBytes)
}

func (c *moviesCache) CacheMovies(ctx context.Context, movies []models.RepositoryMovie, ttl time.Duration) (err error) {
	defer handleError(ctx, &err)
	defer c.logError(err, "CacheMovies")

	tx := c.rdb.Pipeline()
	for _, movie := range movies {
		movieToCache := convertMovieToCacheMovie(movie)
		toCache, merr := json.Marshal(movieToCache)
		if merr != nil {
			err = merr
			return
		}

		tx.Set(ctx, fmt.Sprint(movie.ID), toCache, ttl)
	}
	_, err = tx.Exec(ctx)
	return
}

type cachedMovie struct {
	ID                  int32    `json:"-"`
	TitleRU             string   `json:"title_ru"`
	TitleEN             string   `json:"title_en"`
	Description         string   `json:"description"`
	Genres              []string `json:"genres"`
	Duration            int32    `json:"duration"`
	PosterID            string   `json:"poster_picture_id"`
	BackgroundPictureID string   `json:"background_picture_id"`
	Countries           []string `json:"countries"`
	ReleaseYear         int32    `json:"release_year"`
	AgeRating           string   `json:"age_rating"`
}

func convertMovieToCacheMovie(movie models.RepositoryMovie) cachedMovie {
	return cachedMovie{
		ID:                  movie.ID,
		TitleRU:             movie.TitleRU,
		TitleEN:             movie.TitleEN,
		Description:         movie.Description,
		Genres:              movie.Genres,
		Duration:            movie.Duration,
		PosterID:            movie.PosterID,
		BackgroundPictureID: movie.BackgroundPictureID,
		Countries:           movie.Countries,
		ReleaseYear:         movie.ReleaseYear,
		AgeRating:           movie.AgeRating,
	}
}

func convertCacheMovieToMovie(movie cachedMovie) models.RepositoryMovie {
	return models.RepositoryMovie{
		ID:                  movie.ID,
		TitleRU:             movie.TitleRU,
		TitleEN:             movie.TitleEN,
		Description:         movie.Description,
		Genres:              movie.Genres,
		Duration:            movie.Duration,
		PosterID:            movie.PosterID,
		BackgroundPictureID: movie.BackgroundPictureID,
		Countries:           movie.Countries,
		ReleaseYear:         movie.ReleaseYear,
		AgeRating:           movie.AgeRating,
	}
}

func (c *moviesCache) logError(err error, functionName string) {
	if err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if errors.As(err, &repoErr) {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           repoErr.Msg,
				"error.code":          repoErr.Code,
			},
		).Error("movies cache error occurred")
	} else {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("movies cache error occurred")
	}
}

func (c *moviesCache) updateMetrics(err error, functionName string) {
	if err == nil {
		c.metrics.IncCacheHits(functionName, 1)
		return
	}
	if models.Code(err) == models.NotFound {
		c.metrics.IncCacheMiss(functionName, 1)
	}
}
