package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
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
	AccessToken    string `json:"access_token"`
	TokenType      string `json:"token_type"`
	Scope          string `json:"scope"`
	ExpirationTime int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
}

type spotifyAlbum struct {
	AlbumType            string                `json:"album_type"`
	Artists              []spotifySimpleArtist `json:"artists"`
	AvailableMarkets     []string              `json:"available_markets"`
	Copyrights           []copyright           `json:"copyrights"`
	ExternalIds          []string              `json:"-"`
	ExternalUrls         []externalURL         `json:"-"`
	Genres               []string              `json:"genres"`
	Href                 string                `json:"href"`
	ID                   string                `json:"id"`
	Images               []image               `json:"images"`
	Label                string                `json:"label"`
	Name                 string                `json:"name"`
	Popularity           int                   `json:"popularity"`
	ReleaseDate          string                `json:"release_date"`
	ReleaseDatePrecision string                `json:"release_date_precision"`
	Restrictions         []string              `json:"-"`
	Tracks               []string              `json:"-"`
	ObjectType           string                `json:"type"`
	URI                  string                `json:"uri"`
}

type spotifySimpleArtist struct {
	ExternalUrls []externalURL `json:"-"`
	Href         string        `json:"href"`
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	ObjectType   string        `json:"type"`
	URI          string        `json:"uri"`
}

type spotifySimplePlaylist struct {
	Collaborative bool          `json:"collaborative"`
	Href          string        `json:"href"`
	ExternalUrls  []externalURL `json:"-"`
	ID            string        `json:"id"`
	Images        []image       `json:"images"`
	Name          string        `json:"name"`
	Owner         spotifyUser   `json:"owner"`
	Public        bool          `json:"public"`
	SnapshotID    string        `json:"snapshot_id"`
	Tracks        spotifyTracks `json:"tracks"`
	ObjectType    string        `json:"type"`
	URI           string        `json:"uri"`
}

type spotifyPlaylistTrack struct {
	AddedAt   string       `json:"added_at"`
	AddedBy   spotifyUser  `json:"added_by"`
	LocalFile bool         `json:"is_local"`
	Track     spotifyTrack `json:"track"`
}

type spotifyTrack struct {
	Album            spotifySimpleAlbum    `json:"album"`
	Artists          []spotifySimpleArtist `json:"artists"`
	AvailableMarkets []string              `json:"available_markets"`
	DiscNumber       int                   `json:"disc_number"`
	DurationMS       int                   `json:"duration_ms"`
	Explicit         bool                  `json:"explicit"`
	ExternalID       []string              `json:"-"`
	ExternalUrls     []externalURL         `json:"-"`
	Href             string                `json:"href"`
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Popularity       int                   `json:"popularity"`
	PreviewURL       string                `json:"preview_url"`
	TrackNumber      int                   `json:"track_number"`
	ObjectType       string                `json:"type"`
	URI              string                `json:"uri"`
}

type spotifySimpleAlbum struct {
	AlbumType            string                `json:"album_type"`
	Artists              []spotifySimpleArtist `json:"artists"`
	AvailableMarkets     []string              `json:"available_markets"`
	ExternalUrls         []externalURL         `json:"-"`
	Href                 string                `json:"href"`
	ID                   string                `json:"id"`
	Images               []image               `json:"images"`
	Name                 string                `json:"name"`
	ReleaseDate          string                `json:"release_date"`
	ReleaseDatePrecision string                `json:"release_date_precision"`
	Restrictions         []string              `json:"-"`
	ObjectType           string                `json:"type"`
	URI                  string                `json:"uri"`
}

type spotifyUser struct {
	DisplayName  string        `json:"display_name"`
	ExternalUrls []externalURL `json:"-"`
	Followers    string        `json:"-"`
	Href         string        `json:"href"`
	ID           string        `json:"id"`
	Images       []image       `json:"images"`
	ObjectType   string        `json:"type"`
	URI          string        `json:"uri"`
}

// type spotifyExternalID {
// 	Key string `j`
// 	Value string `json:"-"`
// }

