package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Falokut/movies_service/internal/repository"
	repo_mock "github.com/Falokut/movies_service/internal/repository/mocks"
	"github.com/Falokut/movies_service/internal/service"
	service_mock "github.com/Falokut/movies_service/internal/service/mocks"
	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
	gomock "github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

func getNullLogger() *logrus.Logger {
	l, _ := test.NewNullLogger()
	return l
}

func newServer(t *testing.T, register func(srv *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	return conn
}

func newClient(t *testing.T, s *service.MoviesService) *grpc.ClientConn {
	return newServer(t, func(srv *grpc.Server) { movies_service.RegisterMoviesServiceV1Server(srv, s) })
}

type GetMovieBehavior func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, movieId int32, expectedMovie repository.Movie)
type GetMoviesPreviewBehavior func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
	expectedMovies []repository.MoviePreview, limit, offset uint32)
type GetPictureURLMultipleBehavior func(s *service_mock.MockImagesService, PicturesIDs []string, times int)

func TestGetMovie(t *testing.T) {
	testCases := []struct {
		movieID          int32
		behavior         GetMovieBehavior
		imgBehavior      GetPictureURLMultipleBehavior
		expectedStatus   codes.Code
		expectedResponce *movies_service.Movie
		movie            repository.Movie
		expectedError    error
		urlRequestTimes  int
		msg              string
	}{
		{
			movieID: 1,
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, movieId int32, expectedMovie repository.Movie) {
				m.EXPECT().GetMovie(gomock.Any(), movieId).Return(repository.Movie{}, repository.ErrNotFound).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Times(times)
			},
			expectedStatus:   codes.NotFound,
			expectedError:    service.ErrNotFound,
			expectedResponce: nil,
			msg:              "Test case num %d, must return not found error, if movie not found",
		},
		{
			movieID: 1,
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, movieId int32, expectedMovie repository.Movie) {
				m.EXPECT().GetMovie(gomock.Any(), movieId).Return(repository.Movie{}, context.Canceled).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Times(times)
			},
			expectedStatus:   codes.Internal,
			expectedError:    service.ErrInternal,
			expectedResponce: nil,
			msg:              "Test case num %d, must return internal error, if repo manager return error != ErrNotFound",
		},
		{
			movieID: 10,
			movie: repository.Movie{
				ID:                  10,
				TitleRU:             "TitleRU",
				PosterID:            sql.NullString{String: "1012", Valid: true},
				BackgroundPictureID: sql.NullString{String: "1012", Valid: true},
				Description:         "Description",
				Duration:            100,
				ReleaseYear:         2000,
			},
			expectedResponce: &movies_service.Movie{
				TitleRu:       "TitleRU",
				PosterUrl:     "someurl",
				BackgroundUrl: "someurl",
				Description:   "Description",
				Duration:      100,
				ReleaseYear:   2000,
			},
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, movieId int32, expectedMovie repository.Movie) {
				m.EXPECT().GetMovie(gomock.Any(), movieId).Return(expectedMovie, nil).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Return("someurl").Times(times)
			},
			urlRequestTimes: 2,
			expectedStatus:  codes.OK,
			msg:             "Test case num %d, must return expected responce, if repo manager doesn't return error, service shouldn't change data, except for the link to the poster",
		},
		{
			movieID: 10,
			movie: repository.Movie{
				ID:                  10,
				TitleRU:             "TitleRU",
				Genres:              []string{"genre1", "genre2"},
				PosterID:            sql.NullString{String: "1012", Valid: true},
				BackgroundPictureID: sql.NullString{String: "someStr", Valid: true},
				Description:         "ShortDescription",
				Duration:            100,
				ReleaseYear:         2000,
			},
			expectedResponce: &movies_service.Movie{
				TitleRu:       "TitleRU",
				PosterUrl:     "someurl",
				BackgroundUrl: "someurl",
				Description:   "ShortDescription",
				Genres:        []string{"genre1", "genre2"},
				Duration:      100,
				ReleaseYear:   2000,
			},
			urlRequestTimes: 2,
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, movieId int32, expectedMovie repository.Movie) {
				m.EXPECT().GetMovie(gomock.Any(), movieId).Return(expectedMovie, nil).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Return("someurl").Times(times)
			},
			expectedStatus: codes.OK,
			msg:            "Test case num %d, must return expected responce, if repo manager doesn't return error, service shouldn't change data, except for the link to the poster",
		},
	}

	for i, testCase := range testCases {
		ctrl := gomock.NewController(t)
		repo := repo_mock.NewMockMoviesRepositoryManager(ctrl)
		imgServ := service_mock.NewMockImagesService(ctrl)
		testCase.imgBehavior(imgServ, []string{testCase.movie.PosterID.String, testCase.movie.BackgroundPictureID.String}, testCase.urlRequestTimes)
		testCase.behavior(repo, context.Background(), testCase.movieID, testCase.movie)

		conn := newClient(t, service.NewMoviesService(getNullLogger(), repo, imgServ, service.PicturesUrlConfig{}))
		defer conn.Close()

		client := movies_service.NewMoviesServiceV1Client(conn)
		assert.NotNil(t, client)

		res, err := client.GetMovie(context.Background(), &movies_service.GetMovieRequest{
			MovieID: testCase.movieID,
		})

		testCase.msg = fmt.Sprintf(testCase.msg, i+1)
		if testCase.expectedError != nil {
			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
			}
		} else if assert.NotNil(t, res) {
			var comp assert.Comparison = func() (success bool) {
				return isProtoMoviesEqual(t, testCase.expectedResponce, res)
			}
			assert.Condition(t, comp, testCase.msg)
		}
		assert.Equal(t, testCase.expectedStatus, status.Code(err), testCase.msg)
	}
}

