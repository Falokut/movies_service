package postgresrepository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Falokut/movies_service/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CountriesRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	countriesTableName       = "countries"
	moviesCountriesTableName = "movies_countries"
)

func NewCountriesRepository(db *sqlx.DB, logger *logrus.Logger) *CountriesRepository {
	return &CountriesRepository{db: db, logger: logger}
}

func (r *CountriesRepository) GetCountries(ctx context.Context, movieId int32) (countries []string, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetCountries")

	query := fmt.Sprintf("SELECT name FROM %s JOIN %s ON country_id=id WHERE movie_id=$1 ORDER BY id",
		moviesCountriesTableName, countriesTableName)

	err = r.db.SelectContext(ctx, &countries, query, movieId)
	return
}

func (r *CountriesRepository) GetAllCountries(ctx context.Context) (countries []models.Country, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetAllCountries")

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", countriesTableName)

	err = r.db.SelectContext(ctx, &countries, query)
	return
}

func (r *CountriesRepository) GetCountriesForMovies(ctx context.Context, ids []string) (countries map[int32][]string, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetCountriesForMovies")

	query := fmt.Sprintf("SELECT movie_id, ARRAY_AGG(name) FROM %s JOIN %s ON country_id=id "+
		"WHERE movie_id=ANY(ARRAY[%s]) GROUP BY movie_id",
		moviesCountriesTableName, countriesTableName, strings.Join(ids, ","))
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return
	}

	countries = make(map[int32][]string, len(ids))
	for rows.Next() {
		var id int32
		var names string

		rows.Scan(&id, &names)
		countries[id] = strings.Split(strings.Trim(names, "{}"), ",")
	}

	return
}

func (r *CountriesRepository) logError(err error, functionName string) {
	if err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if errors.As(err, &repoErr) {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           repoErr.Msg,
				"error.code":          repoErr.Code,
			},
		).Error("countries repository error occurred")
	} else {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("countries repository error occurred")
	}
}
