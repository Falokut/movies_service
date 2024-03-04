package service

import (
	"context"

	"github.com/Falokut/movies_service/internal/models"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=service.go -destination=mocks/repository.go
type MoviesRepository interface {
	GetMovie(ctx context.Context, movieId int32) (models.RepositoryMovie, error)
	GetMoviePreview(ctx context.Context, movieId int32) (models.RepositoryMoviePreview, error)
	GetMoviesPreview(ctx context.Context, Filter models.MoviesFilter, limit, offset uint32) ([]models.RepositoryMoviePreview, error)
	GetMoviesPreviewByIDs(ctx context.Context, ids []string) ([]models.RepositoryMoviePreview, error)
	GetAgeRatings(ctx context.Context) ([]string, error)
	GetGenres(ctx context.Context, movieId int32) ([]string, error)
	GetAllGenres(ctx context.Context) ([]models.Genre, error)
	GetCountries(ctx context.Context, movieId int32) ([]string, error)
	GetAllCountries(ctx context.Context) ([]models.Country, error)
}

type MoviesService interface {
	GetMovie(ctx context.Context, id int32) (models.Movie, error)
	GetMoviesPreview(ctx context.Context, filter models.MoviesFilter, limit, offset uint32) ([]models.MoviePreview, error)
	GetMoviesPreviewByIDs(ctx context.Context, ids []string) ([]models.MoviePreview, error)
	GetAgeRatings(ctx context.Context) ([]string, error)
	GetGenres(ctx context.Context) ([]models.Genre, error)
	GetCountries(ctx context.Context) ([]models.Country, error)
}

type PicturesUrlConfig struct {
	BaseUrl                string `yaml:"base_url" env:"BASE_URL"`
	PostersCategory        string `yaml:"posters_category" env:"POSTERS_CATEGORY"`
	PreviewPostersCategory string `yaml:"preview_posters_category" env:"PREVIEW_POSTERS_CATEGORY"`
	BackgroundsCategory    string `yaml:"backgrounds_category" env:"BACKGROUNDS_CATEGORY"`
}

type moviesService struct {
	logger      *logrus.Logger
	repo        MoviesRepository
	picturesCfg PicturesUrlConfig
}

func NewMoviesService(logger *logrus.Logger, repo MoviesRepository,
	picturesCfg PicturesUrlConfig) *moviesService {
	return &moviesService{
		logger:      logger,
		repo:        repo,
		picturesCfg: picturesCfg,
	}
}

func (s *moviesService) GetMovie(ctx context.Context, id int32) (movie models.Movie, err error) {
	res, err := s.repo.GetMovie(ctx, id)
	if err != nil {
		return
	}
	return models.Movie{
		Description: res.Description,
		TitleRU:     res.TitleRU,
		TitleEN:     res.TitleEN,
		Genres:      res.Genres,
		Duration:    res.Duration,
		Countries:   res.Countries,
		PosterUrl: getPictureURL(res.PosterID,
			s.picturesCfg.BaseUrl, s.picturesCfg.PostersCategory),
		BackgroundUrl: getPictureURL(res.BackgroundPictureID,
			s.picturesCfg.BaseUrl, s.picturesCfg.BackgroundsCategory),
		ReleaseYear: res.ReleaseYear,
		AgeRating:   res.AgeRating,
	}, nil
}

func (s *moviesService) GetAgeRatings(ctx context.Context) (ratings []string, err error) {
	ratings, err = s.repo.GetAgeRatings(ctx)
	if err != nil {
		return
	}

	return
}

func (s *moviesService) GetMoviesPreview(ctx context.Context,
	filter models.MoviesFilter, limit, offset uint32) (movies []models.MoviePreview, err error) {

	res, err := s.repo.GetMoviesPreview(ctx, filter, limit, offset)
	if err != nil {
		return
	}
	movies = make([]models.MoviePreview, len(res))
	for i := range res {
		movies[i] = models.MoviePreview{
			ID:        res[i].ID,
			TitleRU:   res[i].TitleRU,
			TitleEN:   res[i].TitleEN,
			Genres:    res[i].Genres,
			Countries: res[i].Countries,
			Duration:  res[i].Duration,
			PreviewPosterUrl: getPictureURL(res[i].PreviewPosterID,
				s.picturesCfg.BaseUrl, s.picturesCfg.PreviewPostersCategory),
			ShortDescription: res[i].ShortDescription,
			ReleaseYear:      res[i].ReleaseYear,
			AgeRating:        res[i].AgeRating,
		}
	}
	return
}

func (s *moviesService) GetMoviesPreviewByIDs(ctx context.Context,
	ids []string) (movies []models.MoviePreview, err error) {

	res, err := s.repo.GetMoviesPreviewByIDs(ctx, ids)
	if err != nil {
		return
	}
	movies = make([]models.MoviePreview, len(res))
	for i := range res {
		movies[i] = models.MoviePreview{
			ID:        res[i].ID,
			TitleRU:   res[i].TitleRU,
			TitleEN:   res[i].TitleEN,
			Genres:    res[i].Genres,
			Countries: res[i].Countries,
			Duration:  res[i].Duration,
			PreviewPosterUrl: getPictureURL(res[i].PreviewPosterID,
				s.picturesCfg.BaseUrl, s.picturesCfg.PreviewPostersCategory),
			ShortDescription: res[i].ShortDescription,
			ReleaseYear:      res[i].ReleaseYear,
			AgeRating:        res[i].AgeRating,
		}
	}
	return
}

func (s *moviesService) GetGenres(ctx context.Context) ([]models.Genre, error) {
	return s.repo.GetAllGenres(ctx)
}

func (s *moviesService) GetCountries(ctx context.Context) (countries []models.Country, err error) {
	return s.repo.GetAllCountries(ctx)
}

func getPictureURL(pictureID, baseUrl, category string) string {
	if pictureID == "" || baseUrl == "" || category == "" {
		return ""
	}

	return baseUrl + "/" + category + "/" + pictureID
}