func TestGetMoviesPreview(t *testing.T) {
	type MoviesRequest struct {
		MoviesIDs    string
		GenresIDs    string
		DiretorsIDs  string
		CountriesIDs string
		Title        string
		Limit        uint32
		Offset       uint32
	}

	testCases := []struct {
		moviesIDs        []string
		request          MoviesRequest
		urlRequestTimes  int
		behavior         GetMoviesPreviewBehavior
		imgBehavior      GetPictureURLMultipleBehavior
		expectedStatus   codes.Code
		expectedResponce *movies_service.MoviesPreview
		movies           []repository.MoviePreview
		expectedError    error
		msg              string
	}{
		{
			moviesIDs: []string{"1"},
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
				expectedMovies []repository.MoviePreview, limit, offset uint32) {
				m.EXPECT().GetMoviesPreview(gomock.Any(), filter, limit, offset).Return([]repository.MoviePreview{}, repository.ErrNotFound).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Times(times)
			},
			expectedStatus:   codes.NotFound,
			expectedError:    service.ErrNotFound,
			expectedResponce: nil,
			msg:              "Test case num %d, must return not found error, if movies not found",
		},
		{
			moviesIDs: []string{"1"},
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
				expectedMovies []repository.MoviePreview, limit, offset uint32) {
				m.EXPECT().GetMoviesPreview(gomock.Any(), filter, limit, offset).Return([]repository.MoviePreview{}, context.Canceled).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Times(times)
			},
			expectedStatus:   codes.Internal,
			expectedError:    service.ErrInternal,
			expectedResponce: nil,
			msg:              "Test case num %d, must return internal error, if repo manager return error != ErrNotFound",
		},
		{
			moviesIDs: []string{"10", "12"},
			request: MoviesRequest{
				MoviesIDs: "10,12",
			},
			movies: []repository.MoviePreview{
				{
					ID:               10,
					TitleRU:          "TitleRU",
					ShortDescription: "ShortDescription",
					Duration:         100,
					ReleaseYear:      2000,
				},
				{
					ID:               12,
					TitleRU:          "TitleRU",
					TitleEN:          sql.NullString{String: "TitleEn", Valid: true},
					PreviewPosterID:  sql.NullString{String: "validString", Valid: true},
					ShortDescription: "ShortDescription",
					Duration:         150,
					ReleaseYear:      2200,
				},
			},

			expectedResponce: &movies_service.MoviesPreview{
				Movies: map[int32]*movies_service.MoviePreview{
					10: {
						TitleRu:          "TitleRU",
						PreviewPosterUrl: "",
						ShortDescription: "ShortDescription",
						Duration:         100,
						ReleaseYear:      2000,
					},

					12: {
						TitleEn:          "TitleEn",
						TitleRu:          "TitleRU",
						PreviewPosterUrl: "",
						ShortDescription: "ShortDescription",
						Duration:         150,
						ReleaseYear:      2200,
					},
				},
			},
			urlRequestTimes: 2,
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
				expectedMovies []repository.MoviePreview, limit, offset uint32) {
				m.EXPECT().GetMoviesPreview(gomock.Any(), filter, limit, offset).Return(expectedMovies, nil).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Return("").Times(times)
			},
			expectedStatus: codes.OK,
			msg:            "Test case num %d, must return expected responce, if repo manager doesn't return error, service shouldn't change data, except for the link to the poster",
		},
		{
			moviesIDs: []string{"1"},
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
				expectedMovies []repository.MoviePreview, limit, offset uint32) {
				m.EXPECT().GetMoviesPreview(gomock.Any(), filter, limit, offset).Return([]repository.MoviePreview{}, context.Canceled).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Times(times)
			},
			expectedStatus:   codes.Internal,
			expectedError:    service.ErrInternal,
			expectedResponce: nil,
			msg:              "Test case num %d, must return internal error, if repo manager return error != ErrNotFound",
		},
		{
			moviesIDs: []string{"10", "12"},
			request: MoviesRequest{
				MoviesIDs:   "10,12",
				GenresIDs:   "10,2",
				DiretorsIDs: "1,9,11,99",
				Limit:       110,
			},

			movies:           []repository.MoviePreview{},
			expectedResponce: nil,
			urlRequestTimes:  0,
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
				expectedMovies []repository.MoviePreview, limit, offset uint32) {
				m.EXPECT().GetMoviesPreview(gomock.Any(), filter, limit, offset).Return(expectedMovies, nil).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Return("").Times(times)
			},
			expectedStatus: codes.NotFound,
			expectedError:  service.ErrNotFound,
			msg: "Test case num %d, must return expected error," +
				"if repo manager return empty Moves slice",
		},

		{
			moviesIDs: []string{"10", "120"},
			request: MoviesRequest{
				MoviesIDs: "10,120",
				Limit:     110,
			},
			movies: []repository.MoviePreview{
				{
					ID:               120,
					TitleRU:          "TitleRU",
					ShortDescription: "ShortDescription",
					Duration:         100,
					ReleaseYear:      2000,
				},
				{
					ID:               12,
					TitleRU:          "TitleRU",
					TitleEN:          sql.NullString{String: "TitleEn", Valid: true},
					ShortDescription: "ShortDescription",
					Duration:         150,
					ReleaseYear:      2200,
				},
			},

			expectedResponce: &movies_service.MoviesPreview{
				Movies: map[int32]*movies_service.MoviePreview{
					120: {
						TitleRu:          "TitleRU",
						ShortDescription: "ShortDescription",
						Duration:         100,
						ReleaseYear:      2000,
					},

					12: {
						TitleEn:          "TitleEn",
						TitleRu:          "TitleRU",
						ShortDescription: "ShortDescription",
						Duration:         150,
						ReleaseYear:      2200,
					},
				},
			},
			urlRequestTimes: 2,
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
				expectedMovies []repository.MoviePreview, limit, offset uint32) {
				m.EXPECT().GetMoviesPreview(gomock.Any(), filter, limit, offset).Return(expectedMovies, nil).Times(1)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Return("").Times(times)
			},
			expectedStatus: codes.OK,
			msg: "Test case num %d, must return expected responce, limit should be in [10,100]," +
				"if repo manager doesn't return error, service shouldn't change data, except for the link to the poster",
		},
		{
			moviesIDs: []string{"10", "12"},
			request: MoviesRequest{
				MoviesIDs: "10-,2-12",
			},
			movies: []repository.MoviePreview{
				{
					ID:               10,
					TitleRU:          "TitleRU",
					ShortDescription: "ShortDescription",
					Duration:         100,
					ReleaseYear:      2000,
				},
				{
					ID:               12,
					TitleRU:          "TitleRU",
					TitleEN:          sql.NullString{String: "TitleEn", Valid: true},
					ShortDescription: "ShortDescription",
					Duration:         150,
					ReleaseYear:      2200,
				},
			},

			expectedResponce: &movies_service.MoviesPreview{
				Movies: map[int32]*movies_service.MoviePreview{
					10: {
						TitleRu:          "TitleRU",
						PreviewPosterUrl: "",
						ShortDescription: "ShortDescription",
						Duration:         100,
						ReleaseYear:      2000,
					},

					12: {
						TitleEn:          "TitleEn",
						PreviewPosterUrl: "TitleRU",
						ShortDescription: "ShortDescription",
						Duration:         150,
						ReleaseYear:      2200,
					},
				},
			},
			behavior: func(m *repo_mock.MockMoviesRepositoryManager, ctx context.Context, filter repository.MoviesFilter,
				expectedMovies []repository.MoviePreview, limit, offset uint32) {
				m.EXPECT().GetMoviesPreview(gomock.Any(), filter, limit, offset).Return(expectedMovies, nil).Times(0)
			},
			imgBehavior: func(s *service_mock.MockImagesService, PicturesIDs []string, times int) {
				s.EXPECT().GetPictureURL(gomock.Any(), gomock.Any(), gomock.Any()).Return("").Times(times)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrInvalidFilter,
			msg:            "Test case num %d, must return expected error, if filter not valid",
		},
	}

	for i, testCase := range testCases {
		ctrl := gomock.NewController(t)
		repo := repo_mock.NewMockMoviesRepositoryManager(ctrl)
		imgServ := service_mock.NewMockImagesService(ctrl)

		var picturesIds = make([]string, 0, len(testCase.movies))
		for _, movie := range testCase.movies {
			picturesIds = append(picturesIds, movie.PreviewPosterID.String)
		}
		testCase.imgBehavior(imgServ, picturesIds, testCase.urlRequestTimes)
		var limit = testCase.request.Limit
		if limit == 0 {
			limit = 10
		} else if limit > 100 {
			limit = 100
		}
		testCase.behavior(repo, context.Background(), repository.MoviesFilter{
			MoviesIDs:    testCase.request.MoviesIDs,
			GenresIDs:    testCase.request.GenresIDs,
			CountriesIDs: testCase.request.CountriesIDs,
			Title:        testCase.request.Title,
		}, testCase.movies,
			limit, testCase.request.Offset)

		conn := newClient(t, service.NewMoviesService(getNullLogger(), repo, imgServ, service.PicturesUrlConfig{}))
		defer conn.Close()

		client := movies_service.NewMoviesServiceV1Client(conn)
		assert.NotNil(t, client)

		res, err := client.GetMoviesPreview(context.Background(),
			&movies_service.GetMoviesPreviewRequest{
				MoviesIDs:    &testCase.request.MoviesIDs,
				GenresIDs:    &testCase.request.GenresIDs,
				CountriesIDs: &testCase.request.CountriesIDs,
				Title:        &testCase.request.Title,
				Limit:        testCase.request.Limit,
				Offset:       testCase.request.Offset,
			})

		testCase.msg = fmt.Sprintf(testCase.msg, i+1)
		if testCase.expectedError != nil {
			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), testCase.expectedError.Error())
			}
		} else if assert.NotNil(t, res) && assert.Equal(t, len(testCase.expectedResponce.Movies), len(res.Movies)) {
			var comp assert.Comparison = func() (success bool) {
				for key, Expectedmovie := range testCase.expectedResponce.Movies {
					if !isProtoMoviesPreviewEqual(t, Expectedmovie, res.Movies[key]) {
						return false
					}
				}
				return true
			}
			assert.Condition(t, comp, testCase.msg)
		}
		assert.Equal(t, testCase.expectedStatus, status.Code(err), testCase.msg)
	}
}

