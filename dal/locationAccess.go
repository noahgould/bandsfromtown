package dal

import (

	//mysql driver
	"log"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

//Location of an artist or band.
type Location struct {
	ID            int     `json:"id"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Country       string  `json:"country"`
	FullLocation  string  `json:"location_string"`
	GooglePlaceID string  `json:"google_place_id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

//LocationStore database access.
type LocationStore struct {
	DB *sql.DB
}

//NewLocationStore returns a new connection to an Location store
func NewLocationStore(db *sql.DB) LocationStore {
	return LocationStore{DB: db}
}

func (ls *LocationStore) AddLocation(location Location) (locationID int, err error) {
	query := `
	INSERT location
	SET City = ?, State = ?, Country = ?, full_location = ?, google_place_id = ?, latitude = ?, longitude = ?
	`
	res, err := ls.DB.Exec(query, location.City, location.State, location.Country, location.FullLocation, location.GooglePlaceID, location.Latitude, location.Longitude)

	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	location.ID = int(id)

	return location.ID, nil
}

func (ls *LocationStore) UpdateLocation(location Location) (locationID int, err error) {
	query := `
	UPDATE location
	SET City = ?, State = ?, Country = ?, full_location = ?, google_place_id = ?, latitude = ?, longitude = ?
	WHERE id = ?
	`

	_, err = ls.DB.Exec(query, location.City, location.State, location.Country, location.FullLocation, location.GooglePlaceID, location.Latitude, location.Longitude, location.ID)

	if err != nil {
		log.Print(err)
		return location.ID, err
	}

	return location.ID, nil
}

func (ls *LocationStore) GetLocationByID(locationID int) (location Location, err error) {
	query := `
	SELECT id, city, state, country, full_location, google_place_id, coalesce(latitude, 0), coalesce(longitude, 0) 
	FROM location
	WHERE 
	id = ?
	`
	res := ls.DB.QueryRow(query, locationID)

	if err != nil {
		log.Fatal(err)
	}

	err = res.Scan(&location.ID, &location.City, &location.State, &location.Country, &location.FullLocation, &location.GooglePlaceID, &location.Latitude, &location.Longitude)

	return location, err
}

func (ls *LocationStore) GetLocationByGoogleID(locationID string) (location Location, err error) {
	query := `
	SELECT id, city, state, country, full_location, google_place_id, coalesce(latitude, 0), coalesce(longitude, 0) 
	FROM location  		
	WHERE 
	google_place_id = ?
	`
	res := ls.DB.QueryRow(query, locationID)

	err = res.Scan(&location.ID, &location.City, &location.State, &location.Country, &location.FullLocation, &location.GooglePlaceID, &location.Latitude, &location.Longitude)

	return location, err
}

func (ls *LocationStore) GetArtistsByLocationID(locationID int) (artists []Artist, err error) {
	query := `
		SELECT * FROM artist
		WHERE 
		hometown = ?
	`
	rows, err := ls.DB.Query(query, locationID)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var artist Artist
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Location, &artist.Genre, &artist.SpotifyID, &artist.WikipediaURL)
		if err != nil {
			log.Fatal(err)
		}

		artists = append(artists, artist)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return artists, nil
}
