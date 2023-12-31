// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockImagesService is a mock of ImagesService interface.
type MockImagesService struct {
	ctrl     *gomock.Controller
	recorder *MockImagesServiceMockRecorder
}

// MockImagesServiceMockRecorder is the mock recorder for MockImagesService.
type MockImagesServiceMockRecorder struct {
	mock *MockImagesService
}

// NewMockImagesService creates a new mock instance.
func NewMockImagesService(ctrl *gomock.Controller) *MockImagesService {
	mock := &MockImagesService{ctrl: ctrl}
	mock.recorder = &MockImagesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImagesService) EXPECT() *MockImagesServiceMockRecorder {
	return m.recorder
}

// GetPictureURL mocks base method.
func (m *MockImagesService) GetPictureURL(pictureID, baseUrl, category string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPictureURL", pictureID, baseUrl, category)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetPictureURL indicates an expected call of GetPictureURL.
func (mr *MockImagesServiceMockRecorder) GetPictureURL(pictureID, baseUrl, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPictureURL", reflect.TypeOf((*MockImagesService)(nil).GetPictureURL), pictureID, baseUrl, category)
}
