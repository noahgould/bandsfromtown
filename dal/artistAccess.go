package dal

import (
	//mysql driver
	"log"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

//Artist data type including some basic information and location.
type Artist struct {
	ID           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Location     ` db:"hometown" json:"location"`
	Genre        string `db:"genre" json:"genre"`
	SpotifyID    string `db:"spotify_id" json:"spotify_id"`
	WikipediaURL string `db:"wikipedia_url" json:"wikipedia_url"`
}

//ArtistStore database access.
type ArtistStore struct {
	DB *sql.DB
}

//NewArtistStore returns a new connection to an Artist store
func NewArtistStore(db *sql.DB) ArtistStore {
	return ArtistStore{DB: db}
}

func (as *ArtistStore) AddArtist(artist Artist) (artistID int, err error) {
	query := `
	INSERT artist
	SET name = ?, hometown = ?, genre = ?, spotify_id = ?, wikipedia_url = ?
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

//GetArtistByID returns the artist with a matching id.
func (as *ArtistStore) GetArtistByID(artistID int) (artist Artist, err error) {
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

//GetArtistBySpotifyID returns the artist with a matching spotify id.
func (as *ArtistStore) GetArtistBySpotifyID(spotifyID string) (artist Artist, err error) {
	query := `
		SELECT * FROM artist
		WHERE 
		spotify_id = ?`

	res := as.DB.QueryRow(query, spotifyID)

	if err != nil {
		log.Fatal(err)
	}

	err = res.Scan(&artist.ID, &artist.Name, &artist.Location, &artist.Genre, &artist.SpotifyID, &artist.WikipediaURL)

	return artist, err
}

//GetArtistsByName returns the artist with a matching name.
func (as *ArtistStore) GetArtistsByName(artistName string) (artists []Artist, err error) {
	query := `
		SELECT * FROM artist
		WHERE 
		name = ?`

	rows, err := as.DB.Query(query, artistName)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var artist Artist
		err = rows.Scan(&artist.ID, &artist.Name, &artist.Location.ID, &artist.Genre, &artist.SpotifyID, &artist.WikipediaURL)
		if err != nil {
			log.Println(err)
		}

		artists = append(artists, artist)
	}

	return artists, err
}

func (as *ArtistStore) UpdateArtist(artist Artist) (artistId int, err error) {
	query := `UPDATE artist
		SET name = ?, hometown = ?, genre = ?, spotify_id = ?, wikipedia_url = ?
		WHERE id = ?`

	_, err = as.DB.Exec(query, artist.Name, artist.Location.ID, artist.Genre, artist.SpotifyID, artist.WikipediaURL, artist.ID)

	if err != nil {
		log.Print(err)
		return artist.ID, err
	}

	return artist.ID, nil
}
