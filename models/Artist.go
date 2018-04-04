package models

//Artist or band saved in database .
type Artist struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Location     `json:"location"`
	Genre        string `json:"genre"`
	SpotifyID    int    `json:"spotify_id"`
	WikipediaURL string `json:"wikipedia_url"`
}

// func newArtist(id, name, origin, genre, spotifyId, WikipediaURL) *artist {
// }
