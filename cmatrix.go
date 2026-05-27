// SPDX-FileCopyrightText: 2026 Reuben Percival
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
    "math/rand"
    "time"

    "github.com/nsf/termbox-go"
)

// column represents a single falling column of characters.
type column struct {
    y      int // current head position
    length int // total length of the column
    speed  int // frames between moves
    tick   int // counter for speed control
}

func initTermbox() error {
    err := termbox.Init()
    if err != nil {
        return err
    }
    termbox.HideCursor()
    return nil
}

func closeTermbox() {
    termbox.Close()
}

func randomChar() rune {
	if rand.Intn(2) == 0 {
		return rune(33 + rand.Intn(94)) // ASCII printable
	}
	return rune(0x30A0 + rand.Intn(96)) // Katakana
}

func main() {
	if err := initTermbox(); err != nil { panic(err) }
    defer closeTermbox()

    // initial size and columns
    w, h := termbox.Size()
    cols := make([]column, w)
    for i := range cols { cols[i] = column{y: rand.Intn(h), length: 5 + rand.Intn(15), speed: 1 + rand.Intn(3)} }

    ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS
    defer ticker.Stop()

    // channel for key press detection
    keyCh := make(chan struct{}, 1)
    go func() {
        for {
            ev := termbox.PollEvent()
            if ev.Type == termbox.EventKey {
                select { case keyCh <- struct{}{}: default: }
                return
            }
        }
    }()

    for {
        select {
        case <-ticker.C:
            termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
            w, h = termbox.Size()
            if len(cols) != w {
                newCols := make([]column, w)
                copy(newCols, cols)
                for i := len(cols); i < w; i++ { newCols[i] = column{y: rand.Intn(h), length: 5 + rand.Intn(15), speed: 1 + rand.Intn(3)} }
                cols = newCols
            }
            for x := 0; x < w; x++ {
                c := &cols[x]
                c.tick++
                if c.tick >= c.speed {
                    c.tick = 0
                    c.y = (c.y + 1) % h
                }
                for i := 0; i < c.length; i++ {
                    y := (c.y - i + h) % h
                    ch := randomChar()
                    fg := termbox.ColorGreen
                    if i == 0 { fg = termbox.ColorWhite } else if i < 3 { fg = termbox.ColorLightGreen }
                    termbox.SetCell(x, y, ch, fg, termbox.ColorDefault)
                }
            }
            termbox.Flush()
        case <-keyCh:
            return
        }
    }
}
