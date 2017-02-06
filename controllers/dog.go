package controllers

import (
	"html/template"
	"net/http"
)

func ViewsHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/little-dogs.html")
	t.Execute(w, nil)
}
