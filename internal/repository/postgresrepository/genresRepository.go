package postgresrepository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Falokut/movies_service/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type GenresRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	moviesGenresTableName = "movies_genres"
	genresTableName       = "genres"
)

func NewGenresRepository(db *sqlx.DB, logger *logrus.Logger) *GenresRepository {
	return &GenresRepository{db: db, logger: logger}
}

func (r *GenresRepository) GetGenres(ctx context.Context, movieId int32) (genres []string, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetGenres")

	query := fmt.Sprintf("SELECT name FROM %s JOIN %s ON genre_id=id WHERE movie_id=$1 ORDER BY id",
		moviesGenresTableName, genresTableName)

	err = r.db.SelectContext(ctx, &genres, query, movieId)
	return
}

func (r *GenresRepository) GetAllGenres(ctx context.Context) (genres []models.Genre, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetAllGenres")

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", genresTableName)

	err = r.db.SelectContext(ctx, &genres, query)

	return
}

func (r *GenresRepository) GetGenresForMovies(ctx context.Context, ids []string) (map[int32][]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GenresRepository.GetGenresForMovies")
	defer span.Finish()

	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT movie_id, ARRAY_AGG(name) FROM %s JOIN %s ON genre_id=id "+
		"WHERE movie_id=ANY(ARRAY[%s]) GROUP BY movie_id",
		moviesGenresTableName, genresTableName, strings.Join(ids, ","))
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return map[int32][]string{}, err
	}

	genres := make(map[int32][]string, len(ids))
	for rows.Next() {
		var id int32
		var names string

		rows.Scan(&id, &names)
		genres[id] = strings.Split(strings.Trim(names, "{}"), ",")
	}
	return genres, nil
}

func (r *GenresRepository) logError(err error, functionName string) {
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
		).Error("genres repository error occurred")
	} else {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("genres repository error occurred")
	}
}
