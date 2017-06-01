package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
)

const (
	White = 1
	Black = 2
)

type Room struct {
	// wsList *list.List
	W     *Role   // 白手
	B     *Role   // 黑手
	OB    []*Role // 吃瓜群众
	board *Board
}

func (r *Room) Join(ws *websocket.Conn) string {
	rl := &Role{Conn: ws}
	if r.B == nil {
		rl.Color = Black
		r.B = rl
		fmt.Println("黑色")
		return "black"
	}
	if r.W == nil {
		rl.Color = White
		r.W = rl
		fmt.Println("白色")
		return "white"
	}
	fmt.Println("吃瓜群众")
	rl.Color = 0
	r.OB = append(r.OB, rl)
	return ""
}

func (r *Room) ClearBoard() {
	r.board = &Board{}
}

// 消息广播
func (r *Room) Broad(data []byte, over bool, sender string) {
	if over {
		r.BroadAll(data, over, sender)
		return
	}

	if sender == "white" {
		r.B.Conn.WriteMessage(websocket.TextMessage, data)
	} else if sender == "black" {
		r.W.Conn.WriteMessage(websocket.TextMessage, data)
	}
	for _, v := range r.OB {
		v.Conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (r *Room) BroadAll(data []byte, over bool, sender string) {
	r.B.Conn.WriteMessage(websocket.TextMessage, data)
	r.W.Conn.WriteMessage(websocket.TextMessage, data)
	for _, v := range r.OB {
		v.Conn.WriteMessage(websocket.TextMessage, data)
	}
}

type Role struct {
	Conn  *websocket.Conn
	Color int
}
