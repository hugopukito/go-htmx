package service

import (
	"fmt"
	"html/template"
	"htmx/entity"
	"htmx/repository"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))

	dogs, err := repository.FindDogs()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var dogTmpls []entity.DogTmpl

	for _, dog := range dogs {
		dogTmpl := entity.DogTmpl{
			Dog:  dog,
			Date: dog.DateCreation.Format("15:04:05"),
		}
		dogTmpls = append(dogTmpls, dogTmpl)
	}

	dogsHtmx := map[string][]entity.DogTmpl{
		"Dogs": dogTmpls,
	}

	tmpl.Execute(w, dogsHtmx)
}

func IncrementDog(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	tmpl, _ := template.New("t").Parse("1")
	tmpl.Execute(w, nil)
}
