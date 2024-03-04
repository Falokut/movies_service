package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/Falokut/movies_service/internal/models"
	"github.com/sirupsen/logrus"
)

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type RepositoryConfig struct {
	MovieTTL        time.Duration
	FilteredTTL     time.Duration
	MoviePreviewTTL time.Duration
}

type MoviesRepository interface {
	GetMovie(ctx context.Context, movieId int32) (models.RepositoryMovie, error)
}

type AgeRatingRepository interface {
	GetAgeRatings(ctx context.Context) ([]string, error)
}

type MoviesPreviewRepository interface {
	GetMoviesPreviewIds(ctx context.Context, Filter models.MoviesFilter, limit, offset uint32) ([]string, error)
	GetMoviePreview(ctx context.Context, movieId int32) (models.RepositoryMoviePreview, error)
	GetMovies(ctx context.Context, ids []int32) ([]models.RepositoryMoviePreview, error)
}

type MoviesCache interface {
	GetMovie(ctx context.Context, movieId int32) (models.RepositoryMovie, error)
	CacheMovies(ctx context.Context, movies []models.RepositoryMovie, ttl time.Duration) error
}

type MoviesPreviewCache interface {
	GetMovie(ctx context.Context, movieId int32) (models.RepositoryMoviePreview, error)
	GetMoviesIDs(ctx context.Context, Filter models.MoviesFilter, limit, offset uint32) ([]string, error)
	GetMovies(ctx context.Context, ids []string) ([]models.RepositoryMoviePreview, []string, error)

	CacheMovies(ctx context.Context, movies []models.RepositoryMoviePreview, ttl time.Duration) error
	CacheFilteredRequest(ctx context.Context, Filter models.MoviesFilter,
		limit, offset uint32, moviesIDs []string, ttl time.Duration) error
}

type GenresRepository interface {
	GetGenres(ctx context.Context, movieId int32) ([]string, error)
	GetGenresForMovies(ctx context.Context, ids []string) (map[int32][]string, error)
	GetAllGenres(ctx context.Context) ([]models.Genre, error)
}

type CountryRepository interface {
	GetCountries(ctx context.Context, movieId int32) ([]string, error)
	GetCountriesForMovies(ctx context.Context, ids []string) (map[int32][]string, error)
	GetAllCountries(ctx context.Context) ([]models.Country, error)
}

type moviesRepository struct {
	moviesRepo         MoviesRepository
	moviesCache        MoviesCache
	moviesPreviewRepo  MoviesPreviewRepository
	moviesPreviewCache MoviesPreviewCache
	ageRatingsRepo     AgeRatingRepository
	genresRepo         GenresRepository
	countriesRepo      CountryRepository

	cfg    RepositoryConfig
	logger *logrus.Logger
}

func NewMoviesRepository(moviesRepo MoviesRepository,
	moviesCache MoviesCache,
	moviesPreviewRepo MoviesPreviewRepository,
	moviesPreviewCache MoviesPreviewCache,
	ageRatingsRepo AgeRatingRepository,
	genresRepo GenresRepository,
	countriesRepo CountryRepository,
	cfg RepositoryConfig,
	logger *logrus.Logger) *moviesRepository {
	return &moviesRepository{
		moviesRepo:         moviesRepo,
		moviesCache:        moviesCache,
		moviesPreviewCache: moviesPreviewCache,
		moviesPreviewRepo:  moviesPreviewRepo,
		ageRatingsRepo:     ageRatingsRepo,
		genresRepo:         genresRepo,
		countriesRepo:      countriesRepo,
		cfg:                cfg,
		logger:             logger,
	}
}

func (r *moviesRepository) GetMovie(ctx context.Context, movieID int32) (movie models.RepositoryMovie, err error) {
	movie, err = r.moviesCache.GetMovie(ctx, movieID)
	if err == nil {
		return
	}
	if models.Code(err) != models.NotFound {
		r.logger.Warn(err)
	}
	err = nil

	movieCh := make(chan models.RepositoryMovie, 1)
	errCh := make(chan error, 1)
	genresAndCountriesCh := make(chan struct{ genres, countries []string }, 1)

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		m, err := r.moviesRepo.GetMovie(reqCtx, movieID)
		if err != nil {
			errCh <- err
			r.logger.Warn(err)

			return
		}
		movieCh <- m
	}()
	go func() {
		genres, countries, err := r.getGenresAndCountriesForMovie(reqCtx, movieID)
		if err != nil {
			errCh <- err
			r.logger.Warn(err)

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
			return
		case str := <-genresAndCountriesCh:
			if genresAndCountriesDone {
				break
			}
			movie.Genres = str.genres
			movie.Countries = str.countries
			genresAndCountriesDone = true
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
		case cherr := <-errCh:
			if cherr != nil {
				err = cherr
				return
			}
		}
	}

	go func(mov models.RepositoryMovie) {
		if err := r.moviesCache.CacheMovies(context.Background(), []models.RepositoryMovie{mov}, r.cfg.MovieTTL); err != nil {
			r.logger.Error(err)
		}
	}(movie)

	return

}

