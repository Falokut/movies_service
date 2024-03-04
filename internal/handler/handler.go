package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Falokut/movies_service/internal/models"
	"github.com/Falokut/movies_service/internal/service"
	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
	"github.com/mennanov/fmutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MoviesServiceHandler struct {
	movies_service.UnimplementedMoviesServiceV1Server
	service service.MoviesService
}

func NewMoviesServiceHandler(service service.MoviesService) *MoviesServiceHandler {
	return &MoviesServiceHandler{service: service}
}

func (h *MoviesServiceHandler) GetAgeRatings(ctx context.Context,
	in *emptypb.Empty) (res *movies_service.AgeRatings, err error) {
	defer h.handleError(&err)

	ratings, err := h.service.GetAgeRatings(ctx)
	if err != nil {
		return
	}
	return &movies_service.AgeRatings{
		Ratings: ratings,
	}, nil
}

func (h *MoviesServiceHandler) GetGenres(ctx context.Context,
	in *emptypb.Empty) (res *movies_service.Genres, err error) {
	defer h.handleError(&err)

	genres, err := h.service.GetGenres(ctx)
	if err != nil {
		return
	}

	res = &movies_service.Genres{}

	res.Genres = make([]*movies_service.Genre, len(genres))
	for i := range genres {
		res.Genres[i] = &movies_service.Genre{
			Id:   genres[i].ID,
			Name: genres[i].Name,
		}
	}
	return
}

func (h *MoviesServiceHandler) GetCountries(ctx context.Context,
	in *emptypb.Empty) (res *movies_service.Countries, err error) {
	defer h.handleError(&err)

	countries, err := h.service.GetCountries(ctx)
	if err != nil {
		return
	}

	res = &movies_service.Countries{}

	res.Countries = make([]*movies_service.Country, len(countries))
	for i := range countries {
		res.Countries[i] = &movies_service.Country{
			Id:   countries[i].ID,
			Name: countries[i].Name,
		}
	}
	return
}

func (h *MoviesServiceHandler) GetMovie(ctx context.Context,
	in *movies_service.GetMovieRequest) (res *movies_service.Movie, err error) {
	defer h.handleError(&err)

	if in.Mask != nil && !in.Mask.IsValid(&movies_service.Movie{}) {
		return nil, status.Error(codes.InvalidArgument, "invalid mask value")
	}

	movie, err := h.service.GetMovie(ctx, in.MovieID)
	if err != nil {
		return
	}

	res = &movies_service.Movie{
		Description:   movie.Description,
		TitleRu:       movie.TitleRU,
		TitleEn:       movie.TitleEN,
		Genres:        movie.Genres,
		Duration:      movie.Duration,
		Countries:     movie.Countries,
		PosterUrl:     movie.PosterUrl,
		BackgroundUrl: movie.BackgroundUrl,
		ReleaseYear:   movie.ReleaseYear,
		AgeRating:     movie.AgeRating,
	}

	if in.Mask == nil || len(in.Mask.Paths) == 0 {
		return
	}

	fmutils.Filter(res, in.Mask.Paths)
	return
}

func (h *MoviesServiceHandler) GetMoviesPreview(ctx context.Context,
	in *movies_service.GetMoviesPreviewRequest) (res *movies_service.MoviesPreview, err error) {
	defer h.handleError(&err)

	if in.Mask != nil && !in.Mask.IsValid(&movies_service.MoviePreview{}) {
		err = status.Error(codes.InvalidArgument, "invalid mask value")
		return
	}

	err = validateFilter(in)
	if err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	filter := models.MoviesFilter{
		MoviesIDs:    ReplaceAllDoubleQuotes(in.GetMoviesIDs()),
		GenresIDs:    ReplaceAllDoubleQuotes(in.GetGenresIDs()),
		CountriesIDs: ReplaceAllDoubleQuotes(in.GetCountriesIDs()),
		AgeRating:    GetAgeRatingsFilter(in.GetAgeRatings()),
		Title:        ReplaceAllDoubleQuotes(in.GetTitle()),
	}
	if in.Limit == 0 {
		in.Limit = 10
	} else if in.Limit > 100 {
		in.Limit = 100
	}

	movies, err := h.service.GetMoviesPreview(ctx, filter, in.Limit, in.Offset)
	if err != nil {
		return
	}

	paths := []string{}
	if in.Mask != nil {
		paths = in.Mask.Paths
	}
	res = &movies_service.MoviesPreview{
		Movies: convertModelsMoviesPreviewToProto(movies, paths),
	}
	return
}

