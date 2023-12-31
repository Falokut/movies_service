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

type moviesPreviewCache struct {
	rdb    *redis.Client
	logger *logrus.Logger
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

func NewMoviesPreviewCache(logger *logrus.Logger, opt *redis.Options) (*moviesPreviewCache, error) {
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

	return &moviesPreviewCache{rdb: rdb, logger: logger}, nil
}

func (c *moviesPreviewCache) GetMovie(ctx context.Context, movieId int32) (MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesPreviewCache.GetMovie")
	defer span.Finish()
	res, err := c.rdb.Get(ctx, fmt.Sprint(movieId)).Bytes()
	if err != nil {
		return MoviePreview{}, err
	}

	var movie cachedMoviePreview
	if err := json.Unmarshal(res, &movie); err != nil {
		return MoviePreview{}, err
	}
	movie.ID = movieId
	return convertCacheMoviePreviewToMoviePreview(movie), nil
}

type cacheMoviePreviewFilterKey struct {
	MoviesIDs    string `json:"1,omitempty"`
	GenresIDs    string `json:"2,omitempty"`
	DiretorsIDs  string `json:"3,omitempty"`
	CountriesIDs string `json:"4,omitempty"`
	Title        string `json:"5,omitempty"`
	Limit        uint32 `json:"6,omitempty"`
	Offset       uint32 `json:"7,omitempty"`
	AgeRating    string `json:"8,omitempty"`
}

func buildPreviewFilterKey(filter MoviesFilter, limit, offset uint32) string {
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

type cachedPreviewFilteredRequest struct {
	Ids []string `json:"-,"`
}

func (c *moviesPreviewCache) GetMovies(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesPreviewCache.GetMovies")
	defer span.Finish()

	key := buildPreviewFilterKey(filter, limit, offset)
	res, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return []string{}, err
	}

	var cache cachedPreviewFilteredRequest
	if err = json.Unmarshal(res, &cache); err != nil {
		return []string{}, err
	}

	return cache.Ids, nil
}

func (c *moviesPreviewCache) CacheMovies(ctx context.Context, movies []MoviePreview, ttl time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesPreviewCache.CacheMovies")
	defer span.Finish()

	tx := c.rdb.Pipeline()
	for _, movie := range movies {
		movieToCache := convertMoviePreviewToCacheMovie(movie)
		toCache, err := json.Marshal(movieToCache)
		if err != nil {
			return err
		}
		tx.Set(ctx, fmt.Sprint(movie.ID), toCache, ttl)
	}
	_, err := tx.Exec(ctx)

	return err
}

func (c *moviesPreviewCache) CacheFilteredRequest(ctx context.Context, filter MoviesFilter,
	limit, offset uint32, moviesIDs []string, ttl time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesPreviewCache.CacheFilteredRequest")
	defer span.Finish()

	key := buildPreviewFilterKey(filter, limit, offset)

	var toCache = cachedPreviewFilteredRequest{
		Ids: moviesIDs,
	}

	marshalled, err := json.Marshal(&toCache)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, marshalled, ttl).Err()
}

type cachedMoviePreview struct {
	ID          int32    `json:"-"`
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

func convertMoviePreviewToCacheMovie(movie MoviePreview) cachedMoviePreview {
	return cachedMoviePreview{
		ID:          movie.ID,
		TitleRU:     movie.TitleRU,
		TitleEN:     movie.TitleEN.String,
		Description: movie.ShortDescription,
		Genres:      movie.Genres,
		Duration:    movie.Duration,
		PosterID:    movie.PreviewPosterID.String,
		Countries:   movie.Countries,
		ReleaseYear: movie.ReleaseYear,
		AgeRating:   movie.AgeRating,
	}
}

func convertCacheMoviePreviewToMoviePreview(movie cachedMoviePreview) MoviePreview {
	return MoviePreview{
		ID:               movie.ID,
		TitleRU:          movie.TitleRU,
		TitleEN:          sql.NullString{String: movie.TitleEN, Valid: true},
		ShortDescription: movie.Description,
		Genres:           movie.Genres,
		Duration:         movie.Duration,
		PreviewPosterID:  sql.NullString{String: movie.PosterID, Valid: true},
		Countries:        movie.Countries,
		ReleaseYear:      movie.ReleaseYear,
		AgeRating:        movie.AgeRating,
	}
}
