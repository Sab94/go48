package ZOAB

import (
	"math/rand"
	"reflect"
	"time"
)

type Board struct {
	Items *[][]int
	Size int
}

func NewBoard(size int) *Board{
	var board = make([][]int, size)
	for i := range board {
		board[i] = make([]int, size)
	}
	return &Board{
		Items: &board,
		Size: size,
	}
}

func (b *Board) PutNextNumber() {
	emptyCells := findEmptyCells(b.Items)
	if len(emptyCells) <= 0 {
		return
	}
	rndSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSrc)
	emptyCell := emptyCells[rnd.Intn(len(emptyCells))]
	board := *b.Items
	board[emptyCell/len(board)][emptyCell%len(board)] = 2
}

func findEmptyCells(board *[][]int) []int {
	var emptyCells = []int{}
	for i, row := range *board {
		for j, cell := range row {
			if cell == 0 {
				emptyCells = append(emptyCells, i*len(*board)+j)
			}
		}
	}
	return emptyCells
}

func (b *Board) SlideLeft() {
	for _, row := range *b.Items {
		stopMerge := 0
		for j := 1; j < len(row); j++ {
			if row[j] != 0 {
				for k := j; k > stopMerge; k-- {
					if row[k-1] == 0 {
						row[k-1] = row[k]
						row[k] = 0
					} else if row[k-1] == row[k] {
						row[k-1] += row[k]
						row[k] = 0
						stopMerge = k
						break
					} else {
						break
					}
				}
			}
		}
	}
}

func (b *Board) RotateBoard(counterClockWise bool) {
	board := *b.Items
	var rotatedBoard = make([][]int, len(board))
	for i, row := range board {
		rotatedBoard[i] = make([]int, len(row))
		for j := range row {
			if counterClockWise {
				rotatedBoard[i][j] = board[j][len(board)-i-1]
			} else {
				rotatedBoard[i][j] = board[len(board)-j-1][i]
			}
		}
	}
	b.Items = &rotatedBoard
}

func (b *Board) IsSame(previousBoard [][]int) bool {
	board := *b.Items
	for i, _ := range board {
		if !reflect.DeepEqual(board[i], previousBoard[i]) {
			return false
		}
	}
	return true
}