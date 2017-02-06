package controllers

import (
	// "gobang"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	// "sync"
	"container/list"
	"time"
)

const (
	staticDir = "./static"
	White     = 1
	Black     = 2
)

var (
	IsFirst = true
	rooms   = map[string]*Room{}
	// board   = [14][16]int{}
)

func Loop(n int) []int {
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		slice[i] = i + 1
	}
	return slice
}
func Add(n int) int {
	n++
	return n
}

type Req struct {
	Role     string // 执手
	Act      int    // 行为
	Position string // 位置
	Over     bool
	Sign     string
}
type Room struct {
	// wsList []*websocket.Conn
	wsList *list.List
	board  [15][17]int
}

func init() {
	room := &Room{}
	room.wsList = list.New()
	rooms["110"] = room

	_, exist := rooms["110"]
	fmt.Println("exist:", exist)
}

// gobang
func GobangHandler(w http.ResponseWriter, r *http.Request) {
	roomid := r.FormValue("roomid")
	if roomid == "" {
		w.Write([]byte("请先创建房间"))
		return
	}
	d := map[string]string{"sign": roomid}
	t, _ := template.ParseFiles("views/game2.html")
	t.Execute(w, d)
}
func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(time.Now(), err)
	}
	// wsList = append(wsList, ws)
	// defer ws.Close()
	go roomListener(ws) // 如果不用go关键字怎么实现既可以for循环监听，又能转换IsFirst？
	// 下一个进入房间的用户颜色改变
	IsFirst = TurnBool(IsFirst)
}

func roomListener(ws *websocket.Conn) {
	for {

		_, p, err := ws.ReadMessage()
		if err != nil {

			return
		}

		// // 解析客户端数据
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
			// room.wsList = append(rooms[req.Sign].wsList, ws)
			room.wsList.PushBack(ws)

			req.Role = CheckWtOrBk() // 分配执手颜色
			data, _ := json.Marshal(req)
			ws.WriteMessage(websocket.TextMessage, data)
		} else if req.Act == 1 { // 点击
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), room.wsList.Len())
			arr := strings.Split(req.Position, "_")
			if len(arr) != 2 {
				return
			} // 请求信息非法
			x, _ := strconv.Atoi(arr[0])
			y, _ := strconv.Atoi(arr[1])

			rooms[req.Sign].board[x][y] = TurnColor(req.Role)
			// 检查胜负
			if req.checkWinX() || req.checkWinY() || req.checkWin13() || req.checkWin24() {
				req.Over = true
				fmt.Println(req.Role, "win!")
				rooms[req.Sign].ClearBoard()
			}
			// 若胜负已分，发给所有人，否则发给除了自己以外的其他人
			for conn := room.wsList.Front(); conn != nil; conn = conn.Next() {
				data, _ := json.Marshal(req)
				conn.Value.(*websocket.Conn).WriteMessage(websocket.TextMessage, data)
			}

		} else if req.Act == 2 { // 连接断开

			for conn := room.wsList.Front(); conn != nil; conn = conn.Next() {
				if conn.Value.(*websocket.Conn) == ws {
					fmt.Println("len1:", room.wsList.Len())
					fmt.Println("断开连接")
					room.wsList.Remove(conn)
					fmt.Println("len2:", room.wsList.Len())
				}
			}
		}

	}
}

func TurnBool(b bool) bool {
	if b {
		return false
	} else {
		return true
	}
}
func CheckWtOrBk() string {
	if IsFirst {
		return "white"
	}
	return "black"
}

func TurnColor(s string) int {
	if s == "white" {
		return 1
	}
	if s == "black" {
		return -1
	}
	return 0
}
func (r *Req) checkWinX() bool {
	board := rooms[r.Sign].board
	count := 0
	arr := strings.Split(r.Position, "_")
	x, _ := strconv.Atoi(arr[0])
	y, _ := strconv.Atoi(arr[1])
	for i := y - 1; i >= 0; i-- {
		if board[x][i] == TurnColor(r.Role) {
			count++
		} else {
			break
		}
	}
	for i := y + 1; i <= 16; i++ {
		if board[x][i] == TurnColor(r.Role) {
			count++
		} else {
			break
		}
	}
	return count >= 4

}

func (r *Req) checkWinY() bool {
	board := rooms[r.Sign].board
	count := 0
	arr := strings.Split(r.Position, "_")
	x, _ := strconv.Atoi(arr[0])
	y, _ := strconv.Atoi(arr[1])
	for i := x - 1; i >= 0; i-- {
		if board[i][y] == TurnColor(r.Role) {
			count++
		} else {
			break
		}
	}
	for i := x + 1; i <= 14; i++ {
		if board[i][y] == TurnColor(r.Role) {
			count++
		} else {
			break
		}
	}
	return count >= 4
}
func (r *Req) checkWin13() bool {
	board := rooms[r.Sign].board
	count := 0
	arr := strings.Split(r.Position, "_")
	x, _ := strconv.Atoi(arr[0])
	y, _ := strconv.Atoi(arr[1])
	for i, j := x-1, y+1; i >= 0 && j <= 16; {
		if board[i][j] == TurnColor(r.Role) {
			count++
			i--
			j++
		} else {
			break
		}
	}
	for i, j := x+1, y-1; i <= 14 && j >= 0; i++ {
		if board[i][j] == TurnColor(r.Role) {
			count++
			i++
			j--
		} else {
			break
		}
	}
	return count >= 4
}
func (r *Req) checkWin24() bool {
	board := rooms[r.Sign].board
	count := 0
	arr := strings.Split(r.Position, "_")
	x, _ := strconv.Atoi(arr[0])
	y, _ := strconv.Atoi(arr[1])
	for i, j := x-1, y-1; i >= 0 && j >= 0; {
		if board[i][j] == TurnColor(r.Role) {
			count++
			i--
			j--
		} else {
			break
		}
	}
	for i, j := x+1, y+1; i <= 14 && j <= 16; {
		if board[i][j] == TurnColor(r.Role) {
			count++
			i++
			j++
		} else {
			break
		}
	}
	return count >= 4
}

func (r *Room) ClearBoard() {
	for i := 0; i <= 14; i++ {
		for j := 0; j <= 16; j++ {
			r.board[i][j] = 0
		}
	}
}
