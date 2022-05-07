package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

var abc string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	s.EnableMouse()
	result := print(s)
	s.Sync()
	start := time.Now()
	for {
		go func() {
			w, h := s.Size()
			elapsed := time.Since(start)
			length := len(elapsed.String())
			for i := w - length; i < w; i++ {
				s.SetContent(i, h-1, rune(elapsed.String()[i-w+length]), nil, tcell.StyleDefault.Reverse(true))
			}
			s.Sync()
		}()
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			result = print(s)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyRune {
				if ev.Rune() == 'r' {
					result = print(s)
				}
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			switch ev.Buttons() {
			case tcell.ButtonPrimary:
				if abc[0] == byte(result[x][y]) {
					w, h := s.Size()
					for i := 0; i < w; i++ {
						for j := 0; j < h; j++ {
							if abc[0] == byte(result[i][j]) {
								s.SetContent(i, j, result[i][j], nil, tcell.StyleDefault.Background(tcell.ColorGreen))
							}
						}
					}
					_, abc = abc[0], abc[1:]
				} else {
					blink(s, result[x][y], x, y)
				}
			}
		}
	}
}

func blink(s tcell.Screen, letter rune, x, y int) {
	for i := 0; i < 5; i++ {
		s.SetContent(x, y, letter, nil, tcell.StyleDefault.Background(tcell.ColorRed))
		s.Sync()
		time.Sleep(50 * time.Millisecond)
		s.SetContent(x, y, letter, nil, tcell.StyleDefault)
		s.Sync()
	}
}

func print(s tcell.Screen) [][]rune {
	w, h := s.Size()
	result := make([][]rune, w)
	for i := range result {
		result[i] = make([]rune, h)
	}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			randomNumber := rand.Int() % len(abc)
			letter := rune(abc[randomNumber])
			result[i][j] = letter
			s.SetContent(i, j, letter, nil, tcell.StyleDefault)
		}
	}
	return result
}