type spotifyTracks struct {
	TracksURI      string `json:"href"`
	NumberOfTracks int    `json:"total"`
}

type image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type externalURL struct {
	Type  string
	Value string
}

type copyright struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type spotifyBasePage struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

type spotifyAlbumPage struct {
	spotifyBasePage
	Albums []savedAlbum `json:"items"`
}

type spotifyPlaylistPage struct {
	spotifyBasePage
	Playlists []spotifySimplePlaylist `json:"items"`
}

type spotifyTrackPage struct {
	spotifyBasePage
	PlaylistTracks []spotifyPlaylistTrack `json:"items"`
}

type savedAlbum struct {
	AddedAt string       `json:"added_at"`
	Album   spotifyAlbum `json:"album"`
}

const redirectURI string = "https://bandsfromtown.herokuapp.com/spotify/login/"

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

	render(w, "./frontend/spotifyResults.html", tokenResult.AccessToken)

}

func render(w http.ResponseWriter, tmpl string, arg string) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, arg)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func (sc *SpotifyController) FindUserArtistLocations(w http.ResponseWriter, r *http.Request) {
	accessToken := mux.Vars(r)["accessToken"]

	usersArtists := sc.getAllUserArtists(accessToken)

	if err := json.NewEncoder(w).Encode(usersArtists); err != nil {
		log.Println(err)
	}
}

func (sc *SpotifyController) getAllUserArtists(userToken string) []dal.Artist {

	//start requesting user albums
	resultPage := makeAlbumRequest(userToken, 0)
	artistList := []dal.Artist{}
	artistMap := make(map[string]bool)
	artistList = append(artistList, getArtistsFromAlbums(resultPage, artistList, artistMap)...)
	//go through each page of user albums, getting artists from albums as we go along.
	for numAlbums := 50; numAlbums <= resultPage.Total; numAlbums += 50 {
		resultPage = makeAlbumRequest(userToken, numAlbums)
		artistList = getArtistsFromAlbums(resultPage, artistList, artistMap)
	}

	//start requesting playlists
	playlistResultPage := makePlaylistRequest(userToken, 0)
	playlists := playlistResultPage.Playlists
	//get all the playlists.
	for numPlaylists := 50; numPlaylists <= playlistResultPage.Total; numPlaylists += 50 {
		playlists = append(playlists, makePlaylistRequest(userToken, numPlaylists).Playlists...)
	}

	for _, playlist := range playlists {
		trackOffset := 0

		trackPage := makePlaylistTrackRequest(userToken, trackOffset, playlist)
		for _, track := range trackPage.PlaylistTracks {
			artistList = append(artistList, spotifyArtistToArtist(artistMap, track.Track.Artists...)...)
		}

		for trackOffset := 100; trackOffset <= trackPage.Total; trackOffset += 100 {
			trackPage = makePlaylistTrackRequest(userToken, trackOffset, playlist)
			for _, track := range trackPage.PlaylistTracks {
				artistList = append(artistList, spotifyArtistToArtist(artistMap, track.Track.Artists...)...)
			}
		}
	}

	artistList = sc.getArtistLocations(artistList)
	return artistList
}

func spotifyArtistToArtist(artistMap map[string]bool, artist ...spotifySimpleArtist) []dal.Artist {
	newArtists := []dal.Artist{}
	for _, a := range artist {
		if _, ok := artistMap[a.ID]; !ok {
			artistMap[a.ID] = true
			newArtist := dal.Artist{
				Name:      a.Name,
				SpotifyID: a.ID,
			}
			newArtists = append(newArtists, newArtist)
		}
	}
	return newArtists
}

