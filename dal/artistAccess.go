package dal

import (
	//mysql driver
	"context"
	"log"

	"github.com/jackc/pgx/v4"
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
	DB *pgx.Conn
}

//NewArtistStore returns a new connection to an Artist store
func NewArtistStore(db *pgx.Conn) ArtistStore {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return ArtistStore{DB: db}
}

//AddArtist saves an artist to the database.
func (as *ArtistStore) AddArtist(artist Artist) (artistID int, err error) {
	if artist.SpotifyID == "" {
		artist.SpotifyID = "-1"
	}

	query := `INSERT INTO bands_from_town.artist (name, hometown, genre, spotify_id, wikipedia_url)
		values ($1, $2, $3, $4, $5) returning id;`

	err = as.DB.QueryRow(context.Background(), query, artist.Name, artist.Location.ID, artist.Genre, artist.SpotifyID, artist.WikipediaURL).Scan(&artistID)

	if err != nil {
		log.Fatal(err)
	}

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// artist.ID = int(id)

	return artist.ID, err
}

//GetArtistByID returns the artist with a matching id.
func (as *ArtistStore) GetArtistByID(artistID int) (artist Artist, err error) {
	query := `
		SELECT * FROM artist
		WHERE 
		id = $1
	`

	res := as.DB.QueryRow(context.Background(), query, artistID)

	if err != nil {
		log.Fatal(err)
	}

	err = res.Scan(&artist.ID, &artist.Name, &artist.Location.ID, &artist.Genre, &artist.SpotifyID, &artist.WikipediaURL)

	return artist, err
}

//GetArtistBySpotifyID returns the artist with a matching spotify id.
func (as *ArtistStore) GetArtistBySpotifyID(spotifyID string) (artist Artist, err error) {
	if spotifyID == "-1" {
		return artist, pgx.ErrNoRows
	}

	query := `
		SELECT * FROM bands_from_town.artist
		WHERE 
		spotify_id = $1`

	res := as.DB.QueryRow(context.Background(), query, spotifyID)

	err = res.Scan(&artist.ID, &artist.Name, &artist.Location.ID, &artist.Genre, &artist.SpotifyID, &artist.WikipediaURL)

	return artist, err
}

//GetArtistsByName returns the artist with a matching name.
func (as *ArtistStore) GetArtistsByName(artistName string) (artists []Artist, err error) {
	query := `
		SELECT * FROM bands_from_town.artist
		WHERE name = $1;`

	rows, err := as.DB.Query(context.Background(), query, artistName)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var artist Artist
		err = rows.Scan(&artist.ID, &artist.Name, &artist.Location.ID, &artist.Genre, &artist.SpotifyID, &artist.WikipediaURL)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		artists = append(artists, artist)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return artists, err
}

//UpdateArtist updates an existing artist.
func (as *ArtistStore) UpdateArtist(artist Artist) (artistId int, err error) {
	query := `UPDATE bands_from_town.artist
		SET name = $1, hometown = $2, genre = $3, spotify_id = $4, wikipedia_url = $5
		WHERE id = $6`

	_, err = as.DB.Exec(context.Background(), query, artist.Name, artist.Location.ID, artist.Genre, artist.SpotifyID, artist.WikipediaURL, artist.ID)

	if err != nil {
		log.Print(err)
		return artist.ID, err
	}

	return artist.ID, nil
}