func (r *moviesRepository) getGenresAndCountriesForMovie(ctx context.Context,
	movieId int32) (genres, countries []string, err error) {
	genresCh, countriesCh := make(chan []string, 1), make(chan []string, 1)
	errCh := make(chan error, 1)

	go func() {
		res, err := r.GetGenres(ctx, movieId)
		if err != nil {
			errCh <- err
			return
		}
		genresCh <- res
	}()
	go func() {
		res, err := r.GetCountries(ctx, movieId)
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
			err = models.Error(models.Canceled, "")
			return
		case g := <-genresCh:
			if genresDone {
				break
			}
			genres = g
			genresDone = true
		case c := <-countriesCh:
			if countriesDone {
				break
			}
			countries = c
			countriesDone = true
		case err = <-errCh:
			return
		}
	}

	return genres, countries, nil
}

func (r *moviesRepository) GetMoviePreview(ctx context.Context, movieID int32) (movie models.RepositoryMoviePreview, err error) {
	movie, err = r.moviesPreviewCache.GetMovie(ctx, movieID)
	if err == nil {
		return
	}
	if models.Code(err) != models.NotFound {
		r.logger.Warn(err)
	}

	movie, err = r.moviesPreviewRepo.GetMoviePreview(ctx, movieID)
	if err != nil {
		return models.RepositoryMoviePreview{}, err
	}

	movieCh := make(chan models.RepositoryMoviePreview, 1)
	errCh := make(chan error, 1)
	genresAndCountriesCh := make(chan struct{ genres, countries []string }, 1)

	reqCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		movie, err := r.moviesPreviewRepo.GetMoviePreview(reqCtx, movieID)
		if err != nil {
			errCh <- err
			return
		}
		movieCh <- movie
	}()
	go func() {
		genres, countries, err := r.getGenresAndCountriesForMovie(reqCtx, movieID)
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
			err = models.Error(models.Canceled, "")
			return
		case struc := <-genresAndCountriesCh:
			if genresAndCountriesDone {
				break
			}
			movie.Genres = struc.genres
			movie.Countries = struc.countries
			genresAndCountriesDone = true
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
		case err = <-errCh:
			if err != nil {
				r.logger.Error(err)
				return models.RepositoryMoviePreview{}, err
			}
		}
	}

	go func(models.RepositoryMoviePreview) {
		if err := r.moviesPreviewCache.CacheMovies(context.Background(),
			[]models.RepositoryMoviePreview{movie}, r.cfg.MoviePreviewTTL); err != nil {
			r.logger.Error(err)
		}
	}(movie)

	return movie, nil
}

func (r *moviesRepository) getGenresAndCountriesForMovies(ctx context.Context,
	ids []string) (genres, countries map[int32][]string, err error) {

	genresCh := make(chan map[int32][]string, 1)
	countriesCh := make(chan map[int32][]string, 1)

	errCh := make(chan error, 1)

	go func() {
		res, err := r.genresRepo.GetGenresForMovies(ctx, ids)
		if err != nil {
			errCh <- err
			return
		}
		genresCh <- res
	}()
	go func() {
		res, err := r.countriesRepo.GetCountriesForMovies(ctx, ids)
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
			err = models.Error(models.Canceled, "")
			return
		case g := <-genresCh:
			if genresDone {
				break
			}
			genres = g
			genresDone = true
		case c := <-countriesCh:
			if countriesDone {
				break
			}
			countries = c
			countriesDone = true
		case err = <-errCh:
			return
		}
	}

	return genres, countries, nil
}

