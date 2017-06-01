package controllers

import (
	// "gobang"
	"encoding/json"
	// "fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	IsFirst = true
	rooms   = map[string]*Room{}
)

// func Loop(n int) []int {
// 	slice := make([]int, n)
// 	for i := 0; i < n; i++ {
// 		slice[i] = i + 1
// 	}
// 	return slice
// }
// func Add(n int) int {
// 	n++
// 	return n
// }

func init() {
	rooms["110"] = &Room{board: &Board{}}
	rooms["111"] = &Room{board: &Board{}}
}

// gobang
func GobangHandler(w http.ResponseWriter, r *http.Request) {
	roomid := r.FormValue("roomid")
	if roomid == "" || rooms[roomid] == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	d := map[string]string{"sign": roomid}
	t, _ := template.ParseFiles("views/board.html")
	t.Execute(w, d)
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(time.Now(), err)
	}
	go wsListener(ws) // 监听
}

func wsListener(ws *websocket.Conn) {
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		// 解析客户端数据
		req := &Req{}
		if err := json.Unmarshal(p, req); err != nil {
			log.Printf("json.Unmarshal err:%v\n", err)
			return
		}
		// 检查房间是否存在
		if _, exist := rooms[req.Sign]; !exist {
			return
		}
		room := rooms[req.Sign]

		if req.Act == 0 && req.Role == "" { // 登录
			req.Role = room.Join(ws) // 分配执手颜色
			data, _ := json.Marshal(req)
			ws.WriteMessage(websocket.TextMessage, data)
		} else if req.Act == 1 { // 点击
			x, y, _ := req.OptInfo()
			room.board.Fill(x, y, req.Role) // 落子
			// 检查胜负
			cl := TurnColor(req.Role)
			bd := room.board
			if bd.checkWin(x, y, cl) {
				req.Over = true
				room.ClearBoard()
			}
			// 广播结果
			data, _ := json.Marshal(req)
			room.Broad(data, req.Over, req.Role)

		} else if req.Act == 2 { // 连接断开
			// todo
		}

	}
}

func TurnBool(b bool) bool {
	return !b
}

func TurnColor(s string) int {
	m := map[string]int{"white": 1, "black": -1}
	return m[s]
}
