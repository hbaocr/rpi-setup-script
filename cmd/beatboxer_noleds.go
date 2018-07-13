package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/siggy/bbox/bbox"
	"github.com/siggy/bbox/beatboxer"
	"github.com/siggy/bbox/beatboxer/programs/ceottk"
	"github.com/siggy/bbox/beatboxer/programs/drums"
	"github.com/siggy/bbox/beatboxer/render/web"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	harness := beatboxer.InitHarness(
		web.InitWeb(),
		bbox.KeyMapsPC,
	)

	harness.Register(&drums.DrumMachine{})
	harness.Register(&ceottk.Ceottk{})

	go harness.Run()

	for {
		select {
		case <-sig:
			fmt.Printf("Received OS signal, exiting\n")
			return
		}
	}
}
