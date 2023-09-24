package service

import (
	"fmt"
	"html/template"
	"htmx/entity"
	"htmx/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))

	dogs, err := repository.FindDogs()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// var dogTmpls []entity.DogTmpl

	// for _, dog := range dogs {
	// 	dogTmpl := entity.DogTmpl{
	// 		Dog:  dog,
	// 		Date: dog.DateCreation.Format("15:04:05"),
	// 	}
	// 	dogTmpls = append(dogTmpls, dogTmpl)
	// }

	dogsHtmx := map[string][]entity.Dog{
		"Dogs": dogs,
	}

	tmpl.Execute(w, dogsHtmx)
}

func IncrementDog(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	uuid, err := uuid.Parse(id)
	if err != nil {
		// will remove score from dom, need to do nothing
		fmt.Println(err)
		return
	}

	score, err := repository.IncrementScore(uuid)
	if err != nil {
		// will remove score from dom, need to do nothing
		fmt.Println(err)
		return
	}

	time.Sleep(2 * 1000)

	tmpl, _ := template.New("t").Parse(strconv.Itoa(score))
	tmpl.Execute(w, nil)
}
