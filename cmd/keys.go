package main

import (
	"github.com/nsf/termbox-go"
	"github.com/siggy/bbox/bbox"
)

func main() {
	defer termbox.Close()

	// beat changes
	//   keyboard => [main]
	msgs := []chan bbox.Beats{
		make(chan bbox.Beats),
	}

	// tempo changes
	//	 keyboard => loop
	tempo := make(chan int)

	// keyboard broadcasts quit with close(msgs)
	keyboard := bbox.InitKeyboard(bbox.WriteonlyBeats(msgs), tempo, bbox.KeyMapsPC, true)

	go keyboard.Run()
	defer keyboard.Close()

	for {
		select {
		case _, more := <-msgs[0]:
			if !more {
				return
			}
		}
	}
}
