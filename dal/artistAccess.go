package dal

import (
	//mysql driver
	"log"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

//Artist data type including some basic information and location.
type Artist struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Location     `json:"location"`
	Genre        string `json:"genre"`
	SpotifyID    int    `json:"spotify_id"`
	WikipediaURL string `json:"wikipedia_url"`
}

//ArtistStore database access.
type ArtistStore struct {
	DB *sql.DB
}

//NewArtistStore returns a new connection to an Artist store
func NewArtistStore(db *sql.DB) ArtistStore {
	return ArtistStore{DB: db}
}

func (as *ArtistStore) addArtist(artist Artist) (artistID int, err error) {
	query := `
	INSERT artist
	SET Name = ?, Location = ?, Genre = ?, SpotifyID = ?, WikipediaURL = ?
	`
	res, err := as.DB.Exec(query, artist.Name, artist.Location.ID, artist.Genre, artist.SpotifyID, artist.WikipediaURL)

	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	artist.ID = int(id)

	return artist.ID, nil
}

func (as *ArtistStore) getArtistByID(artistID int) (artist Artist, err error) {
	query := `
		SELECT * FROM artist
		WHERE 
		id = ?
	`

	res := as.DB.QueryRow(query, artistID)

	if err != nil {
		log.Fatal(err)
	}

	err = res.Scan(&artist.ID, &artist.Name, &artist.Location, &artist.Genre, &artist.SpotifyID, &artist.WikipediaURL)

	return artist, err
}
