package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"spf-playlist/api/spotify/auth"
	"spf-playlist/api/spotify/handler"
	"spf-playlist/api/spotify/models"
	"spf-playlist/pkg/config"
	"spf-playlist/pkg/logger"
	"spf-playlist/utils"
)

type Spotify struct {
	token       models.Token
	ctx         context.Context
	spotifyAuth auth.SpotifyAuth
	cfg         config.GlobalEnv
}

func NewSpotifyHandler(
	token models.Token,
	ctx context.Context,
	spotifyAuth auth.SpotifyAuth,
	cfg config.GlobalEnv,
) *Spotify {
	return &Spotify{
		token:       token,
		ctx:         ctx,
		spotifyAuth: spotifyAuth,
		cfg:         cfg,
	}
}

func (s *Spotify) SpotifyAuth(w http.ResponseWriter, r *http.Request) {
	log := utils.GetLogger(s.ctx)

	utils.TrackRequestID(log, r)

	s.spotifyAuth.SpotifyAuth(w, r)
}

func (s *Spotify) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	log := utils.GetLogger(s.ctx)

	utils.TrackRequestID(log, r)

	if r.URL.Query().Get("state") != "state" {
		log.Errorf("Invalid state parameter")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Errorf("No code parameter")
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	token, err := s.spotifyAuth.ExchangeToken(code)
	if err != nil {
		log.Errorf("Error exchanging token: %v", err)
		http.Error(w, fmt.Sprintf("Error exchanging code for token: %v", err), http.StatusInternalServerError)
		return
	}

	s.token.AccessToken = token.AccessToken
	s.token.RefreshToken = token.RefreshToken

	log.Infof("Access token: %v\n", s.token.AccessToken)
	log.Infof("Refresh token: %v\n", s.token.RefreshToken)
}

func (s *Spotify) ProcessDataHandler(w http.ResponseWriter, r *http.Request) {
	var playlistName string

	log := utils.GetLogger(s.ctx)

	utils.TrackRequestID(log, r)

	if r.Method != http.MethodPost {
		log.Errorf("Invalid method: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload := &models.PayloadRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Errorf("Error decoding payload: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	userID, err := getUserProfile(s.token.AccessToken, s.cfg, log)
	if err != nil {
		log.Errorf("Error getting user profile: %v", err)
		http.Error(w, fmt.Sprintf("Error getting user profile: %v", err), http.StatusInternalServerError)
		return
	}

	s.ctx = context.WithValue(s.ctx, "userID", userID)

	playlistName, hasPlaylist, err := handler.HasPlaylist(payload.PlaylistName, s.token.AccessToken, s.cfg, log)
	if err != nil {
		log.Errorf("Error checking playlist: %v", err)
		http.Error(w, fmt.Sprintf("Error checking playlist: %v", err), http.StatusInternalServerError)
		return
	}

	if !hasPlaylist {
		playlistName, err = handler.CreatePlaylist(payload.PlaylistName, s.token.AccessToken, s.cfg, s.ctx, log)
		if err != nil {
			log.Errorf("Error creating playlist: %v", err)
			http.Error(w, fmt.Sprintf("Error creating playlist: %v", err), http.StatusInternalServerError)
			return
		}
	}

	tracksURI, err := handler.GetTrackURI(payload.TrackNames, s.token.AccessToken, s.cfg, log)
	if err != nil {
		log.Errorf("Error getting track URI: %v", err)
		http.Error(w, fmt.Sprintf("Error getting track URI: %v", err), http.StatusInternalServerError)
		return
	}

	err = handler.AddToPlaylist(playlistName, s.token.AccessToken, tracksURI, s.cfg, log)
	if err != nil {
		log.Errorf("Error adding playlist: %v", err)
		http.Error(w, fmt.Sprintf("Error adding playlist: %v", err), http.StatusInternalServerError)
		return
	}
}

func getUserProfile(accessToken string, cfg config.GlobalEnv, log logger.Logger) (string, error) {
	req, err := http.NewRequest("GET", cfg.BaseHost+"/me", nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	userProfile := &models.UserProfile{}

	if err = json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		log.Errorf("Error decoding user profile: %v", err)
		return "", err
	}

	return userProfile.ID, nil
}
