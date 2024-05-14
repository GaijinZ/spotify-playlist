package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"spf-playlist/api/spotify/models"
	"spf-playlist/pkg/config"
	"spf-playlist/pkg/logger"
)

func HasPlaylist(playlistName, accessToken string, cfg config.GlobalEnv, log logger.Logger) (string, bool, error) {
	req, err := http.NewRequest("GET", cfg.BaseHost+"/me/playlists", nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return "", false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending request: %v", err)
		return "", false, err
	}
	defer resp.Body.Close()

	playlists := &models.Playlists{}

	if err = json.NewDecoder(resp.Body).Decode(&playlists); err != nil {
		log.Errorf("Error decoding response: %v", err)
		return "", false, err
	}

	for _, item := range playlists.Items {
		if item.Name == playlistName {
			log.Infof("Found playlist: %s", item.Name)
			return item.ID, true, nil
		}
	}

	log.Infof("No playlist: %s", playlistName)
	return "", false, err
}

func CreatePlaylist(name, accessToken string, cfg config.GlobalEnv, ctx context.Context, log logger.Logger) (string, error) {
	userID := ctx.Value("userID").(string)

	url := fmt.Sprintf(cfg.BaseHost+"/users/%s/playlists", userID)

	playlistData := map[string]interface{}{
		"name":        name,
		"description": "Created using the Spotify API and Go",
		"public":      true,
	}

	playlistJSON, err := json.Marshal(playlistData)
	if err != nil {
		log.Errorf("Error marshalling playlist data: %s", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(playlistJSON))
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error making request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Errorf("Error creating playlist: %v", resp.Status)
		return "", fmt.Errorf("failed to create playlist (status code: %d)", resp.StatusCode)
	}

	var response map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	playlistID, ok := response["id"].(string)
	if !ok {
		log.Errorf("Error getting playlist ID: %v", response["id"])
		return "", fmt.Errorf("unable to extract playlist ID from response")
	}

	log.Infof("Created playlist: %s", name)
	return playlistID, nil
}

func SearchTrack(trackName, accessToken string, cfg config.GlobalEnv, log logger.Logger) (*models.TrackResponse, error) {
	searchResult := &models.SearchResult{}
	trackResponse := &models.TrackResponse{}

	url := fmt.Sprintf(cfg.BaseHost+"/search?q=track:%s&type=track&limit=50", trackName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
		log.Errorf(errMsg)
		return trackResponse, fmt.Errorf(errMsg)
	}

	if err = json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		log.Errorf("Error decoding response: %v", err)
		return trackResponse, err
	}

	for _, track := range searchResult.Tracks.Items {
		if strings.ToLower(track.Name) == strings.ToLower(trackName) {
			for _, artist := range track.Artists {
				trackResponse.Artist += artist.Name + ", "
			}
			if len(trackResponse.Artist) > 0 {
				trackResponse.Artist = trackResponse.Artist[:len(trackResponse.Artist)-2]
			}

			trackResponse.Album = track.Album.Name
			trackResponse.Name = track.Name
			trackResponse.URI = track.URI
			return trackResponse, nil
		}
	}

	return trackResponse, nil
}

func AddToPlaylist(playlist, accessToken string, trackURI []string, cfg config.GlobalEnv, log logger.Logger) error {
	url := fmt.Sprintf(cfg.BaseHost+"/playlists/%s/tracks", playlist)

	requestBody := map[string]interface{}{
		"uris": trackURI,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		log.Errorf("Error marshaling request body: %s", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		log.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error making request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("failed to add track to playlist (status code: %d)", resp.StatusCode)
		return fmt.Errorf("failed to add track to playlist (status code: %d)", resp.StatusCode)
	}

	return nil
}

func GetTrackURI(trackNames []string, accessToken string, cfg config.GlobalEnv, log logger.Logger) ([]string, error) {
	tracksURI := make([]string, 0, len(trackNames))
	copy(tracksURI, trackNames)

	for _, trackName := range trackNames {
		track := strings.ReplaceAll(trackName, " ", "+")
		tracks, err := SearchTrack(track, accessToken, cfg, log)
		if err != nil {
			log.Errorf("Error searching tracks: %s", err)
			return tracksURI, err
		}

		log.Infof("Search result for track '%s':\n", trackName)
		if tracks.URI == "" {
			log.Warningf("No URI found for track '%s'\n", trackName)
		} else {
			log.Infof("Track found:\n Artist: %s\n Album: %s\n Name: %s\n", tracks.Artist, tracks.Album, tracks.Name)
			tracksURI = append(tracksURI, tracks.URI)
		}
	}

	return tracksURI, nil
}
