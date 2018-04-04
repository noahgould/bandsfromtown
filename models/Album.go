package models

import "time"

//Album of an artist.
type Album struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ArtistID    int       `json:"artist_id"`
	ReleaseDate time.Time `json:"date"`
}
