# Developing Beatboxer Programs

This document describes development of "programs" for the Beatboxer.

## Quick Start

From the root of this repo:

```bash
go run cmd/beatboxer_noleds.go
```

The Beatboxer's 16x4 grid of buttons is emulated via your keyboard. The first
8x4 keys are bound to:
- `1` .. `8`
- `q` .. `i`
- `a` .. `k`
- `z` .. `,`

The shift key modifies the above to act on the second set of 8x4 keys.

Note all indexes in code are 0-based. Button 0x0 is `1`, 15x3 is `shift` + `,`.

Pressing button 15x1 (`shift` + `i`) 5 times will switch programs.

## Adding a Program

1. Add a new file under [`programs`](programs/)
1. Implement the `Program` interface, as defined in [`program.go`](program.go).
    This interface supports the following:
    - `Init`: initialize state and kick off any long-running processes
    - `Pressed`: receive button presses
    - `Render`: render LEDs
    - `Play`: play audio files
1. Add your new module to [`../cmd/beatboxer_noleds.go`](../cmd/beatboxer_noleds.go).

## TODO

- start/stop programs when not in use
- allow programs to yield to the next program
- better development renderer (perhaps web-based?)
- fix concurrency (`go run -race cmd/beatboxer_noleds.go` is not pretty)
- update LED rendering code for harness environment