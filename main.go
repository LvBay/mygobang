package main

import (
	"log"
	_ "mygobang/controllers"
	_ "mygobang/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
