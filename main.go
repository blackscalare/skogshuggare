package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	// Set default style and clear terminal.
	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()

	// Initialize game state.
	w, h := screen.Size()
	game := Game{
		player: Player{5, 5},
		border: Border{0, w - 1, 0, h - 1, 1},
		trees: map[int]*Tree{
			0: {6, 5, TreeStateAdult},
			1: {3, 4, TreeStateSapling},
			2: {10, 10, TreeStateAdult},
			3: {15, 5, TreeStateAdult},
			4: {16, 8, TreeStateAdult},
			5: {20, 15, TreeStateAdult},
			6: {25, 10, TreeStateAdult},
			7: {27, 12, TreeStateSapling},
			8: {17, 17, TreeStateSapling},
		},
		exit: false,
	}

	// Wait for Loop() goroutine to finish before moving on.
	var wg sync.WaitGroup
	wg.Add(1)
	go Ticker(&wg, screen, game)
	wg.Wait()
	screen.Fini()
}

func Ticker(wg *sync.WaitGroup, screen tcell.Screen, game Game) {
	// Wait for this goroutine to finish before resuming main().
	defer wg.Done()

	// Randomly seed map with trees in various states.
	game.PopulateTrees(screen)

	// Perform first draw.
	game.Draw(screen)

	// Initialize game update ticker.
	ticker := time.NewTicker(30 * time.Millisecond)

	// Update game state and re-draw on every tick.
	for range ticker.C {
		game.Update(screen)
		game.Draw(screen)
		if game.exit {
			return
		}
	}
}

func (game *Game) Update(screen tcell.Screen) {
	ev := screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			game.exit = true
			return
		case tcell.KeyUp:
			game.MovePlayer(screen, 1, DirUp)
		case tcell.KeyRight:
			game.MovePlayer(screen, 1, DirRight)
		case tcell.KeyDown:
			game.MovePlayer(screen, 1, DirDown)
		case tcell.KeyLeft:
			game.MovePlayer(screen, 1, DirLeft)
		case tcell.KeyRune:
			switch ev.Rune() {
			case rune(' '):
				game.Chop(screen, DirOmni)
			case rune('w'):
				game.Chop(screen, DirUp)
			case rune('d'):
				game.Chop(screen, DirRight)
			case rune('s'):
				game.Chop(screen, DirDown)
			case rune('a'):
				game.Chop(screen, DirLeft)
			}
		}
	case *tcell.EventResize:
		screen.Sync()
	}
}
