package dal

import (

	//mysql driver
	"context"
	"log"

	"github.com/jackc/pgx/v4"
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
	DB *pgx.Conn
}

//NewLocationStore returns a new connection to an Location store
func NewLocationStore(db *pgx.Conn) LocationStore {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return LocationStore{DB: db}
}

func (ls *LocationStore) AddLocation(location Location) (locationID int, err error) {

	var locationId int

	// query := `
	// INSERT location
	// SET City = ?, State = ?, Country = ?, full_location = ?, google_place_id = ?, latitude = ?, longitude = ?
	// `
	query := `
	INSERT into bands_from_town.location ( City, State, Country, full_location, google_place_id, latitude, longitude)
	values ($1, $2, $3, $4, $5, $6, $7) returning id;`

	// res, err := ls.DB.Exec(query, location.City, location.State, location.Country, location.FullLocation, location.GooglePlaceID, location.Latitude, location.Longitude)

	err = ls.DB.QueryRow(context.Background(), query, location.City, location.State, location.Country, location.FullLocation, location.GooglePlaceID, location.Latitude, location.Longitude).Scan(&locationId)

	if err != nil {
		log.Print(err)
	}

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return locationId, err
}

func (ls *LocationStore) UpdateLocation(location Location) (locationID int, err error) {
	query := `
	UPDATE bands_from_town.location
	SET City = $1, State = $2, Country = $3, full_location = $4, google_place_id = $5, latitude = $6, longitude = $7
	WHERE id = $8`

	_, err = ls.DB.Exec(context.Background(), query, location.City, location.State, location.Country, location.FullLocation, location.GooglePlaceID, location.Latitude, location.Longitude, location.ID)

	if err != nil {
		log.Print(err)
		return location.ID, err
	}

	return location.ID, nil
}

func (ls *LocationStore) GetLocationByID(locationID int) (location Location, err error) {
	query := `
	SELECT id, city, state, country, full_location, google_place_id, coalesce(latitude, 0), coalesce(longitude, 0) 
	FROM bands_from_town.location
	WHERE 
	id = $1`

	res := ls.DB.QueryRow(context.Background(), query, locationID)

	if err != nil {
		log.Fatal(err)
	}

	err = res.Scan(&location.ID, &location.City, &location.State, &location.Country, &location.FullLocation, &location.GooglePlaceID, &location.Latitude, &location.Longitude)

	return location, err
}

func (ls *LocationStore) GetLocationByGoogleID(locationID string) (location Location, err error) {
	query := `
	SELECT id, city, state, country, full_location, google_place_id, coalesce(latitude, 0), coalesce(longitude, 0) 
	FROM bands_from_town.location  		
	WHERE 
	google_place_id = $1
	`
	res := ls.DB.QueryRow(context.Background(), query, locationID)

	err = res.Scan(&location.ID, &location.City, &location.State, &location.Country, &location.FullLocation, &location.GooglePlaceID, &location.Latitude, &location.Longitude)

	return location, err
}

func (ls *LocationStore) GetArtistsByLocationID(locationID int) (artists []Artist, err error) {
	query := `
		SELECT * FROM bands_from_town.artist
		WHERE 
		hometown = $1
	`
	rows, err := ls.DB.Query(context.Background(), query, locationID)

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

func (ls *LocationStore) CheckForExistingLocation(locationToCheck Location) (bool, Location) {
	existingLocation, err := ls.GetLocationByGoogleID(locationToCheck.GooglePlaceID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, locationToCheck
		}

		log.Fatal(err)
		return false, locationToCheck

	}
	return true, existingLocation
}
