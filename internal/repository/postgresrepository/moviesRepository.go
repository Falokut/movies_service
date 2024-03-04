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

type MoviesRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	moviesTableName = "movies"
)

func NewMoviesRepository(db *sqlx.DB, logger *logrus.Logger) *MoviesRepository {
	return &MoviesRepository{db: db, logger: logger}
}

func (r *MoviesRepository) GetMovie(ctx context.Context, movieId int32) (movie models.RepositoryMovie, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetMovie")

	query := fmt.Sprintf(`SELECT %[1]s.id, title_ru, COALESCE(title_en,'') AS title_en,
		description, duration, COALESCE(poster_picture_id,'') AS poster_picture_id,
		COALESCE(background_picture_id,'') AS background_picture_id,
		release_year, COALESCE(%[2]s.name,'') AS age_rating 
		FROM %[1]s LEFT JOIN %[2]s ON age_rating_id=%[2]s.id WHERE %[1]s.id=$1`, moviesTableName, ageRatingsTableName)
	err = r.db.GetContext(ctx, &movie, query, movieId)
	if err != nil {
		return
	}

	return movie, nil
}

func (r *MoviesRepository) GetMoviesPreviewIds(ctx context.Context,
	filter models.MoviesFilter, limit, offset uint32) (ids []string, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetMoviesPreviewIds")

	query := fmt.Sprintf("SELECT %[1]s.id FROM %[1]s LEFT JOIN %[2]s ON age_rating_id=%[2]s.id "+
		"%[3]s ORDER BY %[1]s.id LIMIT %[4]d OFFSET %[5]d;", moviesTableName, ageRatingsTableName, convertFilterToWhere(filter),
		limit, offset)
	err = r.db.SelectContext(ctx, &ids, query)
	if err != nil {
		return
	}

	return ids, nil
}

func (r *MoviesRepository) GetMoviePreview(ctx context.Context, movieId int32) (moviePreview models.RepositoryMoviePreview, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetMoviePreview")

	query := fmt.Sprintf(`SELECT %[1]s.id, title_ru,
		COALESCE(title_en,'') AS title_en, duration, COALESCE(preview_poster_picture_id,'') AS preview_poster_picture_id,
		short_description, release_year, COALESCE(%[2]s.name,'') AS age_rating 
		FROM %[1]s LEFT JOIN %[2]s ON age_rating_id=%[2]s.id WHERE %[1]s.id=$1`, moviesTableName, ageRatingsTableName)
	err = r.db.GetContext(ctx, &moviePreview, query, movieId)
	return
}

func (r *MoviesRepository) GetMovies(ctx context.Context, ids []int32) (movies []models.RepositoryMoviePreview, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetMovies")

	query := fmt.Sprintf(`SELECT %[1]s.id, title_ru, COALESCE(title_en,'') AS title_en,
	duration, COALESCE(preview_poster_picture_id,'') AS preview_poster_picture_id,
	short_description, release_year, COALESCE(%[2]s.name,'') AS age_rating 
	FROM %[1]s LEFT JOIN %[2]s ON age_rating_id=%[2]s.id WHERE %[1]s.id=ANY($1)`, moviesTableName, ageRatingsTableName)
	err = r.db.SelectContext(ctx, &movies, query, ids)

	return movies, nil
}

// if any filter param filled, return string with WHERE statement
func convertFilterToWhere(filter models.MoviesFilter) string {
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

	filter := fmt.Sprintf("%s.id=ANY(SELECT movie_id FROM %s WHERE genre_id=ANY(ARRAY[%s]))",
		moviesTableName, moviesGenresTableName, ids)
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

	filter := fmt.Sprintf("%s.id=ANY(SELECT movie_id FROM %s WHERE country_id=ANY(ARRAY[%s]))", moviesTableName, moviesCountriesTableName, ids)
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

func (r *MoviesRepository) logError(err error, functionName string) {
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
		).Error("movies repository error occurred")
	} else {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("movies repository error occurred")
	}
}
