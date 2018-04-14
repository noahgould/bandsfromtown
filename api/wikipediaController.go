package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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
	return strings.Replace(artistName, " ", "%20", -1)
}

//LookupArtistLocation queries wikipedia for an artists location and returns it.
func LookupArtistLocation(artist string) string {
	url := "http://en.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&format=json&titles=" + artist + "&rvsection=0"

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
	for _, p := range pageInfo.Query.Pages {
		if strings.Contains(p.Revisions[0].Content, "origin =") {
			infoBoxItems = strings.Split(p.Revisions[0].Content, "origin =")
		} else if strings.Contains(p.Revisions[0].Content, "birth_place =") {
			infoBoxItems = strings.Split(p.Revisions[0].Content, "birth_place =")
		} else {
			fmt.Println("No birthplace/origin data")
		}
	}
	endOfLocation := strings.Index(infoBoxItems[1], "|")
	location := strings.Trim(infoBoxItems[1][0:endOfLocation], " ")
	location = strings.Replace(location, "]", "", -1)
	location = strings.Replace(location, "[", "", -1)
	location = strings.Replace(location, "&nbsp;", " ", -1)
	return location
}
