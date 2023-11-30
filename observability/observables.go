package observability

import "github.com/prometheus/client_golang/prometheus"

var FingerprintEvaluations = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "locksmith_fingerprint_evaluations",
		Help: "No of Fingerprints Evaluated",
	},
)

var FingerprintScore = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "locksmith_fingerprint_total_score",
		Help: "All scores added together",
	},
)

var FingerprintPasses = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "locksmith_fingerprint_passes",
		Help: "No of Passing Fingerprint Evaluations",
	},
)

var FingerprintRejections = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "locksmith_fingerprint_rejections",
		Help: "No of Passing Fingerprint Evaluations",
	},
)

var FingerprintHitCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "locksmith_fingerprint_hits",
		Help: "Count of fingerprint hits by category",
	},
	[]string{"category"}, // Only one label
)
