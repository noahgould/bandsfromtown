package api

import (
	"html/template"
	"log"
	"net/http"
)

type HomepageController struct {
}

func NewHomepageController() *HomepageController {

	return &HomepageController{}
}

func (hc *HomepageController) Index(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("./frontend/homepage.html")
	if err != nil {
		log.Println(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func (hc *HomepageController) AboutPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./frontend/aboutpage.html")

	if err != nil {
		log.Println(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}
