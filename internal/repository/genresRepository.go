package repository

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type genresRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	moviesGenresTableName = "movies_genres"
	genresTableName       = "genres"
)

func NewGenresRepository(db *sqlx.DB, logger *logrus.Logger) *genresRepository {
	return &genresRepository{db: db, logger: logger}
}

func (r *genresRepository) Shutdown() {
	r.db.Close()
}

func (r *genresRepository) GetGenres(ctx context.Context, movieId int32) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.GetGenres")
	defer span.Finish()

	query := fmt.Sprintf("SELECT name FROM %s JOIN %s ON genre_id=id WHERE movie_id=$1 ORDER BY id",
		moviesGenresTableName, genresTableName)

	var genres []string
	err := r.db.SelectContext(ctx, &genres, query, movieId)
	if err != nil {
		r.logger.Errorf("error: %v query: %s args: %v", err, query, movieId)
		return []string{}, err
	}

	return genres, nil
}

func (r *genresRepository) GetAllGenres(ctx context.Context) ([]Genre, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genresRepository.GetAllGenres")
	defer span.Finish()

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id", genresTableName)

	var genres []Genre
	err := r.db.SelectContext(ctx, &genres, query)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return []Genre{}, err
	}

	return genres, nil
}
