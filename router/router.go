package router

import (
	"htmx/service"
	"log"
	"net/http"
)

func InitRouter() {

	http.HandleFunc("/", service.GetHome)
	http.HandleFunc("/increment-dog", service.IncrementDog)

	fs := http.FileServer(http.Dir("template"))
	http.Handle("/template/", http.StripPrefix("/template/", fs))

	log.Println("running server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
