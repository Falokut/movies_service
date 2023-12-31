package service

import (
	"context"
	"net/url"

	"github.com/opentracing/opentracing-go"
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
func (s *imageService) GetPictureURL(ctx context.Context, pictureID, baseUrl, category string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "imageService.GetPictureURL")
	defer span.Finish()

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
