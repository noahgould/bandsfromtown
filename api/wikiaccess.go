package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Artist struct {
	Name    string
	City    string
	State   string
	Country string
}

type Page struct {
	PageID    int    `json:"pageid"`
	Ns        int    `json:"ns"`
	Title     string `json:"title"`
	Revisions []struct {
		Contentformat string `json:"contentformat"`
		Contentmodel  string `json:"contentmodel"`
		Content       string `json:"*"`
	}
}

type WikiApiResult struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]*Page
	}
}

func main() {
	artist := "The%20Shins"
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

	pageInfo := WikiApiResult{}
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
	fmt.Println(location)

}