func (h *MoviesServiceHandler) GetMoviesPreviewByIDs(ctx context.Context,
	in *movies_service.GetMoviesPreviewByIDsRequest) (res *movies_service.MoviesPreview, err error) {
	defer h.handleError(&err)

	if in.Mask != nil && !in.Mask.IsValid(&movies_service.MoviePreview{}) {
		err = status.Error(codes.InvalidArgument, "invalid mask value")
		return
	}
	err = checkFilterParam(in.MoviesIDs)
	if err != nil {
		err = status.Error(codes.InvalidArgument, "movies_ids mustn't be empty and must contain only digits and commas")
		return
	}
	ids := strings.Split(ReplaceAllDoubleQuotes(in.MoviesIDs), ",")
	movies, err := h.service.GetMoviesPreviewByIDs(ctx, ids)
	if err != nil {
		return
	}
	paths := []string{}
	if in.Mask != nil {
		paths = in.Mask.Paths
	}

	res = &movies_service.MoviesPreview{
		Movies: convertModelsMoviesPreviewToProto(movies, paths),
	}
	return
}

func convertModelsMoviesPreviewToProto(movies []models.MoviePreview, maskPaths []string) map[int32]*movies_service.MoviePreview {
	res := make(map[int32]*movies_service.MoviePreview, len(movies))
	isMaskNil := len(maskPaths) == 0
	for i := range movies {
		movie := &movies_service.MoviePreview{
			TitleRu:          movies[i].TitleRU,
			TitleEn:          movies[i].TitleEN,
			Genres:           movies[i].Genres,
			Countries:        movies[i].Countries,
			Duration:         movies[i].Duration,
			PreviewPosterUrl: movies[i].PreviewPosterUrl,
			ShortDescription: movies[i].ShortDescription,
			ReleaseYear:      movies[i].ReleaseYear,
			AgeRating:        movies[i].AgeRating,
		}
		if !isMaskNil {
			fmutils.Filter(movie, maskPaths)
		}

		res[movies[i].ID] = movie
	}

	return res
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

func ReplaceAllDoubleQuotes(str string) string {
	return strings.ReplaceAll(str, `"`, "")
}

func (h *MoviesServiceHandler) handleError(err *error) {
	if err == nil || *err == nil {
		return
	}

	serviceErr := &models.ServiceError{}
	if errors.As(*err, &serviceErr) {
		*err = status.Error(convertServiceErrCodeToGrpc(serviceErr.Code), serviceErr.Msg)
	} else if _, ok := status.FromError(*err); !ok {
		e := *err
		*err = status.Error(codes.Unknown, e.Error())
	}
}

func convertServiceErrCodeToGrpc(code models.ErrorCode) codes.Code {
	switch code {
	case models.Internal:
		return codes.Internal
	case models.InvalidArgument:
		return codes.InvalidArgument
	case models.Unauthenticated:
		return codes.Unauthenticated
	case models.Conflict:
		return codes.AlreadyExists
	case models.NotFound:
		return codes.NotFound
	case models.Canceled:
		return codes.Canceled
	case models.DeadlineExceeded:
		return codes.DeadlineExceeded
	case models.PermissionDenied:
		return codes.PermissionDenied
	default:
		return codes.Unknown
	}
}
