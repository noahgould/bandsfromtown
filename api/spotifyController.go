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

	artistChan := make(chan dal.Artist)

	go func() {
		for _, a := range artistList {
			artistChan <- a
		}
		fmt.Println("close artistlist to artist chan")
		close(artistChan)
	}()

	//start requesting playlists
	// playlistResultPage := makePlaylistRequest(userToken, 0)
	// playlists := playlistResultPage.Playlists
	// //get all the playlists.
	// for numPlaylists := 50; numPlaylists <= playlistResultPage.Total; numPlaylists += 50 {
	// 	playlists = append(playlists, makePlaylistRequest(userToken, numPlaylists).Playlists...)
	// }

	// for _, playlist := range playlists {
	// 	trackOffset := 0

	// 	trackPage := makePlaylistTrackRequest(userToken, trackOffset, playlist)
	// 	for _, track := range trackPage.PlaylistTracks {
	// 		artistList = append(artistList, spotifyArtistToArtist(artistMap, track.Track.Artists...)...)
	// 	}

	// 	for trackOffset := 100; trackOffset <= trackPage.Total; trackOffset += 100 {
	// 		trackPage = makePlaylistTrackRequest(userToken, trackOffset, playlist)
	// 		for _, track := range trackPage.PlaylistTracks {
	// 			artistList = append(artistList, spotifyArtistToArtist(artistMap, track.Track.Artists...)...)
	// 		}
	// 	}
	// }

	artistList = sc.getArtistLocations(artistChan)
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

func (sc *SpotifyController) checkSavedWithSpotify(artists <-chan dal.Artist, readyArtists chan<- dal.Artist) <-chan dal.Artist {
	notSaved := make(chan dal.Artist)

	fmt.Println("begin check saved spotify")
	go func() {
		for artist := range artists {
			existingArtist, err := sc.artistStore.GetArtistBySpotifyID(artist.SpotifyID)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Println(err)
				} else {
					notSaved <- artist
				}
			} else {
				fmt.Println(existingArtist)
				readyArtists <- existingArtist
			}
		}
		close(notSaved)
		fmt.Println("end check saved spotify func")
	}()

	return notSaved
}

func (sc *SpotifyController) checkSavedByName(artists <-chan dal.Artist, readyArtists chan<- dal.Artist) <-chan dal.Artist {
	notSavedArtists := make(chan dal.Artist)

	fmt.Println("begin check saved name")
	go func() {
		for artist := range artists {
			existingArtists, err := sc.artistStore.GetArtistsByName(artist.Name)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Println(err)
				} else {
					notSavedArtists <- artist
				}
			} else {
				readyArtists <- existingArtists[0]
			}
		}
		close(notSavedArtists)
		close(readyArtists)
		fmt.Println("end check saved name func")
	}()

	return notSavedArtists
}

func (sc *SpotifyController) getArtistLocations(artists <-chan dal.Artist) []dal.Artist {
	gmc := NewGoogleMapsController()

	readyArtists := make(chan dal.Artist)

	noSpotifyArtists := sc.checkSavedWithSpotify(artists, readyArtists)
	notSavedArtists := sc.checkSavedByName(noSpotifyArtists, readyArtists)

	fmt.Println("get artistLocations after calling funcs.")
	artistList := []dal.Artist{}

	for artist := range notSavedArtists {
		artistList = append(artistList, artist)
	}

	fmt.Println("after not saved artists")

	fmt.Println(artistList[0])

	readyArtistList := []dal.Artist{}
	for artist := range readyArtists {
		readyArtistList = append(readyArtistList, artist)
	}

	for i, artist := range artistList {
		artistList[i].Location = LookupArtistLocation(artist.Name)
		locationPtr, err := gmc.NormalizeLocation(artistList[i].Location)
		if err != nil {
			log.Println(err)
		}
		artistList[i].Location = *locationPtr
		var exists bool
		exists, artistList[i].Location = sc.locationStore.CheckForExistingLocation(artistList[i].Location)
		if !exists {
			if artistList[i].Location.Longitude == 0 && artistList[i].Location.GooglePlaceID != "-1" {
				locationPointer, err := gmc.GetCoordinates(artistList[i].Location)
				if err != nil {
					log.Println(err)
				}
				artistList[i].Location = *locationPointer
			}
			artistList[i].Location.ID, err = sc.locationStore.AddLocation(artistList[i].Location)
			if err != nil {
				log.Println(err)
			}

			artistList[i].ID, err = sc.artistStore.AddArtist(artistList[i])
			if err != nil {
				log.Println(err)
			}
		} else {
			artistList[i].ID, err = sc.artistStore.AddArtist(artistList[i])
			if err != nil {
				log.Println(err)
			}
		}
	}

	return artistList
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
