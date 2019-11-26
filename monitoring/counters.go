package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

/*type Counters struct {
	FooCount	prometheus.Counter
	Hits		*prometheus.CounterVec
}

func NewCounters() *Counters {
	return &Counters{
		FooCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "foo_total",
			Help: "Number of foo successfully processed.",
		}),
		Hits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hits",
		}, []string{"status", "path"}),
	}
}

func (c * Counters) Register() {
	prometheus.MustRegister(c.FooCount, c.Hits)
}*/

var (
	RequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
	})

	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status" , "path"})
)
