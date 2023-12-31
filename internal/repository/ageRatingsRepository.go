package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type ageRatingsRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	ageRatingsTableName = "age_ratings"
)

func NewAgeRatingsRepository(db *sqlx.DB, logger *logrus.Logger) *ageRatingsRepository {
	return &ageRatingsRepository{db: db, logger: logger}
}
func (r *ageRatingsRepository) Shutdown() {
	r.db.Close()
}

func (r *ageRatingsRepository) GetAgeRatings(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ageRatingsRepository.GetAgeRatings")
	defer span.Finish()
	var err error
	defer span.SetTag("error", err != nil)

	query := fmt.Sprintf("SELECT name FROM %s", ageRatingsTableName)
	var ratings = []string{}
	err = r.db.SelectContext(ctx, &ratings, query)
	if err != nil {
		r.logger.Errorf("error: %v query: %s", err, query)
		return []string{}, err
	}
	return ratings, nil
}
