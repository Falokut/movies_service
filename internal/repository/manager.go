package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type CacheMetric interface {
	IncCacheHits(method string, times int32)
	IncCacheMiss(method string, times int32)
}

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
	genresRepo         GenresRepository
	countriesRepo      CountryRepository

	metrics CacheMetric
	cfg     RepositoryManagerConfig
	logger  *logrus.Logger
}

func NewMoviesRepositoryManager(moviesRepo MoviesRepository, moviesCache MoviesCache, moviesPreviewRepo MoviesPreviewRepository,
	moviesPreviewCache MoviesPreviewCache, ageRatingsRepo AgeRatingRepository, genresRepo GenresRepository,
	countriesRepo CountryRepository, metrics CacheMetric,
	cfg RepositoryManagerConfig, logger *logrus.Logger) *RepositoryManager {
	return &RepositoryManager{
		moviesRepo:         moviesRepo,
		moviesCache:        moviesCache,
		moviesPreviewCache: moviesPreviewCache,
		moviesPreviewRepo:  moviesPreviewRepo,
		ageRatingsRepo:     ageRatingsRepo,
		genresRepo:         genresRepo,
		countriesRepo:      countriesRepo,
		metrics:            metrics,
		cfg:                cfg,
		logger:             logger,
	}
}

func (m *RepositoryManager) GetMovie(ctx context.Context, movieID int32) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	movie, err := m.moviesCache.GetMovie(ctx, movieID)
	if err == nil {
		m.metrics.IncCacheHits("GetMovie", 1)
		return movie, nil
	}
	if !errors.Is(err, redis.Nil) {
		m.logger.Warn(err)
	}

	m.metrics.IncCacheMiss("GetMovie", 1)
	movieCh := make(chan Movie, 1)
	errCh := make(chan error, 1)
	genresAndCountriesCh := make(chan struct{ genres, countries []string }, 1)

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		movie, err = m.moviesRepo.GetMovie(reqCtx, movieID)
		if err != nil {
			errCh <- err
			return
		}
		movieCh <- movie
	}()
	go func() {
		genres, countries, err := m.getGenresAndCountriesForMovie(reqCtx, movieID)
		if err != nil {
			errCh <- err
			return
		}
		genresAndCountriesCh <- struct {
			genres    []string
			countries []string
		}{genres: genres, countries: countries}
	}()

	var movieDone, genresAndCountriesDone bool
	for !genresAndCountriesDone || !movieDone {
		select {
		case <-reqCtx.Done():
			return Movie{}, reqCtx.Err()
		case struc := <-genresAndCountriesCh:
			if genresAndCountriesDone {
				break
			}
			movie.Genres = struc.genres
			movie.Countries = struc.countries
			genresAndCountriesDone = true
			m.logger.Debug("Genres and countries done")
		case res := <-movieCh:
			if movieDone {
				break
			}
			movie.ID = res.ID
			movie.TitleRU = res.TitleRU
			movie.TitleEN = res.TitleEN
			movie.Description = res.Description
			movie.Duration = res.Duration
			movie.PosterID = res.PosterID
			movie.BackgroundPictureID = res.BackgroundPictureID
			movie.ReleaseYear = res.ReleaseYear
			movie.AgeRating = res.AgeRating
			movieDone = true
			m.logger.Debug("Movie done")
		case err = <-errCh:
			if err != nil {
				m.logger.Error(err)
				return Movie{}, err
			}
		}
	}

	go func(Movie) {
		if err := m.moviesCache.CacheMovies(context.Background(), []Movie{movie}, m.cfg.MovieTTL); err != nil {
			m.logger.Error(err)
		}
	}(movie)
	return movie, nil

}