// can make this more efficient by limiting results. https://beta.developer.spotify.com/documentation/web-api/reference/playlists/get-playlists-tracks/
func makePlaylistTrackRequest(userToken string, offset int, playlist spotifySimplePlaylist) spotifyTrackPage {

	spotifyClient := &http.Client{
		Timeout: time.Second * 5,
	}
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists/%s/tracks", playlist.Owner.ID, playlist.ID)

	req, err := http.NewRequest("GET", requestURL, nil)
	req.Header.Add("Authorization", "Bearer "+userToken)

	q := req.URL.Query()
	q.Add("limit", "100")
	q.Add("offset", strconv.Itoa(offset))

	req.URL.RawQuery = q.Encode()

	response, err := spotifyClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	firstPage := spotifyTrackPage{}

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

func (sc *SpotifyController) getArtistLocations(artists []dal.Artist) []dal.Artist {
	gmc := NewGoogleMapsController()

	var err error
	stopQueryingGoogle := false
	for i, artist := range artists {
		var existingArtist dal.Artist
		existingArtist, err = sc.artistStore.GetArtistBySpotifyID(artist.SpotifyID)
		if err != nil {
			if err == sql.ErrNoRows {
				possibleArtists, err := sc.artistStore.GetArtistsByName(artist.Name)
				if len(possibleArtists) == 0 {
					if err == nil {
						if !stopQueryingGoogle {
							artists[i].Location = LookupArtistLocation(artist.Name)
							locationPtr, err := gmc.NormalizeLocation(artists[i].Location)
							if err != nil {
								if err.Error() == "maps: OVER_QUERY_LIMIT - You have exceeded your daily request quota for this API.." {
									stopQueryingGoogle = true
								} else {
									log.Println(err)
								}
							}
							artists[i].Location = *locationPtr
							var exists bool
							exists, artists[i].Location = sc.locationStore.CheckForExistingLocation(artists[i].Location)
							if !exists {
								if artists[i].Location.Longitude == 0 && artists[i].Location.GooglePlaceID != "-1" {
									locationPointer, err := gmc.GetCoordinates(artists[i].Location)
									if err != nil {
										log.Println(err)
									}
									artists[i].Location = *locationPointer
								}
								artists[i].Location.ID, err = sc.locationStore.AddLocation(artists[i].Location)
								if err != nil {
									log.Println(err)
								}

								artists[i].ID, err = sc.artistStore.AddArtist(artists[i])
								if err != nil {
									log.Println(err)
								}

							} else {
								sc.artistStore.AddArtist(artists[i])
							}
						}
					} else {
						log.Println(err)
					}
				} else {
					artists[i] = possibleArtists[0]
					artists[i].Location, err = sc.locationStore.GetLocationByID(possibleArtists[0].Location.ID)
					if err != nil {
						log.Println(err)
					}
				}
			} else {
				log.Println(err)
			}
		} else {
			artists[i] = existingArtist
			artists[i].Location, err = sc.locationStore.GetLocationByID(existingArtist.Location.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return artists
}

func makeAlbumRequest(userToken string, offset int) spotifyAlbumPage {

	spotifyClient := &http.Client{
		Timeout: time.Second * 20,
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/albums", nil)

	req.Header.Add("Authorization", "Bearer "+userToken)

	q := req.URL.Query()
	q.Add("limit", "50")
	q.Add("offset", strconv.Itoa(offset))

	req.URL.RawQuery = q.Encode()

	response, err := spotifyClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	firstPage := spotifyAlbumPage{}

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

func makePlaylistRequest(userToken string, offset int) spotifyPlaylistPage {
	spotifyClient := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)

	req.Header.Add("Authorization", "Bearer "+userToken)

	q := req.URL.Query()
	q.Add("limit", "50")
	q.Add("offset", strconv.Itoa(offset))

	req.URL.RawQuery = q.Encode()

	response, err := spotifyClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	firstPage := spotifyPlaylistPage{}

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

func getArtistsFromAlbums(page spotifyAlbumPage, artistList []dal.Artist, artistMap map[string]bool) []dal.Artist {

	for _, savedAlbum := range page.Albums {
		for _, artist := range savedAlbum.Album.Artists {
			if _, ok := artistMap[artist.ID]; !ok {
				artistMap[artist.ID] = true
				newArtist := &dal.Artist{
					Name:      artist.Name,
					SpotifyID: artist.ID,
				}
				artistList = append(artistList, *newArtist)
			}
		}
	}
	return artistList

}
