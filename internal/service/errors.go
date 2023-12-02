package service

import (
	"errors"
	"fmt"

	"github.com/Falokut/grpc_errors"
	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInternal        = errors.New("internal error")
	ErrInvalidArgument = errors.New("invalid input data")
)

var errorCodes = map[error]codes.Code{
	ErrNotFound:        codes.NotFound,
	ErrInvalidArgument: codes.InvalidArgument,
	ErrInternal:        codes.Internal,
}

type errorHandler struct {
	logger *logrus.Logger
}

func newErrorHandler(logger *logrus.Logger) errorHandler {
	return errorHandler{
		logger: logger,
	}
}

func (e *errorHandler) createErrorResponce(err error, developerMessage string) error {
	var msg string
	if len(developerMessage) == 0 {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s. error: %v", developerMessage, err)
	}

	err = status.Error(grpc_errors.GetGrpcCode(err), msg)
	e.logger.Error(err)
	return err
}

func (e *errorHandler) createExtendedErrorResponce(err error, DeveloperMessage, UserMessage string) error {
	var msg string
	if DeveloperMessage == "" {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s. error: %v", DeveloperMessage, err)
	}

	extErr := status.New(grpc_errors.GetGrpcCode(err), msg)
	if len(UserMessage) > 0 {
		extErr, _ = extErr.WithDetails(&movies_service.UserErrorMessage{Message: UserMessage})
		if extErr == nil {
			e.logger.Error(err)
			return err
		}
	}

	e.logger.Error(extErr)
	return extErr.Err()
}

func init() {
	grpc_errors.RegisterErrors(errorCodes)
}