func (r *moviesRepository) GetMoviesPreviewByIDs(ctx context.Context, ids []string) (movies []models.RepositoryMoviePreview, err error) {
	if len(ids) == 1 {
		id, _ := strconv.Atoi(ids[0])
		movie, err := r.GetMoviePreview(ctx, int32(id))
		return []models.RepositoryMoviePreview{movie}, err
	}

	r.logger.Info("Searching previews in cache")
	cachedPreviews, notFoundedIds, err := r.moviesPreviewCache.GetMovies(ctx, ids)
	if err != nil {
		r.logger.Error(err)
	}

	if len(cachedPreviews) == len(ids) {
		return cachedPreviews, nil
	}

	if len(cachedPreviews) != 0 && err == nil {
		ids = notFoundedIds
	}

	moviesPreviewsCh := make(chan []models.RepositoryMoviePreview, 1)
	countriesAndGenresCh := make(chan struct{ genres, countries map[int32][]string }, 1)
	errCh := make(chan error, 1)
	go func() {
		r.logger.Info("Searching genres and countries previews in repo")
		genres, countries, err := r.getGenresAndCountriesForMovies(ctx, ids)
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
		r.logger.Info("Searching previews in repo")
		previews, err := r.moviesPreviewRepo.GetMovies(ctx, convertStringsIntoInt(ids))
		if err != nil {
			errCh <- err
		}
		moviesPreviewsCh <- previews
		close(moviesPreviewsCh)
	}()

	var moviesDone, countriesAndMoviesDone bool
	var countries, genres map[int32][]string

	for !moviesDone || !countriesAndMoviesDone {
		select {
		case <-ctx.Done():
			return []models.RepositoryMoviePreview{}, ctx.Err()
		case mov := <-moviesPreviewsCh:
			if moviesDone {
				break
			}
			movies = mov
			moviesDone = true
		case st := <-countriesAndGenresCh:
			if countriesAndMoviesDone {
				break
			}
			countries, genres = st.countries, st.genres
			countriesAndMoviesDone = true
		}
	}

	r.logger.Info("Filling preview")
	var moviesFromRepo = make([]models.RepositoryMoviePreview, len(movies))
	for i := range movies {
		movies[i].Countries = countries[movies[i].ID]
		movies[i].Genres = genres[movies[i].ID]
		moviesFromRepo[i] = movies[i]
	}

	go func() {
		r.moviesPreviewCache.CacheMovies(context.Background(), moviesFromRepo, r.cfg.MoviePreviewTTL)
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

func (r *moviesRepository) GetMoviesPreview(ctx context.Context,
	filter models.MoviesFilter, limit, offset uint32) (movies []models.RepositoryMoviePreview, err error) {

	r.logger.Info("Trying get movies ids from cache")
	moviesIds, err := r.moviesPreviewCache.GetMoviesIDs(ctx, filter, limit, offset)
	inCache := true
	if err != nil {
		r.logger.Warn(err)
		inCache = false
		err = nil
	}

	if !inCache {
		r.logger.Info("Getting movies ids from repository")
		moviesIds, err = r.moviesPreviewRepo.GetMoviesPreviewIds(ctx, filter, limit, offset)
		if err != nil {
			return
		}

		go func() {
			r.logger.Info("Caching filtered request")
			if err := r.moviesPreviewCache.CacheFilteredRequest(context.Background(),
				filter, limit, offset, moviesIds, r.cfg.FilteredTTL); err != nil {
				r.logger.Error(err)
			}
		}()
	}

	r.logger.Info("Checking movies ids nums")
	if len(moviesIds) == 0 {
		return
	}

	r.logger.Info("Filling movies")
	movies, err = r.GetMoviesPreviewByIDs(ctx, moviesIds)
	if err != nil {
		return
	}

	return movies, nil
}

func (r *moviesRepository) GetAgeRatings(ctx context.Context) ([]string, error) {
	return r.ageRatingsRepo.GetAgeRatings(ctx)
}

func (r *moviesRepository) GetGenres(ctx context.Context, movieId int32) ([]string, error) {
	return r.genresRepo.GetGenres(ctx, movieId)
}

func (r *moviesRepository) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	return r.genresRepo.GetAllGenres(ctx)
}

func (r *moviesRepository) GetCountries(ctx context.Context, movieId int32) ([]string, error) {
	return r.countriesRepo.GetCountries(ctx, movieId)
}

func (r *moviesRepository) GetAllCountries(ctx context.Context) ([]models.Country, error) {
	return r.countriesRepo.GetAllCountries(ctx)
}
