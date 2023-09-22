package main

import (
	"htmx/repository"
	"htmx/router"

	fixture "github.com/hugopukito/golang-fixture"
	"github.com/hugopukito/golang-fixture/database"
)

func main() {
	fixture.RunFixtures("fixtures", database.DatabaseParams{Name: "htmx"})
	repository.InitDB()
	router.InitRouter()
}
