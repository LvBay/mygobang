package controllers

import (
	"strconv"
	"strings"
)

type Req struct {
	Role     string // 执手
	Act      int    // 行为
	Position string // 位置
	Over     bool
	Sign     string
}

func (req *Req) OptInfo() (x, y int, color string) {
	arr := strings.Split(req.Position, "_")
	if len(arr) != 2 {
		return
	} // 请求信息非法
	x, _ = strconv.Atoi(arr[0])
	y, _ = strconv.Atoi(arr[1])
	return x, y, req.Role
}
