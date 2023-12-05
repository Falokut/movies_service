package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/Falokut/grpc_errors"
	"github.com/Falokut/movies_service/internal/repository"
	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=service.go -destination=mocks/imagesService.go
type ImagesService interface {
	GetMoviePosterURL(ctx context.Context, PictureID string) string
}

type MoviesService struct {
	movies_service.UnimplementedMoviesServiceV1Server
	logger        *logrus.Logger
	repoManager   repository.MoviesRepositoryManager
	errorHandler  errorHandler
	imagesService ImagesService
}

func NewMoviesService(logger *logrus.Logger, repoManager repository.MoviesRepositoryManager,
	imagesService ImagesService) *MoviesService {
	errorHandler := newErrorHandler(logger)
	return &MoviesService{
		logger:        logger,
		repoManager:   repoManager,
		errorHandler:  errorHandler,
		imagesService: imagesService,
	}
}

func (s *MoviesService) GetMovie(ctx context.Context, in *movies_service.GetMovieRequest) (*movies_service.Movie, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetMovie")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
	movie, err := s.repoManager.GetMovie(ctx, in.MovieID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponce(ErrNotFound, "")
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	return s.convertDbMoviesToProto(ctx, movie), nil
}

func (s *MoviesService) GetMovies(ctx context.Context, in *movies_service.GetMoviesRequest) (*movies_service.Movies, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MoviesService.GetMovies")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	err = validateFilter(in)
	if errors.Is(err, ErrInvalidFilter) {
		return nil, s.errorHandler.createErrorResponce(ErrInvalidArgument, err.Error())
	}

	filter := repository.MoviesFilter{
		MoviesIDs:    in.GetMoviesIDs(),
		GenresIDs:    in.GetGenresIDs(),
		DirectorsIDs: in.GetDirectorsIDs(),
		CountriesIDs: in.GetCountriesIDs(),
		Title:        getTitle(in.GetTitle()),
	}

	if in.Limit == 0 {
		in.Limit = 10
	} else if in.Limit > 100 {
		in.Limit = 100
	}

	Movies, err := s.repoManager.GetMovies(ctx, filter, in.Limit, in.Offset)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponce(ErrNotFound, "")
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	if len(Movies) == 0 {
		return nil, s.errorHandler.createErrorResponce(ErrNotFound, "")
	}

	movies := make(map[string]*movies_service.Movie, len(Movies))
	for _, movie := range Movies {
		movies[movie.ID] = s.convertDbMoviesToProto(ctx, movie)
	}

	return &movies_service.Movies{Movies: movies}, nil
}

func (s *MoviesService) convertDbMoviesToProto(ctx context.Context, movie repository.Movie) *movies_service.Movie {
	return &movies_service.Movie{
		MovieID:          movie.ID,
		TitleRU:          movie.TitleRU,
		TitleEN:          movie.TitleEN.String,
		Budget:           movie.Budget.String,
		CastID:           movie.CastID,
		GenresIDs:        s.intSliceFromString(movie.Genres.String),
		DirectorsIDs:     s.intSliceFromString(movie.DirectorsIDs.String),
		CountriesIDs:     s.intSliceFromString(movie.CountriesIDs.String),
		Duration:         movie.Duration,
		PosterPictureURL: s.imagesService.GetMoviePosterURL(ctx, movie.PictureID.String),
		Plot:             movie.Plot,
		ReleaseYear:      movie.ReleaseYear,
	}
}

func getTitle(title string) string {
	return strings.ReplaceAll(title, `"`, "")
}

func (s *MoviesService) intSliceFromString(str string) []int32 {
	if str == "" {
		return []int32{}
	}
	str = strings.Trim(str, "{}")
	nums := strings.Split(str, ",")
	var res = make([]int32, len(nums))

	for i, n := range nums {
		num, err := strconv.Atoi(n)
		if err != nil {
			s.logger.Errorf("invalid string, can't convert into int slice err: %s, string %s", err, str)
			return []int32{}
		}
		res[i] = int32(num)
	}
	return res
}
