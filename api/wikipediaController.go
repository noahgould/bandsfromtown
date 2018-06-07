package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/noahgould/bandsfromtown/dal"
)

var infoboxStrings = []string{"genre", "years_active", "alias", "origin", "occupation"}

type page struct {
	PageID    int    `json:"pageid"`
	Ns        int    `json:"ns"`
	Title     string `json:"title"`
	Revisions []struct {
		Contentformat string `json:"contentformat"`
		Contentmodel  string `json:"contentmodel"`
		Content       string `json:"*"`
	}
}

type wikiApiResult struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]*page
	}
}

type wikiSearchResult struct {
	Query        string
	PageTitles   []string
	PageSnippets []string
	PageURLS     []string
}

func wikipediaFormat(artistName string) string {
	artistName = strings.Replace(artistName, " ", "%20", -1)
	artistName = strings.Replace(artistName, "&", "%26", -1)
	return artistName
}

func (w *wikiSearchResult) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&w.Query, &w.PageTitles, &w.PageSnippets, &w.PageURLS}
	correctLen := len(tmp)
	if jsonErr := json.Unmarshal(buf, &tmp); jsonErr != nil {
		return jsonErr
	}
	if g, e := len(tmp), correctLen; g != e {
		return fmt.Errorf("Wrong number of fields for search result")
	}
	return nil
}

func locationStringToStruct(location string) dal.Location {
	locationPieces := strings.Split(location, ",")
	for i := range locationPieces {
		if strings.Contains(locationPieces[i], "(") {
			locationPieces[i] = strings.Split(locationPieces[i], "(")[0]
		}
		if strings.Contains(locationPieces[i], "|") {
			locationPieces[i] = strings.Split(locationPieces[i], "|")[0]
		}

		locationPieces[i] = strings.TrimSpace(locationPieces[i])
	}

	var state, country, city string
	if len(locationPieces) >= 3 {
		city = locationPieces[len(locationPieces)-3]
		state = locationPieces[len(locationPieces)-2]
		country = locationPieces[len(locationPieces)-1]
	} else {
		city = locationPieces[0]
		state = "unknown"
		country = locationPieces[1]
	}

	return dal.Location{
		City:         city,
		State:        state,
		Country:      country,
		FullLocation: location,
	}
}

//LookupArtistLocation queries wikipedia for an artists location and returns it.
func LookupArtistLocation(artist string) dal.Location {

	log.Println(artist)

	url := "http://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&format=json&titles=" + wikipediaFormat(artist) + "&rvsection=0"
	wikiClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "bandsfromtown")

	res, getErr := wikiClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	pageInfo := wikiApiResult{}
	jsonErr := json.Unmarshal(body, &pageInfo)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	var infoBoxItems []string
	originGood := false
	birthPlaceGood := false
	for _, p := range pageInfo.Query.Pages {
		if len(p.Revisions) > 0 {
			if strings.Contains(p.Revisions[0].Content, "Infobox musical artist") || strings.Contains(p.Revisions[0].Content, "Infobox person") {
				if strings.Contains(p.Revisions[0].Content, "origin ") {
					infoBoxItems = strings.Split(p.Revisions[0].Content, "origin")
					infoBoxItems = strings.Split(infoBoxItems[1], "=")
					if strings.Count(infoBoxItems[1], ",") > 0 {
						originGood = true
					}
				}
				if !originGood {
					if strings.Contains(p.Revisions[0].Content, "birth_place") {
						infoBoxItems = strings.Split(p.Revisions[0].Content, "birth_place")
						infoBoxItems = strings.Split(infoBoxItems[1], "=")
						if strings.Count(infoBoxItems[1], ",") > 0 {
							birthPlaceGood = true
						}
					}
				}
			}
		}
	}

	if originGood || birthPlaceGood {
		return getLocationFromResult(infoBoxItems)
	}

	newSearchArtist := searchForPage(wikipediaFormat(artist))
	if newSearchArtist == artist {
		return locationStringToStruct("nil, nil, nil")
	}
	return LookupArtistLocation(newSearchArtist)

}

func getLocationFromResult(infoBox []string) dal.Location {

	endOfLocation := 0
	for _, infoboxCategory := range infoboxStrings {
		wordIndex := strings.Index(infoBox[1], infoboxCategory)
		if wordIndex != -1 {
			endOfLocation = wordIndex
			break
		}
	}

	var location string
	if endOfLocation > 0 {
		location = strings.Trim(infoBox[1][0:endOfLocation], "")
	} else {
		location = infoBox[1]
	}

	locationBits := strings.Split(location, ",")

	for i, bit := range locationBits {
		endOfLocation = strings.LastIndex(bit, "|")
		if endOfLocation > 0 {
			bit = strings.Trim(bit[0:endOfLocation], "")
		}
		bit = strings.Replace(bit, "]", "", -1)
		bit = strings.Replace(bit, "[", "", -1)
		bit = strings.Replace(bit, "&nbsp;", " ", -1)
		bit = strings.Replace(bit, "\n", "", -1)
		locationBits[i] = strings.Replace(bit, "nowrap", "", -1)
	}

	location = strings.Join(locationBits, ",")

	log.Println(location)
	return locationStringToStruct(location)

}

func searchForPage(artist string) string {
	url := "https://en.wikipedia.org/w/api.php?action=opensearch&search=" + artist + "&namespace=0&format=json"
	wikiClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "bandsfromtown")

	res, getErr := wikiClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var searchResult wikiSearchResult

	jsonErr := json.Unmarshal(body, &searchResult)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	for _, title := range searchResult.PageTitles {
		if strings.Contains(title, "band") || strings.Contains(title, "musical artist") || strings.Contains(title, "musician") || strings.Contains(title, "singer") {
			return title
		}
	}
	if len(searchResult.PageTitles) > 0 {
		artist = searchResult.PageTitles[0]
	}

	return artist

}
