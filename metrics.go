package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ddnsVersionGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ddns_build_info",
			Help: "Metric with a constant '1' value labeled by version and goversion from which ddns was built.",
		},
		[]string{"version", "goversion"},
	)
	ddnsStartTimeGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ddns_start_time_seconds",
			Help: "Start time of the process since unix epoch in seconds.",
		},
		[]string{},
	)
	ddnsDNSARecordUpdateTimeGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ddns_dns_a_record_update_time_seconds",
			Help: "Time of last DNS A record update since unix epoch in seconds labeled by IP address and A Record.",
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
