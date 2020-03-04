package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Xeoncross/got"
)

type Page struct {
	Title string
	Name  string
}

func main() {

	templates, err := got.New("templates", ".html", got.DefaultFunctions)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Page{Title: "Home", Name: "Guest"}
		err := templates.Render(w, "home", data, http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		data := Page{Title: "About", Name: "Guest"}
		err := templates.Render(w, "about", data, http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
