package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	labelErrandType = "errand_type"
	labelStatus     = "status"
)

var (
	errandCompletedCounter *prometheus.CounterVec
)

func init() {
	completedCounterOpts := prometheus.CounterOpts{
		Namespace: "poly",
		Subsystem: "errands",
		Name:      "completed_total",
		Help:      "A counter of all the completed errands processed by this server",
	}
	errandCompletedCounter = prometheus.NewCounterVec(completedCounterOpts, []string{labelErrandType, labelStatus})
	prometheus.MustRegister(errandCompletedCounter)
}

func ErrandCompleted(errandType string) {
	errandCompletedCounter.With(prometheus.Labels{labelErrandType: errandType, labelStatus: "completed"}).Inc()
}

func ErrandFailed(errandType string) {
	errandCompletedCounter.With(prometheus.Labels{labelErrandType: errandType, labelStatus: "failed"}).Inc()
}
