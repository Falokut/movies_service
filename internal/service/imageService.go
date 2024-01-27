package service

import (
	"github.com/sirupsen/logrus"
)

type imageService struct {
	logger       *logrus.Logger
	errorHandler errorHandler
}

func NewImageService(logger *logrus.Logger) *imageService {
	errorHandler := newErrorHandler(logger)
	return &imageService{
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// Returns picture url for GET request
func (s *imageService) GetPictureURL(pictureID, baseUrl, category string) string {
	if pictureID == "" || baseUrl == "" || category == "" {
		return ""
	}

	return baseUrl + "/" + category + "/" + pictureID
}
