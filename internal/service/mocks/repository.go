// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	models "github.com/Falokut/movies_service/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockMoviesRepository is a mock of MoviesRepository interface.
type MockMoviesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMoviesRepositoryMockRecorder
}

// MockMoviesRepositoryMockRecorder is the mock recorder for MockMoviesRepository.
type MockMoviesRepositoryMockRecorder struct {
	mock *MockMoviesRepository
}

// NewMockMoviesRepository creates a new mock instance.
func NewMockMoviesRepository(ctrl *gomock.Controller) *MockMoviesRepository {
	mock := &MockMoviesRepository{ctrl: ctrl}
	mock.recorder = &MockMoviesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMoviesRepository) EXPECT() *MockMoviesRepositoryMockRecorder {
	return m.recorder
}

// GetAgeRatings mocks base method.
func (m *MockMoviesRepository) GetAgeRatings(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgeRatings", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgeRatings indicates an expected call of GetAgeRatings.
func (mr *MockMoviesRepositoryMockRecorder) GetAgeRatings(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgeRatings", reflect.TypeOf((*MockMoviesRepository)(nil).GetAgeRatings), ctx)
}

// GetAllCountries mocks base method.
func (m *MockMoviesRepository) GetAllCountries(ctx context.Context) ([]models.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCountries", ctx)
	ret0, _ := ret[0].([]models.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCountries indicates an expected call of GetAllCountries.
func (mr *MockMoviesRepositoryMockRecorder) GetAllCountries(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCountries", reflect.TypeOf((*MockMoviesRepository)(nil).GetAllCountries), ctx)
}

// GetAllGenres mocks base method.
func (m *MockMoviesRepository) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllGenres", ctx)
	ret0, _ := ret[0].([]models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllGenres indicates an expected call of GetAllGenres.
func (mr *MockMoviesRepositoryMockRecorder) GetAllGenres(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllGenres", reflect.TypeOf((*MockMoviesRepository)(nil).GetAllGenres), ctx)
}

// GetCountries mocks base method.
func (m *MockMoviesRepository) GetCountries(ctx context.Context, movieId int32) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountries", ctx, movieId)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountries indicates an expected call of GetCountries.
func (mr *MockMoviesRepositoryMockRecorder) GetCountries(ctx, movieId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountries", reflect.TypeOf((*MockMoviesRepository)(nil).GetCountries), ctx, movieId)
}

// GetGenres mocks base method.
func (m *MockMoviesRepository) GetGenres(ctx context.Context, movieId int32) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres", ctx, movieId)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockMoviesRepositoryMockRecorder) GetGenres(ctx, movieId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockMoviesRepository)(nil).GetGenres), ctx, movieId)
}

// GetMovie mocks base method.
func (m *MockMoviesRepository) GetMovie(ctx context.Context, movieId int32) (models.RepositoryMovie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovie", ctx, movieId)
	ret0, _ := ret[0].(models.RepositoryMovie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovie indicates an expected call of GetMovie.
func (mr *MockMoviesRepositoryMockRecorder) GetMovie(ctx, movieId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovie", reflect.TypeOf((*MockMoviesRepository)(nil).GetMovie), ctx, movieId)
}

// GetMoviePreview mocks base method.
func (m *MockMoviesRepository) GetMoviePreview(ctx context.Context, movieId int32) (models.RepositoryMoviePreview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviePreview", ctx, movieId)
	ret0, _ := ret[0].(models.RepositoryMoviePreview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviePreview indicates an expected call of GetMoviePreview.
func (mr *MockMoviesRepositoryMockRecorder) GetMoviePreview(ctx, movieId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviePreview", reflect.TypeOf((*MockMoviesRepository)(nil).GetMoviePreview), ctx, movieId)
}

// GetMoviesPreview mocks base method.
func (m *MockMoviesRepository) GetMoviesPreview(ctx context.Context, Filter models.MoviesFilter, limit, offset uint32) ([]models.RepositoryMoviePreview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviesPreview", ctx, Filter, limit, offset)
	ret0, _ := ret[0].([]models.RepositoryMoviePreview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviesPreview indicates an expected call of GetMoviesPreview.
func (mr *MockMoviesRepositoryMockRecorder) GetMoviesPreview(ctx, Filter, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviesPreview", reflect.TypeOf((*MockMoviesRepository)(nil).GetMoviesPreview), ctx, Filter, limit, offset)
}

// GetMoviesPreviewByIDs mocks base method.
func (m *MockMoviesRepository) GetMoviesPreviewByIDs(ctx context.Context, ids []string) ([]models.RepositoryMoviePreview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviesPreviewByIDs", ctx, ids)
	ret0, _ := ret[0].([]models.RepositoryMoviePreview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviesPreviewByIDs indicates an expected call of GetMoviesPreviewByIDs.
func (mr *MockMoviesRepositoryMockRecorder) GetMoviesPreviewByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviesPreviewByIDs", reflect.TypeOf((*MockMoviesRepository)(nil).GetMoviesPreviewByIDs), ctx, ids)
}

// MockMoviesService is a mock of MoviesService interface.
type MockMoviesService struct {
	ctrl     *gomock.Controller
	recorder *MockMoviesServiceMockRecorder
}

// MockMoviesServiceMockRecorder is the mock recorder for MockMoviesService.
type MockMoviesServiceMockRecorder struct {
	mock *MockMoviesService
}

// NewMockMoviesService creates a new mock instance.
func NewMockMoviesService(ctrl *gomock.Controller) *MockMoviesService {
	mock := &MockMoviesService{ctrl: ctrl}
	mock.recorder = &MockMoviesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMoviesService) EXPECT() *MockMoviesServiceMockRecorder {
	return m.recorder
}

// GetAgeRatings mocks base method.
func (m *MockMoviesService) GetAgeRatings(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgeRatings", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgeRatings indicates an expected call of GetAgeRatings.
func (mr *MockMoviesServiceMockRecorder) GetAgeRatings(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgeRatings", reflect.TypeOf((*MockMoviesService)(nil).GetAgeRatings), ctx)
}

// GetCountries mocks base method.
func (m *MockMoviesService) GetCountries(ctx context.Context) ([]models.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountries", ctx)
	ret0, _ := ret[0].([]models.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountries indicates an expected call of GetCountries.
func (mr *MockMoviesServiceMockRecorder) GetCountries(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountries", reflect.TypeOf((*MockMoviesService)(nil).GetCountries), ctx)
}

// GetGenres mocks base method.
func (m *MockMoviesService) GetGenres(ctx context.Context) ([]models.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres", ctx)
	ret0, _ := ret[0].([]models.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockMoviesServiceMockRecorder) GetGenres(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockMoviesService)(nil).GetGenres), ctx)
}

// GetMovie mocks base method.
func (m *MockMoviesService) GetMovie(ctx context.Context, id int32) (models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMovie", ctx, id)
	ret0, _ := ret[0].(models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMovie indicates an expected call of GetMovie.
func (mr *MockMoviesServiceMockRecorder) GetMovie(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMovie", reflect.TypeOf((*MockMoviesService)(nil).GetMovie), ctx, id)
}

// GetMoviesPreview mocks base method.
func (m *MockMoviesService) GetMoviesPreview(ctx context.Context, filter models.MoviesFilter, limit, offset uint32) ([]models.MoviePreview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviesPreview", ctx, filter, limit, offset)
	ret0, _ := ret[0].([]models.MoviePreview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviesPreview indicates an expected call of GetMoviesPreview.
func (mr *MockMoviesServiceMockRecorder) GetMoviesPreview(ctx, filter, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviesPreview", reflect.TypeOf((*MockMoviesService)(nil).GetMoviesPreview), ctx, filter, limit, offset)
}

// GetMoviesPreviewByIDs mocks base method.
func (m *MockMoviesService) GetMoviesPreviewByIDs(ctx context.Context, ids []string) ([]models.MoviePreview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMoviesPreviewByIDs", ctx, ids)
	ret0, _ := ret[0].([]models.MoviePreview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMoviesPreviewByIDs indicates an expected call of GetMoviesPreviewByIDs.
func (mr *MockMoviesServiceMockRecorder) GetMoviesPreviewByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMoviesPreviewByIDs", reflect.TypeOf((*MockMoviesService)(nil).GetMoviesPreviewByIDs), ctx, ids)
}
