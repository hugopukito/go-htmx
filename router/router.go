package router

import (
	"htmx/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRouter() {

	r := mux.NewRouter()

	r.HandleFunc("/", service.GetHome)
	r.HandleFunc("/increment-dog/{id}", service.IncrementDog)
	r.HandleFunc("/ws", service.HandleWsConnection)

	fs := http.FileServer(http.Dir("template"))
	r.PathPrefix("/template/").Handler(http.StripPrefix("/template/", fs))

	log.Println("running server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
