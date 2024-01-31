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
		movie, err := m.moviesRepo.GetMovie(reqCtx, movieID)
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
		case str := <-genresAndCountriesCh:
			if genresAndCountriesDone {
				break
			}
			movie.Genres = str.genres
			movie.Countries = str.countries
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

	go func(mov Movie) {
		if err := m.moviesCache.CacheMovies(context.Background(), []Movie{mov}, m.cfg.MovieTTL); err != nil {
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
		movie, err := m.moviesPreviewRepo.GetMoviePreview(reqCtx, movieID)
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

func (m *RepositoryManager) getGenresAndCountriesForMovies(ctx context.Context, ids []string) (genres, countries map[int32][]string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.getGenresAndCountriesForMovie")
	defer span.Finish()
	defer span.SetTag("error", err != nil)

	genresCh, countriesCh := make(chan map[int32][]string, 1), make(chan map[int32][]string, 1)
	errCh := make(chan error, 1)

	go func() {
		res, err := m.genresRepo.GetGenresForMovies(ctx, ids)
		if err != nil {
			errCh <- err
			return
		}
		genresCh <- res
	}()
	go func() {
		res, err := m.countriesRepo.GetCountriesForMovies(ctx, ids)
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
			return map[int32][]string{}, map[int32][]string{}, ctx.Err()
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

func (m *RepositoryManager) GetMoviesPreviewByIDs(ctx context.Context, ids []string) ([]MoviePreview, error) {
	if len(ids) == 1 {
		id, _ := strconv.Atoi(ids[0])
		movie, err := m.GetMoviePreview(ctx, int32(id))
		return []MoviePreview{movie}, err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMoviesPreviewByIDs")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	m.logger.Info("Searching previews in cache")
	cachedPreviews, notFoundedIds, err := m.moviesPreviewCache.GetMovies(ctx, ids)
	if errors.Is(err, redis.Nil) {
		m.metrics.IncCacheMiss("GetMoviesPreviewByIDs", int32(len(ids)))
	} else if err != nil {
		m.logger.Error(err)
	}

	if len(cachedPreviews) == len(ids) {
		m.metrics.IncCacheHits("GetMoviesPreviewByIDs", int32(len(ids)))
		return cachedPreviews, nil
	}

	if len(cachedPreviews) != 0 && err == nil {
		m.metrics.IncCacheHits("GetMoviesPreviewByIDs", int32(len(ids)-len(notFoundedIds)))
		m.metrics.IncCacheMiss("GetMoviesPreviewByIDs", int32(len(notFoundedIds)))
		ids = notFoundedIds
	}

	moviesPreviewsCh := make(chan []MoviePreview, 1)
	countriesAndGenresCh := make(chan struct{ genres, countries map[int32][]string }, 1)
	errCh := make(chan error, 1)
	go func() {
		m.logger.Info("Searching genres and countries previews in repo")
		genres, countries, err := m.getGenresAndCountriesForMovies(ctx, ids)
		if err != nil {
			errCh <- err
			return
		}
		countriesAndGenresCh <- struct {
			genres    map[int32][]string
			countries map[int32][]string
		}{genres: genres, countries: countries}
	}()

	go func() {
		m.logger.Info("Searching previews in repo")
		previews, err := m.moviesPreviewRepo.GetMovies(ctx, convertStringsIntoInt(ids))
		if err != nil {
			errCh <- err
		}
		moviesPreviewsCh <- previews
		close(moviesPreviewsCh)
	}()

	var moviesDone, countriesAndMoviesDone bool
	var movies []MoviePreview
	var countries, genres map[int32][]string

	for !moviesDone || !countriesAndMoviesDone {
		select {
		case <-ctx.Done():
			return []MoviePreview{}, ctx.Err()
		case mov := <-moviesPreviewsCh:
			if moviesDone {
				break
			}
			movies = mov
			moviesDone = true
			m.logger.Info("movies done")
		case st := <-countriesAndGenresCh:
			if countriesAndMoviesDone {
				break
			}
			countries, genres = st.countries, st.genres
			countriesAndMoviesDone = true
			m.logger.Info("countries and genres done")
		}
	}

	m.logger.Info("Filling preview")
	var moviesFromRepo = make([]MoviePreview, 0, len(movies))
	for _, movie := range movies {
		movie.Countries = countries[movie.ID]
		movie.Genres = genres[movie.ID]
		m.logger.Debugf("countries, genres: %v %v", countries, genres)
		moviesFromRepo = append(moviesFromRepo, movie)
	}

	go func() {
		m.moviesPreviewCache.CacheMovies(context.Background(), moviesFromRepo, m.cfg.MoviePreviewTTL)
	}()

	return append(moviesFromRepo, cachedPreviews...), nil
}

func convertStringsIntoInt(s []string) []int32 {
	nums := make([]int32, 0, len(s))
	for _, num := range s {
		n, err := strconv.Atoi(num)
		if err != nil {
			continue
		}
		nums = append(nums, int32(n))
	}
	return nums
}

func (m *RepositoryManager) GetMoviesPreview(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetMoviesPreview")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	m.logger.Info("Trying get movies ids from cache")
	moviesIds, err := m.moviesPreviewCache.GetMoviesIDs(ctx, filter, limit, offset)
	inCache := true
	if err != nil {
		m.logger.Warn(err)
		inCache = false
		m.metrics.IncCacheMiss("GetMoviesPreview", 1)
	}
	if inCache {
		m.metrics.IncCacheHits("GetMoviesPreview", 1)
	} else {
		m.logger.Info("Getting movies ids from repository")
		moviesIds, err = m.moviesPreviewRepo.GetMoviesPreviewIds(ctx, filter, limit, offset)
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
	movies, err := m.GetMoviesPreviewByIDs(ctx, moviesIds)
	if err != nil {
		return []MoviePreview{}, err
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
