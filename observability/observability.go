package observability

import "github.com/prometheus/client_golang/prometheus"

func RegisterLocksmithObservables(registry *prometheus.Registry) {
	registry.MustRegister(FingerprintEvaluations)
	registry.MustRegister(FingerprintScore)
	registry.MustRegister(FingerprintPasses)
	registry.MustRegister(FingerprintRejections)
	registry.MustRegister(FingerprintHitCounter)
}
