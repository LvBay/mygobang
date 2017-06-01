package controllers

import (
	"html/template"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		roomId := r.FormValue("roomId")
		if _, exist := rooms[roomId]; exist {
			w.Write([]byte("房间已创建"))
		} else {
			rooms[roomId] = &Room{board: &Board{}}
			http.Redirect(w, r, "/gobang?roomid="+roomId, 302)
		}
	}
}
