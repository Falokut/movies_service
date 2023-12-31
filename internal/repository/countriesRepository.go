package repository

import (
	"context"
	"fmt"

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
