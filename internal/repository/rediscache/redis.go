package rediscache

import (
	"context"
	"errors"

	"github.com/Falokut/movies_service/internal/models"
	"github.com/redis/go-redis/v9"
)

func handleError(ctx context.Context, err *error) {
	if ctx.Err() != nil {
		var code models.ErrorCode
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			code = models.Canceled
		case errors.Is(ctx.Err(), context.DeadlineExceeded):
			code = models.DeadlineExceeded
		}
		*err = models.Error(code, ctx.Err().Error())
		return
	}

	if err == nil || *err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if !errors.As(*err, &repoErr) {
		var code models.ErrorCode
		switch {
		case errors.Is(*err, redis.Nil):
			code = models.NotFound
			*err = models.Error(code, "entity not found")
		default:
			code = models.Internal
			*err = models.Error(code, "cache internal error")
		}
	}
}

type Metrics interface {
	IncCacheHits(method string, times int32)
	IncCacheMiss(method string, times int32)
}
