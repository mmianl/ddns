package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestMetrics tests that the metrics endpoint returns the expected metrics
func TestMetrics(t *testing.T) {
	// Start the server
	ddnsVersionGauge.WithLabelValues(Version, runtime.Version()).Set(1)
	now := time.Now()
	ddnsStartTimeGauge.WithLabelValues().Set(float64(now.Unix()))
	ddnsDNSARecordInfoGauge.WithLabelValues("127.0.0.1", "test.example.com").SetToCurrentTime()

	router := httprouter.New()
	router.GET("/metrics", Metrics())

	// Make a request to the metrics endpoint
	req, _ := http.NewRequest("GET", "/metrics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Verify that the response status code is 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Verify that the response body contains the expected metrics
	if !strings.Contains(rr.Body.String(), "ddns_build_info") {
		t.Errorf("Expected response body to contain ddns_build_info, got %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ddns_start_time_seconds") {
		t.Errorf("Expected response body to contain ddns_start_time_seconds, got %s", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "ddns_dns_a_record_info") {
		t.Errorf("Expected response body to contain ddns_dns_a_record_info, got %s", rr.Body.String())
	}
}
