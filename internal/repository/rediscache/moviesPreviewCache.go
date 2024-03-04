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
	"golang.org/x/exp/maps"
)

type moviesPreviewCache struct {
	rdb     *redis.Client
	logger  *logrus.Logger
	metrics Metrics
}

func (c *moviesPreviewCache) PingContext(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging movies cache: %w", err)
	}

	return nil
}

func (c *moviesPreviewCache) Shutdown() {
	c.rdb.Close()
}

func NewMoviesPreviewCache(logger *logrus.Logger, opt *redis.Options, metrics Metrics) (*moviesPreviewCache, error) {
	logger.Info("Creating preview movies cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Info("Pinging movies cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	return &moviesPreviewCache{rdb: rdb, logger: logger, metrics: metrics}, nil
}

func (c *moviesPreviewCache) GetMovie(ctx context.Context, movieId int32) (preview models.RepositoryMoviePreview, err error) {
	defer c.updateMetrics(err, "GetMovie")
	defer handleError(ctx, &err)
	defer c.logError(err, "GetMovie")

	res, err := c.rdb.Get(ctx, fmt.Sprint(movieId)).Bytes()
	if err != nil {
		return
	}

	var movie cachedMoviePreview
	if err = json.Unmarshal(res, &movie); err != nil {
		return
	}

	movie.ID = movieId
	return convertCacheMoviePreviewToMoviePreview(movie), nil
}

func (c *moviesPreviewCache) GetMovies(ctx context.Context, ids []string) (movies []models.RepositoryMoviePreview,
	notFoundIds []string, err error) {
	defer handleError(ctx, &err)
	defer c.logError(err, "GetMovies")

	var moviesIDs = make(map[string]struct{}, len(ids))
	for _, id := range ids {
		moviesIDs[id] = struct{}{}
	}

	cached, err := c.rdb.MGet(ctx, ids...).Result()
	if err != nil {
		return
	}

	movies = make([]models.RepositoryMoviePreview, 0, len(cached))
	for _, cache := range cached {
		if cache == nil {
			continue
		}

		movie := cachedMoviePreview{}
		err = json.Unmarshal([]byte(cache.(string)), &movie)
		if err != nil {
			return
		}
		delete(moviesIDs, fmt.Sprint(movie.ID))
		movies = append(movies, convertCacheMoviePreviewToMoviePreview(movie))
	}

	notFoundIds = maps.Keys(moviesIDs)
	c.metrics.IncCacheHits("GetMovies", int32(len(ids)-len(notFoundIds)))
	c.metrics.IncCacheMiss("GetMovies", int32(len(notFoundIds)))

	return
}

type cacheMoviePreviewFilterKey struct {
	MoviesIDs    string `json:"1,omitempty"`
	GenresIDs    string `json:"2,omitempty"`
	CountriesIDs string `json:"3,omitempty"`
	Title        string `json:"4,omitempty"`
	Limit        uint32 `json:"5,omitempty"`
	Offset       uint32 `json:"6,omitempty"`
	AgeRating    string `json:"7,omitempty"`
}

func buildPreviewFilterKey(filter models.MoviesFilter, limit, offset uint32) string {
	key := cacheMoviePreviewFilterKey{
		MoviesIDs:    filter.MoviesIDs,
		GenresIDs:    filter.GenresIDs,
		CountriesIDs: filter.CountriesIDs,
		Title:        filter.Title,
		Limit:        limit,
		Offset:       offset,
		AgeRating:    filter.AgeRating,
	}

	keyBytes, _ := json.Marshal(key)
	return string(keyBytes)
}

func (c *moviesPreviewCache) GetMoviesIDs(ctx context.Context, filter models.MoviesFilter, limit, offset uint32) (ids []string, err error) {
	defer c.updateMetrics(err, "GetMoviesIDs")
	defer handleError(ctx, &err)
	defer c.logError(err, "GetMoviesIDs")

	key := buildPreviewFilterKey(filter, limit, offset)
	ids, err = c.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return
	}

	if len(ids) == 0 {
		err = models.Error(models.NotFound, "")
	}
	return
}

func (c *moviesPreviewCache) CacheMovies(ctx context.Context, movies []models.RepositoryMoviePreview, ttl time.Duration) (err error) {
	defer handleError(ctx, &err)
	defer c.logError(err, "CacheMovies")

	tx := c.rdb.Pipeline()
	for _, movie := range movies {
		movieToCache := convertMoviePreviewToCacheMovie(movie)
		toCache, merr := json.Marshal(movieToCache)
		if err != nil {
			err = merr
			return
		}
		tx.Set(ctx, fmt.Sprint(movie.ID), toCache, ttl)
	}
	_, err = tx.Exec(ctx)
	return
}

func (c *moviesPreviewCache) CacheFilteredRequest(ctx context.Context, filter models.MoviesFilter,
	limit, offset uint32, moviesIDs []string, ttl time.Duration) (err error) {
	defer handleError(ctx, &err)
	defer c.logError(err, "CacheFilteredRequest")

	key := buildPreviewFilterKey(filter, limit, offset)

	tx := c.rdb.Pipeline()

	err = tx.SAdd(ctx, key, moviesIDs).Err()
	if err != nil {
		return
	}
	err = tx.Expire(ctx, key, ttl).Err()
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx)
	return
}

type cachedMoviePreview struct {
	ID          int32    `json:"id"`
	TitleRU     string   `json:"title_ru"`
	TitleEN     string   `json:"title_en"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	Duration    int32    `json:"duration"`
	PosterID    string   `json:"poster_picture_id"`
	Countries   []string `json:"countries"`
	ReleaseYear int32    `json:"release_year"`
	AgeRating   string   `json:"age_rating"`
}

func convertMoviePreviewToCacheMovie(movie models.RepositoryMoviePreview) cachedMoviePreview {
	return cachedMoviePreview{
		ID:          movie.ID,
		TitleRU:     movie.TitleRU,
		TitleEN:     movie.TitleEN,
		Description: movie.ShortDescription,
		Genres:      movie.Genres,
		Duration:    movie.Duration,
		PosterID:    movie.PreviewPosterID,
		Countries:   movie.Countries,
		ReleaseYear: movie.ReleaseYear,
		AgeRating:   movie.AgeRating,
	}
}

func convertCacheMoviePreviewToMoviePreview(movie cachedMoviePreview) models.RepositoryMoviePreview {
	return models.RepositoryMoviePreview{
		ID:               movie.ID,
		TitleRU:          movie.TitleRU,
		TitleEN:          movie.TitleEN,
		ShortDescription: movie.Description,
		Genres:           movie.Genres,
		Duration:         movie.Duration,
		PreviewPosterID:  movie.PosterID,
		Countries:        movie.Countries,
		ReleaseYear:      movie.ReleaseYear,
		AgeRating:        movie.AgeRating,
	}
}

func (c *moviesPreviewCache) logError(err error, functionName string) {
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
		).Error("movies preview cache error occurred")
	} else {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("movies preview cache error occurred")
	}
}

func (c *moviesPreviewCache) updateMetrics(err error, functionName string) {
	if err == nil {
		c.metrics.IncCacheHits(functionName, 1)
		return
	}
	if models.Code(err) == models.NotFound {
		c.metrics.IncCacheMiss(functionName, 1)
	}
}
