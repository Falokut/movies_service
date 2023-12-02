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
	MovieTTL    time.Duration
	FilteredTTL time.Duration
}

type RepositoryManager struct {
	repo   MoviesRepository
	cache  MoviesCache
	cfg    RepositoryManagerConfig
	logger *logrus.Logger
}

func NewMoviesRepositoryManager(repo MoviesRepository, cache MoviesCache,
	cfg RepositoryManagerConfig, logger *logrus.Logger) *RepositoryManager {
	return &RepositoryManager{
		repo:   repo,
		cache:  cache,
		cfg:    cfg,
		logger: logger,
	}
}

func (m *RepositoryManager) GetMovie(ctx context.Context, movieID string) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMovie")
	defer span.Finish()

	movie, err := m.cache.GetMovie(ctx, movieID)
	if err == nil {
		return movie, nil
	}
	if !errors.Is(err, redis.Nil) {
		m.logger.Warn(err)
	}

	movie, err = m.repo.GetMovie(ctx, movieID)
	if err != nil {
		return Movie{}, err
	}

	go func(Movie) {
		m.cache.CacheMovies(context.Background(), []Movie{movie}, m.cfg.MovieTTL)
	}(movie)

	return movie, nil
}

func (m *RepositoryManager) GetMovies(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMovie")
	defer span.Finish()

	m.logger.Info("Trying get movies ids from cache")
	moviesIds, err := m.cache.GetMovies(ctx, Filter, limit, offset)
	inCache := true
	if err != nil {
		m.logger.Warn(err)
		inCache = false
	}

	if !inCache {
		m.logger.Info("Getting movies ids from repository")
		moviesIds, err = m.repo.GetMovies(ctx, Filter, limit, offset)
		if err != nil {
			return []Movie{}, err
		}

		go func() {
			m.logger.Info("Caching filtered request")
			m.cache.CacheFilteredRequest(context.Background(), Filter, limit, offset, moviesIds, m.cfg.FilteredTTL)
		}()
	}

	m.logger.Info("Checking movies ids nums")
	if len(moviesIds) == 0 {
		return []Movie{}, ErrNotFound
	}

	m.logger.Info("Filling movies")
	var movies = make([]Movie, 0, len(moviesIds))
	for _, id := range moviesIds {
		movie, err := m.GetMovie(ctx, id)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return []Movie{}, err
		}

		movies = append(movies, movie)

	}

	return movies, nil
}