func TestGetAgeRatings(t *testing.T) {
	var GetAgeRatingsBehavior = func(s *repo_mock.MockMoviesRepositoryManager, Ratings []string, err error) {
		s.EXPECT().GetAgeRatings(gomock.Any()).Return(Ratings, err).Times(1)
	}

	testCases := []struct {
		RepoError      error
		ExpectError    error
		ExpectCode     codes.Code
		ExpectResponce *movies_service.AgeRatings
		msg            string
	}{{
		ExpectResponce: &movies_service.AgeRatings{Ratings: []string{"Rating1", "Rating2"}},
		ExpectCode:     codes.OK,
		msg:            "test case num %d, service must return expected responce",
	},
		{
			RepoError:   errors.New("new error"),
			ExpectError: service.ErrInternal,
			ExpectCode:  codes.Internal,
			msg:         "test case num %d, service must return ErrInternal, if repo returns error",
		},
	}

	for i, testCase := range testCases {
		ctrl := gomock.NewController(t)
		repo := repo_mock.NewMockMoviesRepositoryManager(ctrl)
		imgServ := service_mock.NewMockImagesService(ctrl)
		GetAgeRatingsBehavior(repo, testCase.ExpectResponce.GetRatings(), testCase.RepoError)

		conn := newClient(t, service.NewMoviesService(getNullLogger(), repo, imgServ, service.PicturesUrlConfig{}))
		defer conn.Close()

		client := movies_service.NewMoviesServiceV1Client(conn)
		assert.NotNil(t, client)

		msg := fmt.Sprintf(testCase.msg, i+1)
		res, err := client.GetAgeRatings(context.Background(), &emptypb.Empty{})
		if testCase.ExpectError != nil {
			assert.ErrorContains(t, err, testCase.ExpectError.Error(), msg)
		}
		assert.True(t, isProtoAgeRatingsEqual(t, testCase.ExpectResponce, res), msg)
		assert.Equal(t, testCase.ExpectCode, status.Code(err), msg)
	}
}

