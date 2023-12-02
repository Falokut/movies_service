package service

import (
	"context"
	"net/url"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type ImageServiceConfig struct {
	BasePosterPictureUrl  string
	PosterPictureCategory string
}

type imageService struct {
	cfg          ImageServiceConfig
	logger       *logrus.Logger
	errorHandler errorHandler
}

func NewImageService(cfg ImageServiceConfig, logger *logrus.Logger) *imageService {
	errorHandler := newErrorHandler(logger)
	return &imageService{
		cfg:          cfg,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// Returns movie poster url for GET request
func (s *imageService) GetMoviePosterURL(ctx context.Context, PictureID string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "imageService.GetMoviePosterURL")
	defer span.Finish()

	if PictureID == "" {
		return ""
	}

	u, err := url.Parse(s.cfg.BasePosterPictureUrl)
	if err != nil {
		s.logger.Errorf("can't parse url. error: %s", err.Error())
		return ""
	}

	q := u.Query()
	q.Add("image_id", PictureID)
	q.Add("category", s.cfg.PosterPictureCategory)
	u.RawQuery = q.Encode()
	return u.String()
}
