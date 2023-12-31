package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Falokut/movies_service/internal/repository"
	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go
type ImagesService interface {
	GetPictureURL(ctx context.Context, pictureID, baseUrl, category string) string
}

type PicturesUrlConfig struct {
	BaseUrl                string `yaml:"base_url" env:"BASE_URL"`
	PostersCategory        string `yaml:"posters_category" env:"POSTERS_CATEGORY"`
	PreviewPostersCategory string `yaml:"preview_posters_category" env:"PREVIEW_POSTERS_CATEGORY"`
	BackgroundsCategory    string `yaml:"backgrounds_category" env:"BACKGROUNDS_CATEGORY"`
}
type MoviesService struct {
	movies_service.UnimplementedMoviesServiceV1Server
	logger        *logrus.Logger
	repoManager   repository.MoviesRepositoryManager
	errorHandler  errorHandler
	imagesService ImagesService
	picturesCfg   PicturesUrlConfig
}

func NewMoviesService(logger *logrus.Logger, repoManager repository.MoviesRepositoryManager,
	imagesService ImagesService, picturesCfg PicturesUrlConfig) *MoviesService {
	errorHandler := newErrorHandler(logger)
	return &MoviesService{
		logger:        logger,
		repoManager:   repoManager,
		errorHandler:  errorHandler,
		imagesService: imagesService,
		picturesCfg:   picturesCfg,
	}
}

func (s *MoviesService) GetMovie(ctx context.Context, in *movies_service.GetMovieRequest) (*movies_service.Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetMovie")
	defer span.Finish()

	movie, err := s.repoManager.GetMovie(ctx, in.MovieID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return s.convertDbMovieToProto(ctx, movie), nil
}

func (s *MoviesService) GetAgeRatings(ctx context.Context, in *emptypb.Empty) (*movies_service.AgeRatings, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetAgeRatings")
	defer span.Finish()

	ratings, err := s.repoManager.GetAgeRatings(ctx)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return &movies_service.AgeRatings{Ratings: ratings}, nil
}

func (s *MoviesService) GetMoviesPreview(ctx context.Context, in *movies_service.GetMoviesPreviewRequest) (*movies_service.MoviesPreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetMoviesPreview")
	defer span.Finish()

	err := validateFilter(in)
	if errors.Is(err, ErrInvalidFilter) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
	}

	filter := repository.MoviesFilter{
		MoviesIDs:    ReplaceAll(in.GetMoviesIDs()),
		GenresIDs:    ReplaceAll(in.GetGenresIDs()),
		CountriesIDs: ReplaceAll(in.GetCountriesIDs()),
		AgeRating:    GetAgeRatingsFilter(in.GetAgeRatings()),
		Title:        ReplaceAll(in.GetTitle()),
	}
	s.logger.Print(filter.AgeRating)
	if in.Limit == 0 {
		in.Limit = 10
	} else if in.Limit > 100 {
		in.Limit = 100
	}

	Movies, err := s.repoManager.GetMoviesPreview(ctx, filter, in.Limit, in.Offset)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}
	if len(Movies) == 0 {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	}
	movies := make(map[int32]*movies_service.MoviePreview, len(Movies))
	for _, movie := range Movies {
		movies[movie.ID] = s.convertDbMoviePreviewToProto(ctx, movie)
	}

	span.SetTag("grpc.status", codes.OK)
	return &movies_service.MoviesPreview{Movies: movies}, nil
}

func (s *MoviesService) GetGenres(ctx context.Context, in *emptypb.Empty) (*movies_service.Genres, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetGenres")
	defer span.Finish()

	genres, err := s.repoManager.GetAllGenres(ctx)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	proto := make([]*movies_service.Genre, 0, len(genres))
	for _, genre := range genres {
		proto = append(proto, &movies_service.Genre{Id: genre.ID, Name: genre.Name})
	}
	return &movies_service.Genres{Genres: proto}, nil
}

func (s *MoviesService) GetCountries(ctx context.Context, in *emptypb.Empty) (*movies_service.Countries, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetCountries")
	defer span.Finish()

	countries, err := s.repoManager.GetAllCountries(ctx)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	proto := make([]*movies_service.Country, 0, len(countries))
	for _, country := range countries {
		proto = append(proto, &movies_service.Country{Id: country.ID, Name: country.Name})
	}

	span.SetTag("grpc.status", codes.OK)
	return &movies_service.Countries{Countries: proto}, nil
}

func GetAgeRatingsFilter(ageRating string) string {
	ageRating = ReplaceAll(strings.ReplaceAll(ageRating, " ", ""))
	str := strings.Split(ageRating, ",")
	for i := 0; i < len(str); i++ {
		if num, err := strconv.Atoi(str[i]); err == nil {
			str[i] = fmt.Sprintf("%d+", num)
		}
	}
	return strings.Join(str, ",")
}

func (s *MoviesService) convertDbMoviePreviewToProto(ctx context.Context, movie repository.MoviePreview) *movies_service.MoviePreview {
	return &movies_service.MoviePreview{
		TitleRU:   movie.TitleRU,
		TitleEN:   movie.TitleEN.String,
		Genres:    movie.Genres,
		Countries: movie.Countries,
		Duration:  movie.Duration,
		PreviewPosterURL: s.imagesService.GetPictureURL(ctx,
			movie.PreviewPosterID.String, s.picturesCfg.BaseUrl, s.picturesCfg.PreviewPostersCategory),
		ShortDescription: movie.ShortDescription,
		ReleaseYear:      movie.ReleaseYear,
		AgeRating:        movie.AgeRating,
	}
}

func (s *MoviesService) convertDbMovieToProto(ctx context.Context, movie repository.Movie) *movies_service.Movie {
	return &movies_service.Movie{
		Description: movie.Description,
		TitleRU:     movie.TitleRU,
		TitleEN:     movie.TitleEN.String,
		Genres:      movie.Genres,
		Duration:    movie.Duration,
		Countries:   movie.Countries,
		PosterURL: s.imagesService.GetPictureURL(ctx,
			movie.PosterID.String, s.picturesCfg.BaseUrl, s.picturesCfg.PostersCategory),
		BackgroundURL: s.imagesService.GetPictureURL(ctx,
			movie.BackgroundPictureID.String, s.picturesCfg.BaseUrl, s.picturesCfg.BackgroundsCategory),
		ReleaseYear: movie.ReleaseYear,
		AgeRating:   movie.AgeRating,
	}
}

func ReplaceAll(str string) string {
	return strings.ReplaceAll(str, `"`, "")
}
