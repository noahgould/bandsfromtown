package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/noahgould/bandsfromtown/dal"
)

type SpotifyController struct {
	clientID      string
	clientSecret  string
	redirectURI   string
	artistStore   dal.ArtistStore
	locationStore dal.LocationStore
}

type spotifyTokenResponse struct {
	accessToken    string `json:"access_token"`
	tokenType      string `json:"token_type"`
	scope          string `json:"scope"`
	expirationTime int    `json:"expires_in"`
	refreshToken   string `json:"refresh_token"`
}

type spotifyAlbum struct {
	AlbumType            string                `json:"album_type"`
	Artists              []spotifySimpleArtist `json:"artists"`
	AvailableMarkets     []string              `json:"available_markets"`
	Copyrights           []string              `json:"copyrights"`
	ExternalIds          []string              `json:"external_ids"`
	ExternalUrls         []string              `json:"external_urls"`
	Genres               []string              `json:"genres"`
	Href                 string                `json:"href"`
	ID                   string                `json:"id"`
	Images               []string              `json:"images"`
	Label                string                `json:"label"`
	Name                 string                `json:"name"`
	Popularity           int                   `json:"popularity"`
	ReleaseDate          string                `json:"release_date"`
	ReleaseDatePrecision string                `json:"release_date_precision"`
	Restrictions         []string              `json:"restrictions"`
	Tracks               []string              `json:"tracks"`
	ObjectType           string                `json:"type"`
	URI                  string                `json:"uri"`
}

type spotifySimpleArtist struct {
	ExternalUrls []string `json:"external_urls"`
	Href         string   `json:"href"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ObjectType   string   `json:"type"`
	URI          string   `json:"uri"`
}

type spotifyPage struct {
	Href     string       `json:"href"`
	Items    []savedAlbum `json:"items"`
	Limit    int          `json:"limit"`
	Next     string       `json:"next"`
	Offset   int          `json:"offset"`
	Previous string       `json:"previous"`
	Total    int          `json:"total"`
}

type savedAlbum struct {
	AddedAt string       `json:"added_at"`
	Album   spotifyAlbum `json:"album"`
}

const redirectURI string = "http://localhost:8080/spotify/login/"

func NewSpotifyController(newArtistStore dal.ArtistStore, newLocationStore dal.LocationStore) *SpotifyController {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	clientID, clientEnvExist := os.LookupEnv("SPOTIFY_ID")
	clientSecret, clientSecretEnvExist := os.LookupEnv("SPOTIFY_SECRET")
	if !clientEnvExist || !clientSecretEnvExist {
		log.Fatal("spotify client id or secret not stored in environment variables.")
	}

	return &SpotifyController{
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURI:   redirectURI,
		artistStore:   newArtistStore,
		locationStore: newLocationStore,
	}
}

func (sc *SpotifyController) AuthorizationRequest(w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Add("client_id", sc.clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", sc.redirectURI)
	q.Add("scope", "user-library-read playlist-read-collaborative playlist-read-private")

	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)

}

func (sc *SpotifyController) AuthorizationCallback(w http.ResponseWriter, r *http.Request) {

	authCode := r.URL.Query()["code"]
	errorCode := r.URL.Query()["error"]

	if len(errorCode) > 0 {
		log.Println(errorCode)
		w.Write([]byte("Error authenticating."))
	}

	if len(authCode) == 0 {
		log.Println("More than 1 authcode.")
	}

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode[0])
	form.Add("redirect_uri", sc.redirectURI)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))

	if err != nil {
		log.Println(err)
	}

	headerString := base64.StdEncoding.EncodeToString([]byte(sc.clientID + ":" + sc.clientSecret))
	headerString = strings.Join([]string{"Basic", headerString}, " ")
	req.Header.Add("Authorization", headerString)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	spotifyClient := &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := spotifyClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Println(readErr)
	}

	tokenResult := spotifyTokenResponse{}

	err = json.Unmarshal(body, &tokenResult)

	if err != nil {
		log.Println(err)
	}

	log.Println(tokenResult.accessToken)
	usersArtists := sc.getAllUserArtists(tokenResult.accessToken)

	if err := json.NewEncoder(w).Encode(usersArtists); err != nil {
		log.Println(err)
	}

	w.Write(body)

}

func (sc *SpotifyController) MapUserArtists(w http.ResponseWriter, r *http.Request) {

	usersArtists := sc.getAllUserArtists(mux.Vars(r)["spotifyID"])

	if err := json.NewEncoder(w).Encode(usersArtists); err != nil {
		log.Println(err)
	}
}

func (sc *SpotifyController) getAllUserArtists(userToken string) []dal.Artist {

	resultPage := makeArtistRequest(userToken, 0)
	artistList := []dal.Artist{}

	if resultPage.Total > 50 {
		for numAlbums := 0; numAlbums <= resultPage.Total; numAlbums += 50 {
			artistList = append(artistList, processArtists(resultPage)...)
			resultPage = makeArtistRequest(userToken, numAlbums)
		}
	} else {
		artistList = processArtists(resultPage)
	}

	artistList = sc.getArtistLocations(artistList)

	return artistList

}

func (sc *SpotifyController) getArtistLocations(artists []dal.Artist) []dal.Artist {

	var err error
	for i, artist := range artists {
		artists[i], err = sc.artistStore.GetArtistBySpotifyID(artist.SpotifyID)
		if err != nil {
			if err == sql.ErrNoRows {
				possibleArtists, err := sc.artistStore.GetArtistsByName(artist.Name)
				if err != nil {
					if err == sql.ErrNoRows {
						artist.Location = LookupArtistLocation(artist.Name)
						gmc := NewGoogleMapsController()
						artist.Location = *gmc.NormalizeLocation(artist.Location)
						var exists bool
						exists, artist.Location = sc.locationStore.CheckForExistingLocation(artist.Location)
						if !exists {
							locationPointer, err := gmc.GetCoordinates(artist.Location)
							if err != nil {
								log.Println(err)
							}
							artist.Location = *locationPointer
							sc.locationStore.AddLocation(artist.Location)
							sc.artistStore.AddArtist(artist)
						} else {
							sc.artistStore.AddArtist(artist)
						}
						artists[i] = artist
					}
				}
				artists[i] = possibleArtists[0]
			} else {
				log.Println(err)
			}
		}
	}

	return artists
}

func makeArtistRequest(userToken string, offset int) spotifyPage {

	spotifyClient := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/albums", nil)

	req.Header.Add("Authorization", userToken)

	q := req.URL.Query()
	q.Add("limit", "50")
	q.Add("response_type", "code")
	q.Add("offset", strconv.Itoa(offset))

	req.URL.RawQuery = q.Encode()

	response, err := spotifyClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	firstPage := spotifyPage{}

	if response.StatusCode == 200 {
		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Println(readErr)
		}

		jsonErr := json.Unmarshal(body, &firstPage)

		if jsonErr != nil {
			log.Println(jsonErr)
		}
	}
	return firstPage
}

func processArtists(page spotifyPage) []dal.Artist {
	artistList := []dal.Artist{}

	for _, album := range page.Items {
		for _, artist := range album.Album.Artists {
			newArtist := &dal.Artist{
				Name:      artist.Name,
				SpotifyID: artist.ID,
			}
			artistList = append(artistList, *newArtist)
		}
	}

	return artistList

}
