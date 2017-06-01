package router

import (
	"fmt"
	"mygobang/controllers"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/login", controllers.LoginHandler) //设置访问的路由
	http.HandleFunc("/gobang", controllers.GobangHandler)
	http.HandleFunc("/gobangws", controllers.WsHandler)
	CheckStatic(http.DefaultServeMux, "/static/", "./static/")
}
func CheckStatic(mux *http.ServeMux, prefix, staticDir string) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.Replace(r.URL.Path, prefix, staticDir, 1)
		fmt.Println("filePath:", filePath)
		http.ServeFile(w, r, filePath)
	})
}
