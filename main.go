package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	term.MakeRaw(0)
	buf := []byte("aaaaaaaqweeeefdsjknfjnfsjdnknsdkjfnjk\n")
	os.Stdout.Write(buf)
	// alternate xterm screen
	fmt.Print("\u001B[?1049h")
	insertMode := false
	defer func() {
		// term.Restore(0, prevState)
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
