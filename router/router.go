package router

import (
	"net/http"

	"spf-playlist/handler"
	"spf-playlist/pkg/tracing"
	"spf-playlist/users/handler/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Router(userAuth auth.UserAuther, spotifyHandler handler.Spotify) http.Handler {
	router := mux.NewRouter()

	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/register", userAuth.Register).Methods(http.MethodPost)
	v1.HandleFunc("/login", userAuth.Login).Methods(http.MethodPost)

	v1.Use(tracing.TraceMiddleware)

	v1.HandleFunc("/logout", userAuth.Logout).Methods(http.MethodPost)
	v1.HandleFunc("/auth", spotifyHandler.SpotifyAuth).Methods(http.MethodGet)
	v1.HandleFunc("/callback", spotifyHandler.CallbackHandler).Methods(http.MethodGet)
	v1.HandleFunc("/create-playlist", spotifyHandler.ProcessDataHandler).Methods(http.MethodPost)

	r := cors.AllowAll()
	h := r.Handler(router)

	return h
}
