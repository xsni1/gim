package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell"
	"golang.org/x/term"
)

type Position struct {
	// terminal x, y coords
	tX int
	tY int
	// file x, y coords
	fX      int
	fY      int
	yScroll int
}

func moveCursorTo(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func eraseLine() {
	fmt.Print("\033[2K")
}

func moveCurorDownBy(n int) {
	fmt.Printf("\033[%dB", n)
}

func moveCurorUpBy(n int) {
	fmt.Printf("\033[%dA", n)
}

func moveCurorLeftBy(n int) {
	fmt.Printf("\033[%dD", n)
}

func moveCurorRightBy(n int) {
	fmt.Printf("\033[%dC", n)
}

func moveCursorToColumn(n int) {
	fmt.Printf("\033[%dG", n)
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

		moveCurorDownBy(1)
		moveCursorToColumn(1)
	}

	moveCursorTo(pos.tX, pos.tY)
}

// TODO: Terminal physical lines != text lines - this is causing LOTS of bugs
// TODO: When text goes to the next line old text gets overwritten
// TODO: \t not interpreted correctly

// 1. Store lines of text in an 2D array of bytes
//    a) place manipulation of this buffer behind some structure/interface - so in the future it is easy to replace array implementation with rope for example
// 2. Run main event loop which first (or not first) refreshes display and then checks for any new events
// 3. tcell will be probably used as an library to handle terminal environment (ui, resizing, ansi sequences etc.) - it may be wise to abstract it behind some
//    interface aswell, so I can easily later replace it with mine implementation of tui library
// 4. I want to scroll both verticaly and horizontaly - need to store some offsets, so it is known exactly which lines (and its exact part) of text should be
//    displayed on the screen
//    a) how can I rerender only part of the screen so not whole thing is refreshed when user changes single letter? - for now it will be handled by tcell
//
//   tcell - use SetContent() to set each separate character and then at the end of the iteration of the main loop call Show() / Sync()

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
    fmt.Println(s)

	// lines := [][]byte{}
	// insertMode := false
	// pos := Position{
	// 	tX: 1,
	// 	tY: 1,
	// 	fX: 1,
	// 	fY: 1,
	// }

	// // we have to restore it, otherwise terminal stays in raw mode
	// // we have to use raw mode, because in cooked (default) mode sends data to stdin when the user presses enter
	// // raw mode is basically a combination of flags set in termios
	// prevState, _ := term.MakeRaw(0)

	// f, _ := os.Open("file")
	// s := bufio.NewScanner(f)
	// for s.Scan() {
	// 	lines = append(lines, s.Bytes())
	// }

	// // alternate xterm screen
	// fmt.Print("\u001B[?1049h")
	// defer func() {
	// 	term.Restore(0, prevState)
	// 	fmt.Print("\u001B[?1049l")
	// }()

	// redraw(lines, pos)

	// for {
	// 	// TODO: make it right
	// 	in := make([]byte, 10)
	// 	os.Stdin.Read(in)

	// 	if in[0] == 27 {
	// 		insertMode = false
	// 	}

	// 	if insertMode {
	// 		fmt.Print(string(in[0]))
	// 		lines[pos.fY-1] = append(lines[pos.fY-1][:pos.fX], lines[pos.fY-1][pos.fX-1:]...)
	// 		lines[pos.fY-1][pos.fX-1] = in[0]
	// 		pos.tX++
	// 		// TODO: handle inserting last character in line
	// 		pos.fX++

	// 		eraseLine()
	// 		moveCursorToColumn(1)

	// 		fmt.Print(string(lines[pos.tY-1]))

	// 		moveCursorTo(pos.tX, pos.tY)

	// 		redraw(lines, pos)
	// 		continue
	// 	}

	// 	if in[0] == 'q' {
	// 		break
	// 	}

	// 	if in[0] == 'j' {
	// // TODO: Can it be called only once? or only on resize?
	// 		width, _, _ := term.GetSize(0)
	// 		// TODO: Differntiate console x,y and file x,y to handle multi-line single lines
	// 		// Here we could then:
	// 		if len(lines[]) > width {
	// 		  // file.y stays the same
	// 		  console.y++
	// 		}

	// 		moveCurorDownBy(1)
	// 		pos.tY++
	// 		pos.fY++
	// 	}

	// 	if in[0] == 'k' {
	// 		moveCurorUpBy(1)
	// 		pos.tY--
	// 		pos.fY--
	// 	}

	// 	if in[0] == 'l' {
	// 		moveCurorRightBy(1)
	// 		pos.tX++
	// 		pos.fX++
	// 	}

	// 	if in[0] == 'h' {
	// 		moveCurorLeftBy(1)
	// 		pos.tX--
	// 		pos.fX--
	// 	}

	// 	if in[0] == 'i' {
	// 		insertMode = true
	// 	}

	// 	redraw(lines, pos)
	// }
}
