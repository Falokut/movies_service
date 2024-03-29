package metrics

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncCacheHits(method string, times int32)
	IncCacheMiss(method string, times int32)
	IncHits(status int, method, path string)
	IncRestPanicsTotal()
	IncGrpcPanicsTotal()
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type PrometheusMetrics struct {
	HitsTotal             prometheus.Counter
	Hits                  *prometheus.CounterVec
	CacheHits             *prometheus.CounterVec
	CacheMiss             *prometheus.CounterVec
	Times                 *prometheus.HistogramVec
	RestPanicRecoverTotal prometheus.Counter
	GrpcPanicRecoverTotal prometheus.Counter
}

func CreateMetrics(name string) (*PrometheusMetrics, error) {
	var metr PrometheusMetrics
	metr.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})
	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, err
	}

	metr.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.Hits); err != nil {
		return nil, err
	}

	metr.CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_cache_hits",
		},
		[]string{"method"},
	)

	if err := prometheus.Register(metr.CacheHits); err != nil {
		return nil, err
	}

	metr.CacheMiss = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_cache_miss",
		},
		[]string{"method"},
	)
	if err := prometheus.Register(metr.CacheMiss); err != nil {
		return nil, err
	}

	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"status", "method", "path"},
	)
	if err := prometheus.Register(metr.Times); err != nil {
		return nil, err
	}

	metr.RestPanicRecoverTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_rest_panic_recover_total",
	})
	if err := prometheus.Register(metr.RestPanicRecoverTotal); err != nil {
		return nil, err
	}

	metr.GrpcPanicRecoverTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_grpc_panic_recover_total",
	})
	if err := prometheus.Register(metr.GrpcPanicRecoverTotal); err != nil {
		return nil, err
	}
	if err := prometheus.Register(collectors.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	return &metr, nil
}

type MetricsServerConfig struct {
	Host string `yaml:"host" env:"METRIC_HOST"`
	Port string `yaml:"port" env:"METRIC_PORT"`
}

func RunMetricServer(cfg MetricsServerConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	return http.Serve(lis, mux)
}
func (metr *PrometheusMetrics) IncHits(status int, method, path string) {
	metr.HitsTotal.Inc()
	metr.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (metr *PrometheusMetrics) IncCacheHits(method string, times int32) {
	metr.CacheHits.WithLabelValues(method).Add(float64(times))
}

func (metr *PrometheusMetrics) IncCacheMiss(method string, times int32) {
	metr.CacheMiss.WithLabelValues(method).Add(float64(times))
}

func (metr *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	metr.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

func (metr *PrometheusMetrics) IncRestPanicsTotal() {
	metr.RestPanicRecoverTotal.Inc()
}

func (metr *PrometheusMetrics) IncGrpcPanicsTotal() {
	metr.GrpcPanicRecoverTotal.Inc()
}
