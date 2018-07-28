package render

import (
	"fmt"

	"github.com/nsf/termbox-go"
	"github.com/siggy/bbox/beatboxer/color"
	"github.com/siggy/bbox/beatboxer/keyboard"
)

type Terminal struct {
	kb *keyboard.Keyboard
}

func InitTerminal(kb *keyboard.Keyboard) *Terminal {
	return &Terminal{kb: kb}
}

// TODO: dry up
func (t *Terminal) TBprint(x, y int, msg string) {
	for _, c := range msg {
		t.kb.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

func (t *Terminal) Render(state State) {

	t.TBprint(0, 10, fmt.Sprintf("______render.State:_%+v______", state))

	for row := 0; row < ROWS; row++ {
		for col := 0; col < COLUMNS; col++ {
			// clear everything
			for i := 0; i < COLS_PER_BEAT; i++ {
				t.kb.SetCell(col*COLS_PER_BEAT+i, row, ' ', termbox.ColorBlack, termbox.ColorBlack)
			}

			fgColor := termbox.ColorBlack
			rune := ' '
			if state.LEDs[row][col] == color.Make(0, 0, 0, 127) {
				fgColor = termbox.ColorWhite
				rune = 'X'
			} else if state.LEDs[row][col] == color.Make(127, 0, 0, 0) {
				fgColor = termbox.ColorRed
				rune = 'O'
			} else if state.LEDs[row][col] == color.Make(250, 143, 94, 0) {
				// ceottk
				fgColor = termbox.ColorYellow
				rune = 'X'
			}

			if state.Transitions[row][col].Color == color.Make(0, 0, 0, 127) {
				location := int(COLS_PER_BEAT * state.Transitions[row][col].Location)
				t.kb.SetCell(col*COLS_PER_BEAT+location, row, 'X', termbox.ColorBlack, termbox.ColorWhite)
			}

			t.kb.SetCell(col*COLS_PER_BEAT, row, rune, termbox.ColorBlack, fgColor)
		}
	}

	t.kb.Flush()
	// TODO: map to actual LEDs
}