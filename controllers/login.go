package controllers

import (
	// "fmt"
	// "gobang"
	// "html/template"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// fmt.Println("method:", r.Method)
		// t, _ := template.ParseFiles("views/login.html")
		// t.Execute(w, nil)
		w.Write([]byte("123"))
	} else if r.Method == "POST" {
		uid := r.FormValue("uid")
		passwd := r.FormValue("passwd")
		if uid == "zyh" && passwd == "ww" {
			http.Redirect(w, r, "/gabang", 302)
		} else {
			w.Write([]byte("嘿嘿嘿"))
		}
	}
}
