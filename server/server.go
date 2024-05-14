package server

import (
	"net/http"

	"spf-playlist/pkg/logger"
)

func Run(host, port string, srv *http.Server, log logger.Logger) {
	log.Infof("api running on: %s:%v", host, port)

	if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
		log.Errorf("api server failed: %v", err)
	}
}
