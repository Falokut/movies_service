package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("entity not found")

type MoviePreview struct {
	ID               int32          `db:"id"`
	TitleRU          string         `db:"title_ru"`
	TitleEN          sql.NullString `db:"title_en"`
	Duration         int32          `db:"duration"`
	PreviewPosterID  sql.NullString `db:"preview_poster_picture_id"`
	Genres           []string       `db:"genres"`
	ShortDescription string         `db:"short_description"`
	Countries        []string       `db:"countries"`
	ReleaseYear      int32          `db:"release_year"`
	AgeRating        string         `db:"age_rating"`
}

type Movie struct {
	ID                  int32          `db:"id"`
	TitleRU             string         `db:"title_ru"`
	TitleEN             sql.NullString `db:"title_en"`
	Description         string         `db:"description"`
	Genres              []string       `db:"genres"`
	Duration            int32          `db:"duration"`
	PosterID            sql.NullString `db:"poster_picture_id"`
	BackgroundPictureID sql.NullString `db:"background_picture_id"`
	Countries           []string       `db:"countries"`
	ReleaseYear         int32          `db:"release_year"`
	AgeRating           string         `db:"age_rating"`
}

type MoviesFilter struct {
	MoviesIDs    string
	GenresIDs    string
	CountriesIDs string
	Title        string
	AgeRating    string
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

//go:generate mockgen -source=repository.go -destination=mocks/repository.go
type MoviesRepositoryManager interface {
	GetMovie(ctx context.Context, movieId int32) (Movie, error)
	GetMoviePreview(ctx context.Context, movieId int32) (MoviePreview, error)
	GetMoviesPreview(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]MoviePreview, error)
	GetAgeRatings(ctx context.Context) ([]string, error)
	GetGenres(ctx context.Context, movieId int32) ([]string, error)
	GetAllGenres(ctx context.Context) ([]Genre, error)
	GetCountries(ctx context.Context, movieId int32) ([]string, error)
	GetAllCountries(ctx context.Context) ([]Country, error)
}

type MoviesRepository interface {
	GetMovie(ctx context.Context, movieId int32) (Movie, error)
}

type AgeRatingRepository interface {
	GetAgeRatings(ctx context.Context) ([]string, error)
}

type MoviesPreviewRepository interface {
	GetMoviesPreviewIds(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]string, error)
	GetMoviePreview(ctx context.Context, movieId int32) (MoviePreview, error)
	GetMovies(ctx context.Context, ids []int32) ([]MoviePreview, error)
}

type MoviesCache interface {
	GetMovie(ctx context.Context, movieId int32) (Movie, error)
	CacheMovies(ctx context.Context, movies []Movie, ttl time.Duration) error
}

type MoviesPreviewCache interface {
	GetMovie(ctx context.Context, movieId int32) (MoviePreview, error)
	GetMoviesIDs(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]string, error)
	GetMovies(ctx context.Context, ids []string) ([]MoviePreview, []string, error)

	CacheMovies(ctx context.Context, movies []MoviePreview, ttl time.Duration) error
	CacheFilteredRequest(ctx context.Context, Filter MoviesFilter, limit, offset uint32, moviesIDs []string, ttl time.Duration) error
}

type Genre struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type GenresRepository interface {
	GetGenres(ctx context.Context, movieId int32) ([]string, error)
	GetGenresForMovies(ctx context.Context, ids []string) (map[int32][]string, error)
	GetAllGenres(ctx context.Context) ([]Genre, error)
}

type Country struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type CountryRepository interface {
	GetCountries(ctx context.Context, movieId int32) ([]string, error)
	GetCountriesForMovies(ctx context.Context, ids []string) (map[int32][]string, error)
	GetAllCountries(ctx context.Context) ([]Country, error)
}
