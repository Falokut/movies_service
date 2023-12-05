package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type postgreRepository struct {
	db *sqlx.DB
}

const (
	moviesTableName = "movies"
)

func NewMoviesRepository(db *sqlx.DB) *postgreRepository {
	return &postgreRepository{db: db}
}

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *postgreRepository) Shutdown() {
	r.db.Close()
}

func (r *postgreRepository) GetMovie(ctx context.Context, movieId string) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgreRepository.GetMovie")
	defer span.Finish()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", moviesTableName)
	var movie Movie
	err := r.db.GetContext(ctx, &movie, query, movieId)
	if errors.Is(err, sql.ErrNoRows) {
		return Movie{}, ErrNotFound
	}
	if err != nil {
		return Movie{}, err
	}

	return movie, nil
}

func (r *postgreRepository) GetMovies(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgreRepository.GetMovies")
	defer span.Finish()

	query := fmt.Sprintf("SELECT id FROM %s %s ORDER BY id LIMIT %d OFFSET %d;", moviesTableName,
		convertFilterToWhere(filter), limit, offset)
	movies := []string{}
	err := r.db.SelectContext(ctx, &movies, query)
	if errors.Is(err, sql.ErrNoRows) {
		return []string{}, ErrNotFound
	}
	if err != nil {
		return []string{}, err
	}

	return movies, nil

}

// if any filter param filled, return string with WHERE statement
func convertFilterToWhere(filter MoviesFilter) string {
	statement := ""
	var first = true
	if len(filter.MoviesIDs) > 0 {
		statement += fmt.Sprintf(" id IN(%s) ", filter.MoviesIDs)
		first = false
	}

	str, first := arrayContains("genres", filter.GenresIDs, first)
	statement += str

	str, first = arrayContains("countries", filter.CountriesIDs, first)
	statement += str

	str, first = arrayContains("directors", filter.DirectorsIDs, first)
	statement += str

	if filter.Title != "" {
		filter.Title = strings.ReplaceAll(strings.ToLower(filter.Title), "'", "''")
		str = fmt.Sprintf(" LOWER(title_ru) LIKE('%[1]s%[2]s') OR LOWER(title_en) LIKE('%[1]s%[2]s')", filter.Title, "%")
		if !first {
			statement += " AND " + str
		} else {
			statement += str
		}
	}

	if statement == "" {
		return statement
	}

	return "WHERE" + statement
}

func arrayContains(fieldname string, parameter string, first bool) (string, bool) {
	if fieldname == "" || parameter == "" {
		return "", first
	}

	str := fmt.Sprintf(" %[1]s @> array[%[2]s] AND cardinality(%[1]s) > 0", fieldname, parameter)
	if first {
		return str, !first
	}
	return " AND " + str, false
}
