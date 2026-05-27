// SPDX-FileCopyrightText: 2026 Reuben Percival
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type cell struct {
	val  int
	head bool
}

const (
	cEmpty = -1
	cSpace = ' '
)

var (
	grid    [][]cell
	lengths []int
	gaps    []int
	speeds  []int

	async     bool
	boldMode  int
	classic   bool
	rainbow   bool
	lambda    bool
	changes   bool
	oldstyle  bool
	saver     bool
	paused    bool
	locked    bool
	msg       string
	updateDel int
	mainColor termbox.Attribute
	randMin   int
	randNum   int
	frame     int
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: cmatrix [-abBcflLnorsxkmV] [-u delay] [-C color] [-M message]\n")
	flag.PrintDefaults()
}

func version() {
	fmt.Fprintf(os.Stderr, "GOmatrix 1.0\n")
}

func randChar() int {
	if randNum <= 0 {
		return ' '
	}
	return rand.Intn(randNum) + randMin
}

func randColor() termbox.Attribute {
	pick := rand.Intn(6)
	switch pick {
	case 0:
		return termbox.ColorGreen
	case 1:
		return termbox.ColorBlue
	case 2:
		return termbox.ColorBlack
	case 3:
		return termbox.ColorYellow
	case 4:
		return termbox.ColorCyan
	default:
		return termbox.ColorMagenta
	}
}

func initGrid(w, h int) {
	rows := h + 1
	grid = make([][]cell, rows)
	for i := range grid {
		grid[i] = make([]cell, w)
		for j := range grid[i] {
			grid[i][j].val = cEmpty
		}
	}
	lengths = make([]int, w)
	gaps = make([]int, w)
	speeds = make([]int, w)

	for j := 0; j < w; j += 2 {
		gaps[j] = rand.Intn(h) + 1
		lengths[j] = rand.Intn(h-3) + 3
		if lengths[j] < 3 {
			lengths[j] = 3
		}
		if rows > 1 {
			grid[1][j].val = cSpace
		}
		speeds[j] = rand.Intn(3) + 1
	}
}

func updateGrid(w, h int) {
	rows := h + 1
	highnum := randMin + randNum

	for j := 0; j < w; j += 2 {
		if (!async || frame > speeds[j]) && !paused {
			if oldstyle {
				for i := rows - 1; i >= 1; i-- {
					grid[i][j].val = grid[i-1][j].val
					grid[i][j].head = false
				}
				r := rand.Intn(randNum+8) + randMin

				if grid[1][j].val == 0 {
					grid[0][j].val = 1
				} else if grid[1][j].val == cSpace || grid[1][j].val == cEmpty {
					if gaps[j] > 0 {
						grid[0][j].val = cSpace
						gaps[j]--
					} else {
						if rand.Intn(3) == 1 {
							grid[0][j].val = 0
						} else {
							grid[0][j].val = rand.Intn(randNum) + randMin
						}
						gaps[j] = rand.Intn(h) + 1
					}
				} else if r > highnum && grid[1][j].val != 1 {
					grid[0][j].val = cSpace
				} else {
					grid[0][j].val = rand.Intn(randNum) + randMin
				}
			} else {
				if grid[0][j].val == cEmpty && grid[1][j].val == cSpace && gaps[j] > 0 {
					gaps[j]--
				} else if grid[0][j].val == cEmpty && grid[1][j].val == cSpace {
					lengths[j] = rand.Intn(h-3) + 3
					if lengths[j] < 3 {
						lengths[j] = 3
					}
					grid[0][j].val = randChar()
					gaps[j] = rand.Intn(h) + 1
				}

				i := 0
				var segLen int
				firstSeg := false
				for i <= h {
					for i <= h && (grid[i][j].val == cSpace || grid[i][j].val == cEmpty) {
						i++
					}
					if i > h {
						break
					}

					segStart := i
					segLen = 0
					for i <= h && grid[i][j].val != cSpace && grid[i][j].val != cEmpty {
						grid[i][j].head = false
						if changes && rand.Intn(8) == 0 {
							grid[i][j].val = randChar()
						}
						i++
						segLen++
					}

					if i > h {
						grid[segStart][j].val = cSpace
						continue
					}

					grid[i][j].val = randChar()
					grid[i][j].head = true

					if segLen > lengths[j] || firstSeg {
						grid[segStart][j].val = cSpace
						grid[0][j].val = cEmpty
					}
					firstSeg = true
					i++
				}
			}
		}
	}
}

