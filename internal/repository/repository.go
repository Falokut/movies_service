package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("entity not found")

type Movie struct {
	ID           string         `db:"id"`
	TitleRU      string         `db:"title_ru"`
	TitleEN      sql.NullString `db:"title_en"`
	Budget       sql.NullString `db:"budget"`
	Plot         string         `db:"plot"`
	Genres       sql.NullString `db:"genres"`
	CastID       int32          `db:"cast_id"`
	Duration     int32          `db:"duration"`
	PictureID    sql.NullString `db:"poster_picture_id"`
	DirectorsIDs sql.NullString `db:"directors"`
	CountriesIDs sql.NullString `db:"countries"`
	ReleaseYear  int32          `db:"release_year"`
}

type MoviesFilter struct {
	MoviesIDs    string
	GenresIDs    string
	DirectorsIDs string
	CountriesIDs string
	Title        string
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
	GetMovies(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]Movie, error)
}

type MoviesRepository interface {
	GetMovie(ctx context.Context, movieId string) (Movie, error)
	GetMovies(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]string, error)
}

type MoviesCache interface {
	GetMovie(ctx context.Context, movieId string) (Movie, error)
	GetMovies(ctx context.Context, Filter MoviesFilter, limit, offset uint32) ([]string, error)
	CacheMovies(ctx context.Context, movies []Movie, TTL time.Duration) error
	CacheFilteredRequest(ctx context.Context, Filter MoviesFilter, limit, offset uint32, moviesIDs []string, TTL time.Duration) error
}
