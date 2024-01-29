package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"
	"github.com/Falokut/movies_service/internal/config"
	"github.com/Falokut/movies_service/internal/repository"
	"github.com/Falokut/movies_service/internal/service"
	jaegerTracer "github.com/Falokut/movies_service/pkg/jaeger"
	"github.com/Falokut/movies_service/pkg/metrics"
	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	logging.NewEntry(logging.ConsoleOutput)
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Error(err)
	}
	logger.Logger.SetLevel(logLevel)

	tracer, closer, err := jaegerTracer.InitJaeger(cfg.JaegerConfig)
	if err != nil {
		logger.Errorf("Shutting down, error while creating tracer %v", err)
		return
	}
	logger.Info("Jaeger connected")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(cfg.PrometheusConfig.Name)
	if err != nil {
		logger.Errorf("Shutting down, error while creating metrics %v", err)
		return
	}

	go func() {
		logger.Info("Metrics server running")
		if err := metrics.RunMetricServer(cfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Errorf("Shutting down, error while running metrics server %v", err)
			return
		}
	}()

	logger.Info("Database initializing")
	moviesDatabase, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}

	logger.Info("Repository initializing")
	moviesRepo := repository.NewMoviesRepository(moviesDatabase, logger.Logger)
	defer moviesRepo.Shutdown()

	genresDatabase, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}
	genresRepo := repository.NewGenresRepository(genresDatabase, logger.Logger)
	defer genresRepo.Shutdown()

	countriesDatabase, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}
	countriesRepo := repository.NewCountriesRepository(countriesDatabase, logger.Logger)
	defer countriesRepo.Shutdown()

	ageRatingsDatabase, err := repository.NewPostgreDB(cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}
	ageRatingsRepo := repository.NewAgeRatingsRepository(ageRatingsDatabase, logger.Logger)
	defer ageRatingsRepo.Shutdown()

	moviesCache, err := repository.NewMoviesCache(logger.Logger, getMoviesCacheOptions(cfg))
	if err != nil {
		logger.Errorf("Shutting down, connection to the movies cache is not established: %s", err.Error())
		return
	}
	defer moviesCache.Shutdown()

	moviesPreviewCache, err := repository.NewMoviesPreviewCache(logger.Logger, getMoviesPreviewCacheOptions(cfg))
	if err != nil {
		logger.Errorf("Shutting down, connection to the movies preview cache is not established: %s", err.Error())
		return
	}
	defer moviesPreviewCache.Shutdown()

	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{moviesDatabase, ageRatingsDatabase, genresDatabase,
			countriesDatabase, moviesCache, moviesPreviewCache}, cfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Errorf("Shutting down, error while running healthcheck endpoint %s", err.Error())
			return
		}
	}()

	imagesService := service.NewImageService(logger.Logger)

	repoManager := repository.NewMoviesRepositoryManager(moviesRepo, moviesCache, moviesRepo,
		moviesPreviewCache, ageRatingsRepo, genresRepo, countriesRepo, metric,
		repository.RepositoryManagerConfig{
			MovieTTL:        cfg.RepositoryManager.MovieTTL,
			FilteredTTL:     cfg.RepositoryManager.FilteredTTL,
			MoviePreviewTTL: cfg.RepositoryManager.MoviePreviewTTL,
		}, logger.Logger)
	logger.Info("Service initializing")
	service := service.NewMoviesService(logger.Logger, repoManager, imagesService, cfg.PicturesUrlConfig)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	s.Run(getListenServerConfig(cfg), metric, nil, nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.Shutdown()
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:        cfg.Listen.Mode,
		Host:        cfg.Listen.Host,
		Port:        cfg.Listen.Port,
		ServiceDesc: &movies_service.MoviesServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(movies_service.MoviesServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}
			return movies_service.RegisterMoviesServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
	}
}

func getMoviesCacheOptions(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Network:  cfg.MoviesCache.Network,
		Addr:     cfg.MoviesCache.Addr,
		Password: cfg.MoviesCache.Password,
		DB:       cfg.MoviesCache.DB,
	}
}
func getMoviesPreviewCacheOptions(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Network:  cfg.MoviesPreviewCache.Network,
		Addr:     cfg.MoviesPreviewCache.Addr,
		Password: cfg.MoviesPreviewCache.Password,
		DB:       cfg.MoviesPreviewCache.DB,
	}
}
