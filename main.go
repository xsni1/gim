package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Position struct {
	x int
	y int
}

// TODO: Terminal physical lines != text lines - this is causing LOTS of bugs
// TODO: When text goes to the next line old text gets overwritten
// TODO: \t not interpreted correctly

func main() {
	lines := [][]byte{}
	insertMode := false
	cursorPosition := Position{
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

	// move to the top
	fmt.Print("\033[1;1H")

	for _, line := range lines {
		fmt.Print(string(line))

        // go down
		fmt.Print("\033[1B")
		// move to the first column
		fmt.Print("\033[1G")
	}

	fmt.Print("\033[1;1H")

	defer func() {
		term.Restore(0, prevState)
		fmt.Print("\u001B[?1049l")
	}()

	for {
		// TODO: make it right
		in := make([]byte, 10)
		os.Stdin.Read(in)

		if in[0] == 27 {
			insertMode = false
		}

		if insertMode {
			fmt.Print(string(in[0]))
			lines[cursorPosition.y-1] = append(lines[cursorPosition.y-1][:cursorPosition.x], lines[cursorPosition.y-1][cursorPosition.x-1:]...)
			lines[cursorPosition.y-1][cursorPosition.x-1] = in[0]
            cursorPosition.x++

			// erase entire line
			fmt.Print("\033[2K")
			// move to the first column
			fmt.Print("\033[1G")

			fmt.Print(string(lines[cursorPosition.y-1]))

			fmt.Printf("\033[%d;%dH", cursorPosition.y, cursorPosition.x)

			continue
		}

		if in[0] == 'q' {
			break
		}

		if in[0] == 'j' {
			os.Stdout.Write([]byte("\033[1B"))
			cursorPosition.y++
		}

		if in[0] == 'k' {
			os.Stdout.Write([]byte("\033[1A"))
			cursorPosition.y--
		}

		if in[0] == 'l' {
			os.Stdout.Write([]byte("\033[1C"))
			cursorPosition.x++
		}

		if in[0] == 'h' {
			os.Stdout.Write([]byte("\033[1D"))
			cursorPosition.x--
		}

		if in[0] == 'i' {
			insertMode = true
		}

	}
}

