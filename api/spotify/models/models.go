package models

type Token struct {
	AccessToken  string
	RefreshToken string
}

type UserProfile struct {
	ID string `json:"id"`
}

type Playlists struct {
	Items []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"items"`
}

type SearchResult struct {
	Tracks struct {
		Items []TrackRequest `json:"items"`
	} `json:"tracks"`
}

type TrackRequest struct {
	Artists []Artist `json:"artists"`
	Album   Album    `json:"album"`
	Name    string   `json:"name"`
	URI     string   `json:"uri"`
}

type Album struct {
	Name string `json:"name"`
}

type Artist struct {
	Name string `json:"name"`
}

type TrackResponse struct {
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Name   string `json:"name"`
	URI    string `json:"uri"`
}

type PayloadRequest struct {
	PlaylistName string   `json:"playlist"`
	TrackNames   []string `json:"values"`
}
