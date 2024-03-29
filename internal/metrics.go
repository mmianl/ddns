package internal

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	VersionGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ddns_build_info",
			Help: "Metric with a constant '1' value labeled by version and goversion from which ddns was built.",
		},
		[]string{"version", "goversion"},
	)
	StartTimeGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ddns_start_time_seconds",
			Help: "Start time of the process since unix epoch in seconds.",
		},
		[]string{},
	)
	DNSARecordInfoGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ddns_dns_a_record_info",
			Help: "Metric with a constant '1' value showing the current a records and their ip addresses.",
		},
		[]string{"ip_address", "a_record"},
	)
)

// Metrics Return a httprouter.Handle function that handles metrics requests
func Metrics() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}
