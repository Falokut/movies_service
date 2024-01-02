package repository

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type countriesRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	countriesTableName       = "countries"
	moviesCountriesTableName = "movies_countries"
)

func NewCountriesRepository(db *sqlx.DB, logger *logrus.Logger) *countriesRepository {
	return &countriesRepository{db: db, logger: logger}
}

func (r *countriesRepository) Shutdown() {
	r.db.Close()
}

func (r *countriesRepository) GetCountries(ctx context.Context, movieId int32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetAllCountries")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT name FROM %s JOIN %s ON country_id=id WHERE movie_id=$1 ORDER BY id",
		moviesCountriesTableName, countriesTableName)
	var countries []string
	err = r.db.SelectContext(ctx, &countries, query, movieId)
	if err != nil {
		r.logger.Errorf("error: %v query: %s args: %v", err, query, movieId)
		return []string{}, err
	}

	return countries, nil
}

func (r *countriesRepository) GetAllCountries(ctx context.Context) ([]Country, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetAllCountries")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", countriesTableName)
	var countries []Country
	err = r.db.SelectContext(ctx, &countries, query)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return []Country{}, err
	}

	return countries, nil
}

func (r *countriesRepository) GetCountriesForMovies(ctx context.Context, ids []string) (map[int32][]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "countriesRepository.GetCountriesForMovies")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT movie_id, ARRAY_AGG(name) FROM %s JOIN %s ON country_id=id "+
		"WHERE movie_id=ANY(ARRAY[%s]) GROUP BY movie_id",
		moviesCountriesTableName, countriesTableName, strings.Join(ids, ","))
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return map[int32][]string{}, err
	}

	countries := make(map[int32][]string, len(ids))
	for rows.Next() {
		var id int32
		var names string

		rows.Scan(&id, &names)
		countries[id] = strings.Split(strings.Trim(names, "{}"), ",")
	}
	return countries, nil
}
