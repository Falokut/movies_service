package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type moviesCache struct {
	rdb    *redis.Client
	logger *logrus.Logger
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

func NewMoviesCache(logger *logrus.Logger, opt *redis.Options) (*moviesCache, error) {
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

	return &moviesCache{rdb: rdb, logger: logger}, nil
}

func (c *moviesCache) GetMovie(ctx context.Context, movieId string) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesCache.GetMovie")
	defer span.Finish()
	res, err := c.rdb.Get(ctx, fmt.Sprint(movieId)).Bytes()
	if err != nil {
		return Movie{}, err
	}

	var movie cachedMovie
	if err := json.Unmarshal(res, &movie); err != nil {
		return Movie{}, err
	}
	movie.ID = movieId
	return convertCacheMovieToMovie(movie), nil
}

type cacheMovieFilterKey struct {
	MoviesIDs    string `json:"1,omitempty"`
	GenresIDs    string `json:"2,omitempty"`
	DiretorsIDs  string `json:"3,omitempty"`
	CountriesIDs string `json:"4,omitempty"`
	Title        string `json:"5,omitempty"`
	Limit        uint32 `json:"6,omitempty"`
	Offset       uint32 `json:"7,omitempty"`
}

func buildFilterKey(filter MoviesFilter, limit, offset uint32) string {
	key := cacheMovieFilterKey{
		MoviesIDs:    filter.MoviesIDs,
		GenresIDs:    filter.GenresIDs,
		DiretorsIDs:  filter.DirectorsIDs,
		CountriesIDs: filter.CountriesIDs,
		Title:        filter.Title,
		Limit:        limit,
		Offset:       offset,
	}

	keyBytes, _ := json.Marshal(key)
	return string(keyBytes)
}

type cachedFilteredRequest struct {
	Ids []string `json:"-,"`
}

func (c *moviesCache) GetMovies(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesCache.GetMovies")
	defer span.Finish()

	key := buildFilterKey(filter, limit, offset)
	res, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return []string{}, err
	}

	var cache cachedFilteredRequest
	if err = json.Unmarshal(res, &cache); err != nil {
		return []string{}, err
	}

	return cache.Ids, nil
}

func (c *moviesCache) CacheMovies(ctx context.Context, movies []Movie, ttl time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesCache.CacheMovies")
	defer span.Finish()

	tx := c.rdb.Pipeline()
	for _, movie := range movies {
		movieToCache := convertMovieToCacheMovie(movie)
		toCache, err := json.Marshal(movieToCache)
		if err != nil {
			return err
		}
		tx.Set(ctx, fmt.Sprint(movie.ID), toCache, ttl)
	}
	_, err := tx.Exec(ctx)
	return err
}

func (c *moviesCache) CacheFilteredRequest(ctx context.Context, filter MoviesFilter,
	limit, offset uint32, moviesIDs []string, ttl time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesCache.CacheFilteredRequest")
	defer span.Finish()

	key := buildFilterKey(filter, limit, offset)

	var toCache = cachedFilteredRequest{
		Ids: moviesIDs,
	}

	marshalled, err := json.Marshal(&toCache)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, marshalled, ttl).Err()
}

type cachedMovie struct {
	ID                  string `json:"-"`
	TitleRU             string `json:"title_ru"`
	TitleEN             string `json:"title_en"`
	Description         string `json:"description"`
	Genres              string `json:"genres"`
	Duration            int32  `json:"duration"`
	PosterID            string `json:"poster_picture_id"`
	BackgroundPictureID string `json:"background_picture_id"`
	DiretorsIDs         string `json:"directors"`
	CountriesIDs        string `json:"countries"`
	ReleaseYear         int32  `json:"release_year"`
	AgeRating           string `json:"age_rating"`
}

func convertMovieToCacheMovie(movie Movie) cachedMovie {
	return cachedMovie{
		ID:                  movie.ID,
		TitleRU:             movie.TitleRU,
		TitleEN:             movie.TitleEN.String,
		Description:         movie.Description,
		Genres:              movie.Genres.String,
		Duration:            movie.Duration,
		PosterID:            movie.PosterID.String,
		BackgroundPictureID: movie.BackgroundPictureID.String,
		DiretorsIDs:         movie.DirectorsIDs.String,
		CountriesIDs:        movie.CountriesIDs.String,
		ReleaseYear:         movie.ReleaseYear,
		AgeRating:           movie.AgeRating,
	}
}

func convertCacheMovieToMovie(movie cachedMovie) Movie {
	return Movie{
		ID:           movie.ID,
		TitleRU:      movie.TitleRU,
		TitleEN:      sql.NullString{String: movie.TitleEN, Valid: true},
		Description:  movie.Description,
		Genres:       sql.NullString{String: movie.Genres, Valid: true},
		Duration:     movie.Duration,
		PosterID:     sql.NullString{String: movie.PosterID, Valid: true},
		DirectorsIDs: sql.NullString{String: movie.DiretorsIDs, Valid: true},
		CountriesIDs: sql.NullString{String: movie.CountriesIDs, Valid: true},
		ReleaseYear:  movie.ReleaseYear,
		AgeRating:    movie.AgeRating,
	}
}
