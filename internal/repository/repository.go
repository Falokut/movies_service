package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("entity not found")

type MoviePreview struct {
	ID               string         `db:"id"`
	TitleRU          string         `db:"title_ru"`
	TitleEN          sql.NullString `db:"title_en"`
	Duration         int32          `db:"duration"`
	PreviewPosterID  sql.NullString `db:"preview_poster_picture_id"`
	Genres           sql.NullString `db:"genres"`
	ShortDescription string         `db:"short_description"`
	CountriesIDs     sql.NullString `db:"countries"`
	ReleaseYear      int32          `db:"release_year"`
	AgeRating        string         `db:"age_rating"`
}

type Movie struct {
	ID                  string         `db:"id"`
	TitleRU             string         `db:"title_ru"`
	TitleEN             sql.NullString `db:"title_en"`
	Description         string         `db:"description"`
	Genres              sql.NullString `db:"genres"`
	CastID              int32          `db:"cast_id"`
	Duration            int32          `db:"duration"`
	PosterID            sql.NullString `db:"poster_picture_id"`
	BackgroundPictureID sql.NullString `db:"background_picture_id"`
	DirectorsIDs        sql.NullString `db:"directors"`
	CountriesIDs        sql.NullString `db:"countries"`
	ReleaseYear         int32          `db:"release_year"`
	AgeRating           string         `db:"age_rating"`
}

type MoviesFilter struct {
	MoviesIDs    string
	GenresIDs    string
	DirectorsIDs string
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

//go:generate mockgen -source=repository.go -destination=mocks/moviesRepositoryManager.go
type MoviesRepositoryManager interface {
	GetMovie(ctx context.Context, movieId string) (Movie, error)
	GetMoviePreview(ctx context.Context, movieId string) (MoviePreview, error)
	GetMoviesPreview(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]MoviePreview, error)
	GetAgeRatings(ctx context.Context) ([]string, error)
}

type MoviesRepository interface {
	GetMovie(ctx context.Context, movieId string) (Movie, error)
}

type AgeRatingRepository interface {
	GetAgeRatings(ctx context.Context) ([]string, error)
}

type MoviesPreviewRepository interface {
	GetMoviesPreview(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]string, error)
	GetMoviePreview(ctx context.Context, movieId string) (MoviePreview, error)
}

type MoviesCache interface {
	GetMovie(ctx context.Context, movieId string) (Movie, error)
	CacheMovies(ctx context.Context, movies []Movie, TTL time.Duration) error
}

type MoviesPreviewCache interface {
	GetMovie(ctx context.Context, movieId string) (MoviePreview, error)
	GetMovies(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]string, error)
	CacheMovies(ctx context.Context, movies []MoviePreview, TTL time.Duration) error
	CacheFilteredRequest(ctx context.Context, Filter MoviesFilter, limit, offset uint32, moviesIDs []string, TTL time.Duration) error
}
