package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

var (
	abc       string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	misclicks        = 0
	lent             = 52
)

func init() {
	if len(os.Args) != 1 {
		abc = os.Args[1]
		lent = len(abc)
	} else {
		fmt.Printf("Usage: %s [set of characters]\nPress any key to continue with default set, \"q\" to quit.\n", os.Args[0])
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Error setting terminal to raw mode:", err)
			os.Exit(1)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		buf := make([]byte, 1)
		_, err = os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("Error reading from standard input: %s", err)
			os.Exit(1)
		}
		if buf[0] == 'q' {
			os.Exit(1)
		}
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
	w, h := s.Size()
	if w < 80 || h < 20 {
		s.Fini()
		fmt.Println("minimum window size is 80x20")
		os.Exit(1)
	}
	result, marks := print(s)
	s.Sync()
	start := time.Now()
	go func() {
		for {
			timestring := fmt.Sprintf("%s", time.Since(start).Abs().Round(time.Second))
			length := len(timestring)
			for i := w - length; i < w; i++ {
				s.SetContent(i, h-1, rune(timestring[i-w+length]), nil, tcell.StyleDefault.Reverse(true))
			}
			time.Sleep(1 * time.Second)
			s.Sync()
		}
	}()
	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			//		s.Fini()
			//		fmt.Println("Resizing is unsupported")
			//		os.Exit(1)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Key() == tcell.KeyCtrlC {
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
								marks[i][j] = true
							}
						}
					}
					_, abc = abc[0], abc[1:]
				} else {
					blink(s, result[x][y], x, y, marks[x][y])
				}
			}
		}
		if len(abc) > 0 {
			nextchar := fmt.Sprintf("next char %c", abc[0])
			for i := 0; i < len(nextchar); i++ {
				s.SetContent(i, 0, rune(nextchar[i]), nil, tcell.StyleDefault.Reverse(true))
			}
		} else {
			break
		}
	}
	s.Fini()
	fmt.Printf("You win! %s. %dx%d matrix. %d characters long set. %d misclicks\n", time.Since(start), w, h, lent, misclicks)
}

func blink(s tcell.Screen, letter rune, x, y int, mark bool) {
	misclicks++
	for i := 0; i < 5; i++ {
		s.SetContent(x, y, letter, nil, tcell.StyleDefault.Background(tcell.ColorRed))
		s.Sync()
		time.Sleep(50 * time.Millisecond)
		if mark {
			s.SetContent(x, y, letter, nil, tcell.StyleDefault.Background(tcell.ColorGreen))
		} else {
			s.SetContent(x, y, letter, nil, tcell.StyleDefault)
		}
		s.Sync()
	}
}

func print(s tcell.Screen) ([][]rune, [][]bool) {
	w, h := s.Size()
	result := make([][]rune, w)
	marks := make([][]bool, w)
	for i := range result {
		result[i] = make([]rune, h)
	}
	for i := range marks {
		marks[i] = make([]bool, h)
	}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			randomNumber := rand.Int() % len(abc)
			letter := rune(abc[randomNumber])
			result[i][j] = letter
			s.SetContent(i, j, letter, nil, tcell.StyleDefault)
		}
	}
	return result, marks
}
