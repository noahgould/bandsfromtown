package api

import (
	"net/http"
)

func LookupArtist(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}
