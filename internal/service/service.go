package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Falokut/movies_service/internal/repository"
	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
	"github.com/mennanov/fmutils"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go
type ImagesService interface {
	GetPictureURL(pictureID, baseUrl, category string) string
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

	if in.Mask != nil && !in.Mask.IsValid(&movies_service.Movie{}) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "invalid mask value")
	}

	movie, err := s.repoManager.GetMovie(ctx, in.MovieID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	if in.Mask == nil || len(in.Mask.Paths) == 0 {
		return s.convertDbMovieToProto(movie), nil
	}

	res := s.convertDbMovieToProto(movie)
	fmutils.Filter(res, in.Mask.Paths)
	return res, nil
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

	if in.Mask != nil && !in.Mask.IsValid(&movies_service.MoviePreview{}) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "invalid mask value")
	}

	err := validateFilter(in)
	if errors.Is(err, ErrInvalidFilter) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
	}

	filter := repository.MoviesFilter{
		MoviesIDs:    ReplaceAllDoubleQuotes(in.GetMoviesIDs()),
		GenresIDs:    ReplaceAllDoubleQuotes(in.GetGenresIDs()),
		CountriesIDs: ReplaceAllDoubleQuotes(in.GetCountriesIDs()),
		AgeRating:    GetAgeRatingsFilter(in.GetAgeRatings()),
		Title:        ReplaceAllDoubleQuotes(in.GetTitle()),
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
	isMaskNil := in.Mask == nil || len(in.Mask.Paths) == 0
	for _, movie := range Movies {
		if isMaskNil {
			movies[movie.ID] = s.convertDbMoviePreviewToProto(movie)
		} else {
			filteredMovie := s.convertDbMoviePreviewToProto(movie)
			fmutils.Filter(filteredMovie, in.Mask.Paths)
			movies[movie.ID] = filteredMovie
		}
	}

	span.SetTag("grpc.status", codes.OK)
	return &movies_service.MoviesPreview{Movies: movies}, nil
}

func (s *MoviesService) GetMoviesPreviewByIDs(ctx context.Context,
	in *movies_service.GetMoviesPreviewByIDsRequest) (*movies_service.MoviesPreview, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetMoviesPreviewByIDs")
	defer span.Finish()

	if in.Mask != nil && !in.Mask.IsValid(&movies_service.MoviePreview{}) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "invalid mask value")
	}
	if in.MoviesIDs == "" {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "movies_ids mustn't be empty")
	}
	err := checkFilterParam(in.MoviesIDs)
	if errors.Is(err, ErrInvalidFilter) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, err.Error())
	}

	Movies, err := s.repoManager.GetMoviesPreviewByIDs(ctx,
		strings.Split(ReplaceAllDoubleQuotes(in.MoviesIDs), ","))
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
	isMaskNil := in.Mask == nil || len(in.Mask.Paths) == 0
	for _, movie := range Movies {
		if isMaskNil {
			movies[movie.ID] = s.convertDbMoviePreviewToProto(movie)
		} else {
			filteredMovie := s.convertDbMoviePreviewToProto(movie)
			fmutils.Filter(filteredMovie, in.Mask.Paths)
			movies[movie.ID] = filteredMovie
		}
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
	ageRating = ReplaceAllDoubleQuotes(strings.ReplaceAll(ageRating, " ", ""))
	str := strings.Split(ageRating, ",")
	for i := 0; i < len(str); i++ {
		if num, err := strconv.Atoi(str[i]); err == nil {
			str[i] = fmt.Sprintf("%d+", num)
		}
	}
	return strings.Join(str, ",")
}

func (s *MoviesService) convertDbMoviePreviewToProto(movie repository.MoviePreview) *movies_service.MoviePreview {
	posterURL := s.imagesService.GetPictureURL(movie.PreviewPosterID.String,
		s.picturesCfg.BaseUrl, s.picturesCfg.PreviewPostersCategory)
	return &movies_service.MoviePreview{
		TitleRu:          movie.TitleRU,
		TitleEn:          movie.TitleEN.String,
		Genres:           movie.Genres,
		Countries:        movie.Countries,
		Duration:         movie.Duration,
		PreviewPosterUrl: posterURL,
		ShortDescription: movie.ShortDescription,
		ReleaseYear:      movie.ReleaseYear,
		AgeRating:        movie.AgeRating,
	}
}

func (s *MoviesService) convertDbMovieToProto(movie repository.Movie) *movies_service.Movie {
	return &movies_service.Movie{
		Description: movie.Description,
		TitleRu:     movie.TitleRU,
		TitleEn:     movie.TitleEN.String,
		Genres:      movie.Genres,
		Duration:    movie.Duration,
		Countries:   movie.Countries,
		PosterUrl: s.imagesService.GetPictureURL(movie.PosterID.String,
			s.picturesCfg.BaseUrl, s.picturesCfg.PostersCategory),
		BackgroundUrl: s.imagesService.GetPictureURL(movie.BackgroundPictureID.String,
			s.picturesCfg.BaseUrl, s.picturesCfg.BackgroundsCategory),
		ReleaseYear: movie.ReleaseYear,
		AgeRating:   movie.AgeRating,
	}
}

func ReplaceAllDoubleQuotes(str string) string {
	return strings.ReplaceAll(str, `"`, "")
}
