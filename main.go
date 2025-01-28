package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

var abc string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	if len(os.Args) != 1 {
		abc = os.Args[1]
	}
}

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
	for abc != "" {
		go func() {
			w, h := s.Size()
			elapsed := time.Since(start).Abs().Seconds()
			minutes := 0
			seconds := 0
			timeline := ""
			if elapsed >= 60 {
				minutes = int(elapsed) / 60
				seconds = int(elapsed) % 60
				timeline = fmt.Sprintf("%d min %d sec", minutes, seconds)
			} else {
				timeline = fmt.Sprintf("%.0f sec", elapsed)
			}
			length := len(timeline)
			for i := w - length; i < w; i++ {
				s.SetContent(i, h-1, rune(timeline[i-w+length]), nil, tcell.StyleDefault.Reverse(true))
			}
			time.Sleep(1 * time.Second)
			s.Sync()
		}()
		switch ev := s.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		case *tcell.EventMouse:
			w, h := s.Size()
			x, y := ev.Position()
			switch ev.Buttons() {
			case tcell.ButtonPrimary:
				if abc[0] == byte(result[x][y]) {
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
		if len(abc) > 0 {
			s.SetContent(0, 0, rune(abc[0]), nil, tcell.StyleDefault.Reverse(true))
		} else {
			break
		}
	}
	s.Fini()
	fmt.Printf("You win! %s\n", time.Since(start))
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
