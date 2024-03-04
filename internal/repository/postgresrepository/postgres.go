package postgresrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Falokut/movies_service/internal/models"
	"github.com/Falokut/movies_service/internal/repository"
	"github.com/jmoiron/sqlx"
)

// NewPostgreDB creates a new connection to the PostgreSQL database.
func NewPostgreDB(cfg repository.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))

	return db, err
}

func handleError(ctx context.Context, err *error) {
	if ctx.Err() != nil {
		var code models.ErrorCode
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			code = models.Canceled
		case errors.Is(ctx.Err(), context.DeadlineExceeded):
			code = models.DeadlineExceeded
		}
		*err = models.Error(code, ctx.Err().Error())
		return
	}

	if err == nil || *err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if !errors.As(*err, &repoErr) {
		var code models.ErrorCode
		switch {
		case errors.Is(*err, sql.ErrNoRows):
			code = models.NotFound
			*err = models.Error(code, "entity not found")
		case *err != nil:
			code = models.Internal
			*err = models.Error(code, "repository internal error")
		}

	}
}
