package main

import (
	"os"
	"os/signal"

	"github.com/siggy/bbox/bbox"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// beat changes
	//   keyboard => loop
	//   keyboard => render
	msgs := []chan bbox.Beats{
		make(chan bbox.Beats),
		make(chan bbox.Beats),
	}

	// tempo changes
	//	 keyboard => loop
	tempo := make(chan int)

	// ticks
	//   loop => render
	ticks := []chan int{
		make(chan int),
	}

	// interval changes
	//   loop => render
	intervals := []chan bbox.Interval{
		make(chan bbox.Interval),
	}

	// keyboard broadcasts quit with close(msgs)
	keyboard := bbox.InitKeyboard(bbox.WriteonlyBeats(msgs), tempo, bbox.KeyMapsPC, false)
	loop := bbox.InitLoop(msgs[0], tempo, bbox.WriteonlyInt(ticks), bbox.WriteonlyInterval(intervals))
	render := bbox.InitRender(msgs[1], ticks[0], intervals[0])

	go keyboard.Run()
	go loop.Run()
	go render.Run()

	defer keyboard.Close()
	defer loop.Close()
	defer render.Close()

	for {
		select {
		case <-sig:
			return
		}
	}
}