func isProtoAgeRatingsEqual(t *testing.T, expected, result *movies_service.AgeRatings) bool {
	if expected == nil && result == nil {
		return true
	} else if expected == nil && result != nil ||
		expected != nil && result == nil {
		return false
	}

	return assert.Equal(t, expected.Ratings, result.Ratings, "Ratings not equal")
}

func isProtoMoviesPreviewEqual(t *testing.T, expected, result *movies_service.MoviePreview) bool {
	if expected == nil && result == nil {
		return true
	} else if expected == nil && result != nil ||
		expected != nil && result == nil {
		return false
	}
	return assert.Equal(t, expected.ShortDescription, result.ShortDescription, "short descriptions not equals") &&
		assert.Equal(t, expected.TitleRu, result.TitleRu, "ru titles not equals") &&
		assert.Equal(t, expected.TitleEn, result.TitleEn, "en titles not equals") &&
		assert.Equal(t, expected.Genres, result.Genres, "genres ids not equals") &&
		assert.Equal(t, expected.Duration, result.Duration, "duration not equals") &&
		assert.Equal(t, expected.Countries, result.Countries, "countries ids not equals") &&
		assert.Equal(t, expected.PreviewPosterUrl, result.PreviewPosterUrl, "preview posters urls not equals") &&
		assert.Equal(t, expected.ReleaseYear, result.ReleaseYear, "release years not equals")
}

func isProtoMoviesEqual(t *testing.T, expected, result *movies_service.Movie) bool {
	if expected == nil && result == nil {
		return true
	} else if expected == nil && result != nil ||
		expected != nil && result == nil {
		return false
	}
	return assert.Equal(t, expected.Description, result.Description, "descriptions not equals") &&
		assert.Equal(t, expected.TitleRu, result.TitleRu, "ru titles not equals") &&
		assert.Equal(t, expected.TitleEn, result.TitleEn, "en titles not equals") &&
		assert.Equal(t, expected.Genres, result.Genres, "genres ids not equals") &&
		assert.Equal(t, expected.Duration, result.Duration, "duration not equals") &&
		assert.Equal(t, expected.Countries, result.Countries, "countries ids not equals") &&
		assert.Equal(t, expected.PosterUrl, result.PosterUrl, "posters urls not equals") &&
		assert.Equal(t, expected.BackgroundUrl, result.BackgroundUrl, "backgrounds not equals") &&
		assert.Equal(t, expected.ReleaseYear, result.ReleaseYear, "release years not equals")
}
