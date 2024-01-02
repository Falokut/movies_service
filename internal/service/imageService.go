package service

import (
	"net/url"

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

	u, err := url.Parse(baseUrl)
	if err != nil {
		s.logger.Errorf("can't parse url. error: %s", err.Error())
		return ""
	}

	q := u.Query()
	q.Add("image_id", pictureID)
	q.Add("category", category)
	u.RawQuery = q.Encode()
	return u.String()
}
