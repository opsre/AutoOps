package agent

import (
	"log"
	"net/http"
	"time"
)

func monitorEndpoints(config *Config) {
	for _, endpoint := range config.Endpoints {
		start := time.Now()
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(endpoint.URL)
		duration := time.Since(start).Milliseconds()

		status := 0
		if err == nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				status = 1
			}
			resp.Body.Close()
		}

		endpointResponseTime.WithLabelValues(hostname, endpoint.URL, endpoint.Method).Set(float64(duration))
		endpointStatus.WithLabelValues(hostname, endpoint.URL, endpoint.Method).Set(float64(status))
		log.Printf("Endpoint %s %s - Status: %d, Response: %dms",
			endpoint.Method, endpoint.URL, status, duration)
	}
}
