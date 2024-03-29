package config

import (
	"sync"
	"time"

	"github.com/Falokut/movies_service/internal/repository"
	"github.com/Falokut/movies_service/internal/service"
	"github.com/Falokut/movies_service/pkg/jaeger"
	"github.com/Falokut/movies_service/pkg/logging"
	"github.com/Falokut/movies_service/pkg/metrics"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	HealthcheckPort string `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT"`
	Listen          struct {
		Host string `yaml:"host" env:"HOST"`
		Port string `yaml:"port" env:"PORT"`
		Mode string `yaml:"server_mode" env:"SERVER_MODE"` // support GRPC, REST, BOTH
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" ENV:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`

	PicturesUrlConfig service.PicturesUrlConfig `yaml:"pictures"`

	MoviesCache struct {
		Network  string `yaml:"network" env:"MOVIES_CACHE_NETWORK"`
		Addr     string `yaml:"addr" env:"MOVIES_CACHE_ADDR"`
		Password string `yaml:"password" env:"MOVIES_CACHE_PASSWORD"`
		DB       int    `yaml:"db" env:"MOVIES_CACHE_DB"`
	} `yaml:"movies_cache"`
	MoviesPreviewCache struct {
		Network  string `yaml:"network" env:"MOVIES_PREVIEW_CACHE_NETWORK"`
		Addr     string `yaml:"addr" env:"MOVIES_PREVIEW_CACHE_ADDR"`
		Password string `yaml:"password" env:"MOVIES_PREVIEW_CACHE_PASSWORD"`
		DB       int    `yaml:"db" env:"MOVIES_PREVIEW_CACHE_DB"`
	} `yaml:"movies_preview_cache"`

	RepositoryManager struct {
		MovieTTL        time.Duration `yaml:"movie_ttl"`
		FilteredTTL     time.Duration `yaml:"filtered_ttl"`
		MoviePreviewTTL time.Duration `yaml:"movie_preview_ttl"`
	} `yaml:"repository_manager"`

	DBConfig     repository.DBConfig `yaml:"db_config"`
	JaegerConfig jaeger.Config       `yaml:"jaeger"`
}

var instance *Config
var once sync.Once

const configsPath = "configs/"

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})

	return instance
}
