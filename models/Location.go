package models

//Location of an artist or band.
type Location struct {
	ID      int    `json:"id"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}
