package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	lines := [][]byte{}
	// we have to restore it, otherwise terminal stays in raw mode
	prevState, _ := term.MakeRaw(0)

	f, _ := os.Open("file")
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Bytes())
	}

	// alternate xterm screen
	fmt.Print("\u001B[?1049h")
	insertMode := false

	// move to the top
	fmt.Print("\033[1;1H")

	for _, line := range lines {
		fmt.Print(string(line))
		fmt.Print("\033[1B")
		// move to the first column
        fmt.Print("\033[1G")
		// fmt.Print("\033[1;1H")
	}

	defer func() {
		term.Restore(0, prevState)
		fmt.Print("\u001B[?1049l")
	}()

	for {
		in := make([]byte, 10)
		os.Stdin.Read(in)

		if in[0] == 27 {
			insertMode = false
		}

		if insertMode {
			fmt.Print(string(in[0]))
			continue
		}

		if in[0] == 'q' {
			break
		}

		if in[0] == 'j' {
			os.Stdout.Write([]byte("\033[1B"))
		}

		if in[0] == 'k' {
			os.Stdout.Write([]byte("\033[1A"))
		}

		if in[0] == 'l' {
			os.Stdout.Write([]byte("\033[1C"))
		}

		if in[0] == 'h' {
			os.Stdout.Write([]byte("\033[1D"))
		}

		if in[0] == 'i' {
			insertMode = true
		}

	}
}

// przeczytac line by line caly plik - rozdzielajac po \n albo \r\n
// printowac kazda linie jako osobna linie w terminalu