func (m *RepositoryManager) getGenresAndCountriesForMovie(ctx context.Context, movieId int32) (genres, countries []string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.getGenresAndCountriesForMovie")
	defer span.Finish()
	defer span.SetTag("error", err != nil)

	genresCh, countriesCh := make(chan []string, 1), make(chan []string, 1)
	errCh := make(chan error, 1)

	go func() {
		res, err := m.GetGenres(ctx, movieId)
		if err != nil {
			errCh <- err
			return
		}
		genresCh <- res
	}()
	go func() {
		res, err := m.GetCountries(ctx, movieId)
		if err != nil {
			errCh <- err
			return
		}
		countriesCh <- res
	}()

	var countriesDone, genresDone bool
	for !genresDone || !countriesDone {
		select {
		case <-ctx.Done():
			return []string{}, []string{}, ctx.Err()
		case g := <-genresCh:
			if genresDone {
				break
			}
			genres = g
			genresDone = true
			m.logger.Debug("genres done")
		case c := <-countriesCh:
			if countriesDone {
				break
			}
			countries = c
			countriesDone = true
			m.logger.Debug("countries done")
		case err = <-errCh:
			return
		}
	}

	return genres, countries, nil
}

func (m *RepositoryManager) GetMoviePreview(ctx context.Context, movieID int32) (MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMoviePreview")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

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

	m.metrics.IncCacheMiss("GetMovie", 1)
	movieCh := make(chan MoviePreview, 1)
	errCh := make(chan error, 1)
	genresAndCountriesCh := make(chan struct{ genres, countries []string }, 1)

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		movie, err = m.moviesPreviewRepo.GetMoviePreview(reqCtx, movieID)
		if err != nil {
			errCh <- err
			return
		}
		movieCh <- movie
	}()
	go func() {
		genres, countries, err := m.getGenresAndCountriesForMovie(reqCtx, movieID)
		if err != nil {
			errCh <- err
			return
		}
		genresAndCountriesCh <- struct {
			genres    []string
			countries []string
		}{genres: genres, countries: countries}
	}()

	var movieDone, genresAndCountriesDone bool
	for !genresAndCountriesDone || !movieDone {
		select {
		case <-reqCtx.Done():
			return MoviePreview{}, reqCtx.Err()
		case struc := <-genresAndCountriesCh:
			if genresAndCountriesDone {
				break
			}
			movie.Genres = struc.genres
			movie.Countries = struc.countries
			genresAndCountriesDone = true
			m.logger.Debug("Genres and countries done")
		case res := <-movieCh:
			if movieDone {
				break
			}
			movie.ID = res.ID
			movie.TitleRU = res.TitleRU
			movie.TitleEN = res.TitleEN
			movie.ShortDescription = res.ShortDescription
			movie.Duration = res.Duration
			movie.PreviewPosterID = res.PreviewPosterID
			movie.ReleaseYear = res.ReleaseYear
			movie.AgeRating = res.AgeRating
			movieDone = true
			m.logger.Debug("Movie done")
		case err = <-errCh:
			if err != nil {
				m.logger.Error(err)
				return MoviePreview{}, err
			}
		}
	}

	go func(MoviePreview) {
		if err := m.moviesPreviewCache.CacheMovies(context.Background(), []MoviePreview{movie}, m.cfg.MoviePreviewTTL); err != nil {
			m.logger.Error(err)
		}
	}(movie)

	return movie, nil
}

func (m *RepositoryManager) GetMoviesPreview(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMoviesPreview")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

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
		movieID, _ := strconv.Atoi(id)
		movie, err := m.GetMoviePreview(ctx, int32(movieID))
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
	return m.ageRatingsRepo.GetAgeRatings(ctx)
}

func (m *RepositoryManager) GetGenres(ctx context.Context, movieId int32) ([]string, error) {
	return m.genresRepo.GetGenres(ctx, movieId)
}

func (m *RepositoryManager) GetAllGenres(ctx context.Context) ([]Genre, error) {
	return m.genresRepo.GetAllGenres(ctx)
}

func (m *RepositoryManager) GetCountries(ctx context.Context, movieId int32) ([]string, error) {
	return m.countriesRepo.GetCountries(ctx, movieId)
}

func (m *RepositoryManager) GetAllCountries(ctx context.Context) ([]Country, error) {
	return m.countriesRepo.GetAllCountries(ctx)
}
