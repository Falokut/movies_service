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

type moviesRepository struct {
	db *sqlx.DB
}

const (
	moviesTableName     = "movies"
	ageRatingsTableName = "age_ratings"
)

func NewMoviesRepository(db *sqlx.DB) *moviesRepository {
	return &moviesRepository{db: db}
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

func (r *moviesRepository) Shutdown() {
	r.db.Close()
}

func (r *moviesRepository) GetMovie(ctx context.Context, movieId string) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("has_errors", err != nil)

	query := fmt.Sprintf("SELECT id, title_ru, title_en,"+
		"description, genres, cast_id, duration, poster_picture_id,"+
		"background_picture_id, directors, countries, release_year, age_rating FROM %s WHERE id=$1", moviesTableName)
	var movie Movie
	err = r.db.GetContext(ctx, &movie, query, movieId)
	if errors.Is(err, sql.ErrNoRows) {
		return Movie{}, ErrNotFound
	}
	if err != nil {
		return Movie{}, err
	}

	return movie, nil
}

func (r *moviesRepository) GetMoviesPreview(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMoviesPreview")
	defer span.Finish()

	var err error
	defer span.SetTag("has_errors", err != nil)

	query := fmt.Sprintf("SELECT id FROM %s %s ORDER BY id LIMIT %d OFFSET %d;", moviesTableName,
		convertFilterToWhere(filter), limit, offset)
	movies := []string{}
	err = r.db.SelectContext(ctx, &movies, query)
	if errors.Is(err, sql.ErrNoRows) {
		return []string{}, ErrNotFound
	}
	if err != nil {
		return []string{}, err
	}

	return movies, nil

}

func (r *moviesRepository) GetAgeRatings(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetAgeRatings")
	defer span.Finish()
	var err error
	defer span.SetTag("has_errors", err != nil)

	query := fmt.Sprintf("SELECT name FROM %s", ageRatingsTableName)
	var ratings = []string{}
	err = r.db.SelectContext(ctx, &ratings, query)
	return ratings, err
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

	str, first = containsInArray("age_rating", filter.AgeRating, first)
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

func containsInArray(fieldname string, array string, first bool) (string, bool) {
	if fieldname == "" || array == "" {
		return "", first
	}
	str := fmt.Sprintf(" %s=ANY('{%s}')", fieldname, array)
	if first {
		return str, !first
	}
	return " AND " + str, false
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

func (r *moviesRepository) GetMoviePreview(ctx context.Context, movieId string) (MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMoviePreview")
	defer span.Finish()

	var err error
	defer span.SetTag("has_errors", err != nil)

	query := fmt.Sprintf("SELECT id, title_ru, title_en, duration, preview_poster_picture_id,"+
		"genres, short_description, countries, release_year, age_rating FROM %s WHERE id=$1", moviesTableName)
	var movie MoviePreview
	err = r.db.GetContext(ctx, &movie, query, movieId)
	if errors.Is(err, sql.ErrNoRows) {
		return MoviePreview{}, ErrNotFound
	}
	if err != nil {
		return MoviePreview{}, err
	}

	return movie, nil
}
