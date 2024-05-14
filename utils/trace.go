package utils

import (
	"net/http"

	"spf-playlist/pkg/logger"
)

func TrackRequestID(log logger.Logger, r *http.Request) {
	cookie, err := r.Cookie("requestID")
	if err != nil {
		log.Warningf("Trace ID cookie not found in the request: %v", err)
		return
	}

	traceID := cookie.Value

	log.Infof("Received request with trace ID: %s", traceID)
}
