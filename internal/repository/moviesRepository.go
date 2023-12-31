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
	"github.com/sirupsen/logrus"
)

type moviesRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	moviesTableName = "movies"
)

func NewMoviesRepository(db *sqlx.DB, logger *logrus.Logger) *moviesRepository {
	return &moviesRepository{db: db, logger: logger}
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

func (r *moviesRepository) GetMovie(ctx context.Context, movieId int32) (Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT %[1]s.id, title_ru, title_en,"+
		"description, duration, poster_picture_id,"+
		"background_picture_id, release_year, COALESCE(%[2]s.name,'') AS age_rating "+
		" FROM %[1]s LEFT JOIN %[2]s ON age_rating_id=%[2]s.id WHERE %[1]s.id=$1", moviesTableName, ageRatingsTableName)
	var movie Movie
	err = r.db.GetContext(ctx, &movie, query, movieId)
	if errors.Is(err, sql.ErrNoRows) {
		return Movie{}, ErrNotFound
	}
	if err != nil {
		r.logger.Errorf("error: %v query: %s args: %v", err, query, movieId)
		return Movie{}, err
	}

	return movie, nil
}

func (r *moviesRepository) GetMoviesPreview(ctx context.Context, filter MoviesFilter, limit, offset uint32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMoviesPreview")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT %[1]s.id AS id FROM %[1]s LEFT JOIN %[2]s ON age_rating_id=%[2]s.id "+
		"%[3]s ORDER BY id LIMIT %[4]d OFFSET %[5]d;", moviesTableName, ageRatingsTableName, convertFilterToWhere(filter),
		limit, offset)
	var ids []string
	err = r.db.SelectContext(ctx, &ids, query)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return []string{}, err
	}

	return ids, nil

}

// if any filter param filled, return string with WHERE statement
func convertFilterToWhere(filter MoviesFilter) string {
	statement := ""
	var first = true
	if len(filter.MoviesIDs) > 0 {
		statement += fmt.Sprintf(" %s.id=ANY(ARRAY[%s]) ", moviesTableName, filter.MoviesIDs)
		first = false
	}

	str, first := getGenresFilter(filter.GenresIDs, first)
	statement += str

	str, first = getCountriesFilter(filter.CountriesIDs, first)
	statement += str

	str, first = containsInArray(ageRatingsTableName+".name", filter.AgeRating, first)
	statement += str

	if filter.Title != "" {
		filter.Title = strings.ReplaceAll(strings.ToLower(filter.Title), "'", "''") + "%"
		str = fmt.Sprintf(" LOWER(title_ru) LIKE('%[1]s') OR LOWER(title_en) LIKE('%[1]s')", filter.Title)
		if !first {
			statement += " AND " + str
		} else {
			statement += str
		}
	}

	if statement == "" {
		return statement
	}

	return " WHERE " + statement
}

func getGenresFilter(ids string, first bool) (string, bool) {
	if len(ids) == 0 {
		return "", first
	}

	filter := fmt.Sprintf("id=ANY(SELECT movie_id FROM %s WHERE genre_id=ANY(ARRAY[%s]))", moviesGenresTableName, ids)
	if first {
		first = false
	} else {
		filter = " AND " + filter
	}
	return filter, first
}

func getCountriesFilter(ids string, first bool) (string, bool) {
	if len(ids) == 0 {
		return "", first
	}

	filter := fmt.Sprintf("id=ANY(SELECT movie_id FROM %s WHERE country_id=ANY(ARRAY[%s]))", moviesCountriesTableName, ids)
	if first {
		first = false
	} else {
		filter = " AND " + filter
	}
	return filter, first
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

func (r *moviesRepository) GetMoviePreview(ctx context.Context, movieId int32) (MoviePreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "moviesRepository.GetMoviePreview")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT %[1]s.id, title_ru, title_en, duration, preview_poster_picture_id,"+
		"short_description, release_year, COALESCE(%[2]s.name,'') AS age_rating "+
		"FROM %[1]s LEFT JOIN %[2]s ON age_rating_id=%[2]s.id WHERE %[1]s.id=$1", moviesTableName, ageRatingsTableName)
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
