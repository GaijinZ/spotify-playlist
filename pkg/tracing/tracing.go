package tracing

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Tracing struct {
	RequestID     string
	RequestMethod string
	RequestURL    *url.URL
}

func NewTrace(r *http.Request) *Tracing {
	return &Tracing{
		RequestID:     uuid.NewString(),
		RequestURL:    r.URL,
		RequestMethod: r.Method,
	}
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trace := NewTrace(r)

		logrus.WithFields(logrus.Fields{
			"request_id":     trace.RequestID,
			"request_url":    trace.RequestURL.String(),
			"request_method": trace.RequestMethod,
		}).Info("Request received")

		traceJSON, err := json.Marshal(trace)
		if err != nil {
			logrus.Errorf("Error marshaling trace object: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		traceBase64 := base64.StdEncoding.EncodeToString(traceJSON)

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:    "requestID",
			Value:   traceBase64,
			Expires: expiration,
		}

		http.SetCookie(w, &cookie)
		next.ServeHTTP(w, r)
	})
}