func renderGrid(w, h int) {
	for j := 0; j < w; j += 2 {
		startRow := 1
		endRow := h
		if oldstyle {
			startRow = 0
			endRow = h - 1
		}
		for i := startRow; i <= endRow; i++ {
			c := grid[i][j]
			if c.val == cEmpty && !c.head {
				continue
			}
			if c.val == cSpace {
				continue
			}

			color := mainColor
			var attr termbox.Attribute

			if c.val == 0 || (c.head && !rainbow) {
				color = termbox.ColorWhite
				attr = termbox.AttrBold
			} else if rainbow {
				color = randColor()
			}

			if c.val == 1 {
				termbox.SetCell(j, i-startRow, '|', color|attr, termbox.ColorDefault)
				continue
			}

			ch := rune(c.val)
			if lambda && ch != cSpace {
				ch = 'λ'
			}

			if boldMode == 2 || (boldMode == 1 && c.val%2 == 0) {
				attr |= termbox.AttrBold
			}
			if boldMode == -1 {
				attr &^= termbox.AttrBold
			}

			termbox.SetCell(j, i-startRow, ch, color|attr, termbox.ColorDefault)
		}
	}
}

func renderMessage(w, h int) {
	msgRunes := []rune(msg)
	msgX := h / 2
	msgY := w/2 - len(msgRunes)/2

	for x := msgY - 2; x < msgY+len(msgRunes)+2 && x < w; x++ {
		termbox.SetCell(x, msgX-1, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	for x := 0; x < w; x++ {
		termbox.SetCell(x, msgX, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	for i, r := range msgRunes {
		termbox.SetCell(msgY+i, msgX, r, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
	}
	for x := msgY - 2; x < msgY+len(msgRunes)+2 && x < w; x++ {
		termbox.SetCell(x, msgX+1, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
}

func main() {
	flagAsync := flag.Bool("a", false, "Asynchronous scroll")
	flagBold := flag.Bool("b", false, "Bold characters on")
	flagBoldAll := flag.Bool("B", false, "All bold characters (overrides -b)")
	flagClassic := flag.Bool("c", false, "Use Japanese characters")
	flagNoBold := flag.Bool("n", false, "No bold characters")
	flagOld := flag.Bool("o", false, "Old-style scrolling")
	flagSaver := flag.Bool("s", false, "Screensaver mode (exit on any key)")
	flagRainbow := flag.Bool("r", false, "Rainbow mode")
	flagLambda := flag.Bool("m", false, "Lambda mode")
	flagChanges := flag.Bool("k", false, "Characters change while scrolling")
	flagNoOp := flag.Bool("f", false, "Force linux $TERM (no-op in Go version)")
	flagNoOp2 := flag.Bool("l", false, "Linux console mode (no-op in Go version)")
	flagNoOp3 := flag.Bool("x", false, "X window mode (no-op in Go version)")
	flagTty := flag.String("t", "", "Set tty (no-op in Go version)")
	flagLock := flag.Bool("L", false, "Lock mode")
	flagMsg := flag.String("M", "", "Display message")
	flagUpdate := flag.Int("u", 4, "Screen update delay (0-10)")
	flagColor := flag.String("C", "green", "Matrix color (green, red, blue, white, yellow, cyan, magenta, black)")
	flagHelp := flag.Bool("h", false, "Print usage")
	flagVersion := flag.Bool("V", false, "Print version")
	flag.Parse()

	if *flagHelp {
		usage()
		return
	}
	if *flagVersion {
		version()
		return
	}

	async = *flagAsync
	oldstyle = *flagOld
	saver = *flagSaver
	classic = *flagClassic
	rainbow = *flagRainbow
	lambda = *flagLambda
	changes = *flagChanges
	locked = *flagLock
	msg = *flagMsg
	updateDel = *flagUpdate

	_ = flagNoOp
	_ = flagNoOp2
	_ = flagNoOp3
	_ = flagTty

	if locked && msg == "" {
		msg = "Computer locked."
	}

	if *flagNoBold {
		boldMode = -1
	} else if *flagBoldAll {
		boldMode = 2
	} else if *flagBold {
		boldMode = 1
	}

	switch *flagColor {
	case "green":
		mainColor = termbox.ColorGreen
	case "red":
		mainColor = termbox.ColorRed
	case "blue":
		mainColor = termbox.ColorBlue
	case "white":
		mainColor = termbox.ColorWhite
	case "yellow":
		mainColor = termbox.ColorYellow
	case "cyan":
		mainColor = termbox.ColorCyan
	case "magenta":
		mainColor = termbox.ColorMagenta
	case "black":
		mainColor = termbox.ColorBlack
	default:
		fmt.Fprintf(os.Stderr, "Invalid color. Valid: green, red, blue, white, yellow, cyan, magenta, black.\n")
		return
	}

	if classic {
		randMin = 0xFF66
		randNum = 0xFF9D - 0xFF66
	} else {
		randMin = 33
		randNum = 123 - 33
	}

	if err := termbox.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "termbox init: %v\n", err)
		return
	}
	defer termbox.Close()
	termbox.HideCursor()

	// ignore unused flags: f, l, x, t
	_ = flag.Lookup("f")
	_ = flag.Lookup("l")
	_ = flag.Lookup("x")
	_ = flag.Lookup("t")

	w, h := termbox.Size()
	if h < 10 {
		h = 10
	}
	if w < 10 {
		w = 10
	}
	initGrid(w, h)

	d := time.Duration(updateDel*10) * time.Millisecond
	if d < time.Millisecond {
		d = time.Millisecond
	}
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	eventCh := make(chan termbox.Event, 8)
	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventCh:
			switch ev.Type {
			case termbox.EventKey:
				if saver {
					return
				}
				switch ev.Ch {
				case 'q':
					if !locked {
						return
					}
				case 'a':
					async = !async
				case 'b':
					boldMode = 1
				case 'B':
					boldMode = 2
				case 'n':
					boldMode = 0
				case 'p', 'P':
					paused = !paused
				case 'r':
					rainbow = !rainbow
				case 'm':
					lambda = !lambda
				case '!':
					mainColor = termbox.ColorRed
					rainbow = false
				case '@':
					mainColor = termbox.ColorGreen
					rainbow = false
				case '#':
					mainColor = termbox.ColorYellow
					rainbow = false
				case '$':
					mainColor = termbox.ColorBlue
					rainbow = false
				case '%':
					mainColor = termbox.ColorMagenta
					rainbow = false
				case '^':
					mainColor = termbox.ColorCyan
					rainbow = false
				case '&':
					mainColor = termbox.ColorWhite
					rainbow = false
				case 'L':
					locked = true
				}
				if ev.Ch >= '0' && ev.Ch <= '9' {
					updateDel = int(ev.Ch - '0')
					d := time.Duration(updateDel*10) * time.Millisecond
					if d < time.Millisecond {
						d = time.Millisecond
					}
					ticker.Reset(d)
				}
				if ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc {
					return
				}

			case termbox.EventResize:
				w, h = ev.Width, ev.Height
				if h < 10 {
					h = 10
				}
				if w < 10 {
					w = 10
				}
				initGrid(w, h)

			case termbox.EventError:
				return
			}

		case <-ticker.C:
			if paused {
				termbox.Flush()
				continue
			}

			w, h = termbox.Size()
			if h < 10 {
				h = 10
			}
			if w < 10 {
				w = 10
			}

			if len(grid) == 0 || len(grid[0]) != w || len(grid) != h+1 {
				initGrid(w, h)
			}

			frame++
			if frame > 4 {
				frame = 1
			}

			updateGrid(w, h)

			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			renderGrid(w, h)
			if msg != "" {
				renderMessage(w, h)
			}
			termbox.Flush()
		}
	}
}
