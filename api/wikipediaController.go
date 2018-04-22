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

func wikipediaFormat(artistName string) string {
	artistName = strings.Replace(artistName, " ", "%20", -1)
	return strings.Replace(artistName, "_", "%20", -1)
}

func locationStringToStruct(location string) dal.Location {
	fmt.Println(location)
	locationPieces := strings.Split(location, ",")
	for i := range locationPieces {
		if strings.Contains(locationPieces[i], "(") {
			locationPieces[i] = strings.Split(locationPieces[i], "(")[0]
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

	url := "http://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&format=json&titles=" + wikipediaFormat(artist) + "&rvsection=0"
	wikiClient := http.Client{
		Timeout: time.Second * 2,
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
		if strings.Contains(p.Revisions[0].Content, "origin") {
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

	if originGood || birthPlaceGood {
		endOfLocation := strings.Index(infoBoxItems[1], "| ")
		location := strings.Trim(infoBoxItems[1][0:endOfLocation], "")
		location = strings.Replace(location, "]", "", -1)
		location = strings.Replace(location, "[", "", -1)
		location = strings.Replace(location, "&nbsp;", " ", -1)
		location = strings.Replace(location, "\n", "", -1)
		location = strings.Replace(location, "nowrap", "", -1)
		return locationStringToStruct(location)
	}

	return locationStringToStruct("null, null, null")

}
