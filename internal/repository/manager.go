package repository

import (
	"context"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RepositoryManagerConfig struct {
	MovieTTL        time.Duration
	FilteredTTL     time.Duration
	MoviePreviewTTL time.Duration
}

type RepositoryManager struct {
	moviesRepo         MoviesRepository
	moviesCache        MoviesCache
	moviesPreviewRepo  MoviesPreviewRepository
	moviesPreviewCache MoviesPreviewCache
	ageRatingsRepo     AgeRatingRepository
	cfg                RepositoryManagerConfig
	logger             *logrus.Logger
}

func NewMoviesRepositoryManager(moviesRepo MoviesRepository, moviesCache MoviesCache, moviesPreviewRepo MoviesPreviewRepository,
	moviesPreviewCache MoviesPreviewCache, ageRatingsRepo AgeRatingRepository,
	cfg RepositoryManagerConfig, logger *logrus.Logger) *RepositoryManager {
	return &RepositoryManager{
		moviesRepo:         moviesRepo,
		moviesCache:        moviesCache,
		moviesPreviewCache: moviesPreviewCache,
		moviesPreviewRepo:  moviesPreviewRepo,
		ageRatingsRepo:     ageRatingsRepo,
		cfg:                cfg,
		logger:             logger,
	}
}

func (m *RepositoryManager) GetMovie(ctx context.Context, movieID string) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("has_error", err != nil)

	movie, err := m.moviesCache.GetMovie(ctx, movieID)
	if err == nil {
		return movie, nil
	}
	if !errors.Is(err, redis.Nil) {
		m.logger.Warn(err)
	}

	movie, err = m.moviesRepo.GetMovie(ctx, movieID)
	if err != nil {
		return Movie{}, err
	}

	go func(Movie) {
		if err := m.moviesCache.CacheMovies(context.Background(), []Movie{movie}, m.cfg.MovieTTL); err != nil {
			m.logger.Error(err)
		}
	}(movie)

	return movie, nil
}

func (m *RepositoryManager) GetMoviePreview(ctx context.Context, movieID string) (MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMoviePreview")
	defer span.Finish()
	var err error
	defer span.SetTag("has_error", err != nil)

	movie, err := m.moviesPreviewCache.GetMovie(ctx, movieID)
	if err == nil {
		return movie, nil
	}
	if !errors.Is(err, redis.Nil) {
		m.logger.Warn(err)
	}

	movie, err = m.moviesPreviewRepo.GetMoviePreview(ctx, movieID)
	if err != nil {
		return MoviePreview{}, err
	}

	go func(MoviePreview) {
		if err := m.moviesPreviewCache.CacheMovies(context.Background(), []MoviePreview{movie}, m.cfg.MoviePreviewTTL); err != nil {
			m.logger.Error(err)
		}
	}(movie)

	return movie, nil
}

func (m *RepositoryManager) GetMoviesPreview(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("has_error", err != nil)

	m.logger.Info("Trying get movies ids from cache")
	moviesIds, err := m.moviesPreviewCache.GetMovies(ctx, filter, limit, offset)
	inCache := true
	if err != nil {
		m.logger.Warn(err)
		inCache = false
	}

	if !inCache {
		m.logger.Info("Getting movies ids from repository")
		moviesIds, err = m.moviesPreviewRepo.GetMoviesPreview(ctx, filter, limit, offset)
		if err != nil {
			return []MoviePreview{}, err
		}

		go func() {
			m.logger.Info("Caching filtered request")
			if err := m.moviesPreviewCache.CacheFilteredRequest(context.Background(),
				filter, limit, offset, moviesIds, m.cfg.FilteredTTL); err != nil {
				m.logger.Error(err)
			}
		}()
	}

	m.logger.Info("Checking movies ids nums")
	if len(moviesIds) == 0 {
		return []MoviePreview{}, ErrNotFound
	}

	m.logger.Info("Filling movies")
	var movies = make([]MoviePreview, 0, len(moviesIds))
	for _, id := range moviesIds {
		movie, err := m.GetMoviePreview(ctx, id)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return []MoviePreview{}, err
		}

		movies = append(movies, movie)

	}

	return movies, nil
}

func (m *RepositoryManager) GetAgeRatings(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetAgeRatings")
	defer span.Finish()
	var err error
	defer span.SetTag("has_error", err != nil)

	ratings, err := m.ageRatingsRepo.GetAgeRatings(ctx)
	if err != nil {
		return []string{}, err
	}

	return ratings, nil
}
