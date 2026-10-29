package metrics

import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"

	"server/internal/utils/errors"
)

const (
	isSuccessLabel = "is_success"
	methodLabel    = "method"
)

type metrics struct {
	ResponseTimeMetric *prometheus.HistogramVec
}

func (m *metrics) register() error {

	if err := prometheus.Register(m.ResponseTimeMetric); err != nil {
		return errors.InternalServer.Wrap(err)
	}

	return nil
}

var metricsInitialized atomic.Bool
var GlobalMetrics *metrics

func Init(namespace string) error {

	if metricsInitialized.Load() {
		return nil
	}
	metricsInitialized.Store(true)

	GlobalMetrics = &metrics{

		// Метрика для измерения времени ответа
		ResponseTimeMetric: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace:                       namespace,
			Subsystem:                       "",
			Name:                            "http_response_time_seconds",
			Help:                            "Histogram of response time in seconds.",
			ConstLabels:                     nil,
			Buckets:                         nil,
			NativeHistogramBucketFactor:     0,
			NativeHistogramZeroThreshold:    0,
			NativeHistogramMaxBucketNumber:  0,
			NativeHistogramMinResetDuration: 0,
			NativeHistogramMaxZeroThreshold: 0,
		}, []string{methodLabel, isSuccessLabel}),
	}

	// Регистрируем метрики
	return GlobalMetrics.register()
}
