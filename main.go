package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sab94/go48/ZOAB"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)
var quit = make(chan struct{})

func main() {
	size := flag.Int("size",4,"Set board size")
	flag.Parse()

	board := ZOAB.NewBoard(*size)

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	table := tview.NewTable().
		SetBorders(true)
	table.SetRect(0, 0, 600, 600)

	tempBoard := ZOAB.NewBoard(*size)
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyLeft:
					clone(*tempBoard.Items, *board.Items)
					board.SlideLeft()
					if !board.IsSame(*tempBoard.Items) {
						board.PutNextNumber()
					}
					drawTable(s, table, board)
				case tcell.KeyUp:
					clone(*tempBoard.Items, *board.Items)
					board.RotateBoard(true)
					board.SlideLeft()
					board.RotateBoard(false)
					if !board.IsSame(*tempBoard.Items) {
						board.PutNextNumber()
					}
					drawTable(s, table, board)
				case tcell.KeyRight:
					clone(*tempBoard.Items, *board.Items)
					board.RotateBoard(true)
					board.RotateBoard(true)
					board.SlideLeft()
					board.RotateBoard(false)
					board.RotateBoard(false)
					if !board.IsSame(*tempBoard.Items) {
						board.PutNextNumber()
					}
					drawTable(s, table, board)
				case tcell.KeyDown:
					clone(*tempBoard.Items, *board.Items)
					board.RotateBoard(false)
					board.SlideLeft()
					board.RotateBoard(true)
					if !board.IsSame(*tempBoard.Items) {
						board.PutNextNumber()
					}
					drawTable(s, table, board)
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	board.PutNextNumber()
	board.PutNextNumber()
	drawTable(s, table, board)

	select {
	case <-quit:
		break
	}
	s.Fini()
}

func drawTable(s tcell.Screen, table *tview.Table, board *ZOAB.Board) {
	for rIdx, r := range *board.Items {
		for cIdx, c := range r {
			color := tcell.ColorWhite
			value := "   "
			if c != 0 {
				if c < 10 {
					value = fmt.Sprintf(" %d ",c)
				} else {
					value = fmt.Sprintf("%d",c)
				}
			}
			cell := tview.NewTableCell(value).
				SetTextColor(color).
				SetAlign(tview.AlignCenter)
			table.SetCell(rIdx, cIdx,cell)
		}
	}
	table.Draw(s)
	s.Show()
}

func clone(new, source [][]int) {
	for i, r := range source {
		copy(new[i], r)
	}
}
