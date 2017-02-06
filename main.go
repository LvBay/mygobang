package main

import (
	"log"
	_ "myweb/controllers"
	_ "myweb/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//hahahah
}
