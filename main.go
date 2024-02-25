package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Position struct {
	x       int
	y       int
	yScroll int
}

func moveCursorTo(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func eraseLine() {
	fmt.Print("\033[2K")
}

func moveCurorDownBy(n int) {
	fmt.Print("\033[1B")
}

func moveCurorUpBy(n int) {
	fmt.Print("\033[1A")
}

func moveCurorLeftBy(n int) {
	fmt.Print("\033[1D")
}

func moveCurorRightBy(n int) {
	fmt.Print("\033[1C")
}

func redraw(lines [][]byte, pos Position) {
	_, height, _ := term.GetSize(0)

	moveCursorTo(1, 1)

	for i := 0; i < height; i++ {
		// terminal can be taller than the amount of lines
		if pos.yScroll+i >= len(lines) {
			break
		}
		fmt.Print(string(lines[pos.yScroll+i]))

		// go down
		fmt.Print("\033[1B")
		// move to the first column
		fmt.Print("\033[1G")
	}

	moveCursorTo(pos.x, pos.y)
}

// TODO: Terminal physical lines != text lines - this is causing LOTS of bugs
// TODO: When text goes to the next line old text gets overwritten
// TODO: \t not interpreted correctly
func main() {
	lines := [][]byte{}
	insertMode := false
	pos := Position{
		x: 1,
		y: 1,
	}

	// we have to restore it, otherwise terminal stays in raw mode
	prevState, _ := term.MakeRaw(0)

	f, _ := os.Open("file")
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Bytes())
	}

	// alternate xterm screen
	fmt.Print("\u001B[?1049h")
	defer func() {
		term.Restore(0, prevState)
		fmt.Print("\u001B[?1049l")
	}()

	redraw(lines, pos)

	for {
		// TODO: make it right
		in := make([]byte, 10)
		os.Stdin.Read(in)

		if in[0] == 27 {
			insertMode = false
		}

		if insertMode {
			fmt.Print(string(in[0]))
			lines[pos.y-1] = append(lines[pos.y-1][:pos.x], lines[pos.y-1][pos.x-1:]...)
			lines[pos.y-1][pos.x-1] = in[0]
			pos.x++

			eraseLine()
			// move to the first column
			fmt.Print("\033[1G")

			fmt.Print(string(lines[pos.y-1]))

			moveCursorTo(pos.x, pos.y)

			continue
		}

		if in[0] == 'q' {
			break
		}

		if in[0] == 'j' {
			moveCurorDownBy(1)
			pos.y++
		}

		if in[0] == 'k' {
            moveCurorUpBy(1)
			pos.y--
		}

		if in[0] == 'l' {
            moveCurorRightBy(1)
			pos.x++
		}

		if in[0] == 'h' {
            moveCurorLeftBy(1)
			pos.x--
		}

		if in[0] == 'i' {
			insertMode = true
		}

		redraw(lines, pos)
	}
}
