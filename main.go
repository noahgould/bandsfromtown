package main

import (
	"github.com/noahgould/bandsfromtown/dal"
)

func main() {
	db, err := dal.StartDB("noah:bigDBpass@localhost/bandsfromtown")
}
