package dal

import (
	//mysql driver
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"log"
	"time"
)

//AlbumStore database access.
type AlbumStore struct {
	DB *sql.DB
}

//Album of an artist.
type Album struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ArtistID    int       `json:"artist_id"`
	ReleaseDate time.Time `json:"date"`
}

//NewAlbumStore returns a new connection to an album store
func NewAlbumStore(db *sql.DB) AlbumStore {
	return AlbumStore{DB: db}
}

func (as *AlbumStore) addAlbums(albums []Album) (newAlbums []Album) {
	for _, a := range albums {
		as.addAlbum(a)
	}
	return albums
}

func (as *AlbumStore) addAlbum(album Album) (albumID int, err error) {
	query := `
	INSERT album
	SET title = ,? artist_id = ?, release_date = ?
	`
	res, err := as.DB.Exec(query, album)

	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	album.ID = int(id)

	return album.ID, nil
}

// GetAlbumsByArtistID returns all the albums that are linked to an artist.
func (as *AlbumStore) GetAlbumsByArtistID(artistID int) (a []Album, err error) {
	query := `
		SELECT *
		FROM 
			album
		WHERE
			artist_id  = ?
	`
	rows, err := as.DB.Query(query, artistID)

	if err != nil {
		log.Fatal(err)
	}
	var albums []Album

	for rows.Next() {
		var artistAlbum Album
		err := rows.Scan(&artistAlbum.ID, &artistAlbum.Title, &artistAlbum.ArtistID, &artistAlbum.ReleaseDate)
		if err != nil {
			log.Fatal(err)
		}

		albums = append(albums, artistAlbum)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return albums, nil
}
