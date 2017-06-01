package controllers

// 棋盘
type Board [15][17]int

func (b *Board) Fill(x, y int, color string) {
	m := map[string]int{"white": 1, "black": -1}
	b[x][y] = m[color]
}

func (bd *Board) checkWin(x, y, color int) bool {
	return bd.checkWinX(x, y, color) || bd.checkWinY(x, y, color) || bd.checkWin13(x, y, color) || bd.checkWin24(x, y, color)
}

func (bd *Board) checkWinX(x, y, color int) bool {
	count := 0
	for i := y - 1; i >= 0; i-- {
		if bd[x][i] == color {
			count++
		} else {
			break
		}
	}
	for i := y + 1; i <= 16; i++ {
		if bd[x][i] == color {
			count++
		} else {
			break
		}
	}
	return count >= 4
}

func (bd *Board) checkWinY(x, y, color int) bool {
	count := 0
	for i := x - 1; i >= 0; i-- {
		if bd[i][y] == color {
			count++
		} else {
			break
		}
	}
	for i := x + 1; i <= 14; i++ {
		if bd[i][y] == color {
			count++
		} else {
			break
		}
	}
	return count >= 4
}
func (bd *Board) checkWin13(x, y, color int) bool {
	count := 0
	for i, j := x-1, y+1; i >= 0 && j <= 16; {
		if bd[i][j] == color {
			count++
			i--
			j++
		} else {
			break
		}
	}
	for i, j := x+1, y-1; i <= 14 && j >= 0; i++ {
		if bd[i][j] == color {
			count++
			i++
			j--
		} else {
			break
		}
	}
	return count >= 4
}
func (bd *Board) checkWin24(x, y, color int) bool {
	count := 0
	for i, j := x-1, y-1; i >= 0 && j >= 0; {
		if bd[i][j] == color {
			count++
			i--
			j--
		} else {
			break
		}
	}
	for i, j := x+1, y+1; i <= 14 && j <= 16; {
		if bd[i][j] == color {
			count++
			i++
			j++
		} else {
			break
		}
	}
	return count >= 4
}
