package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell"
	"github.com/xsni1/gim/editor"
)

// TODO: Terminal physical lines != text lines - this is causing LOTS of bugs
// TODO: When text goes to the next line old text gets overwritten
// TODO: \t not interpreted correctly

// 1. Store lines of text in an 2D array of bytes
//    a) place manipulation of this buffer behind some structure/interface - so in the future it is easy to replace array implementation with rope for example
// 2. Run main event loop, which first (or not first) refreshes display and then checks for any new events
// 3. tcell will be probably used as an library to handle terminal environment (ui, resizing, ansi sequences etc.) - it may be wise to abstract it behind some
//    interface aswell, so I can easily later replace it with mine implementation of tui library
// 4. I want to scroll both verticaly and horizontaly - need to store some offsets, so it is known exactly which lines (and its exact part) of text should be
//    displayed on the screen
//    a) how can I rerender only part of the screen so not whole thing is refreshed when user changes single letter? - for now it will be handled by tcell
//
//   tcell - use SetContent() to set each separate character and then at the end of the iteration of the main loop call Show() / Sync()

// tcell uses tput somehow to manipulate the terminal? or not it kind of implements its own tput by loading all the entries from terminfo?
// Fini() method does rmcup/smcup - alternate view or something like that

func main() {
	var fileContent []byte
	if len(os.Args) >= 2 {
		fileName := os.Args[1]
		if fileName != "" {
			bytes, err := os.ReadFile(fileName)
			if err != nil {
				fmt.Println("err reading file", err)
				os.Exit(1)
			}
			fileContent = bytes
		}
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	err = s.Init()
	if err != nil {
		log.Fatal(err)
	}
    // make sure the terminal gets recovered to its proper state in case of exiting program not via os.Exit() (panic is such a case)
    defer s.Fini()
    sigs := make(chan os.Signal)
    signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
    go func() {
        switch sig := <-sigs; sig {
        case syscall.SIGINT:
            s.Fini()
            os.Exit(0)
        case syscall.SIGTERM:
            s.Fini()
            os.Exit(0)
        }
    }()

	editor := editor.NewEditor(s, fileContent)
	go editor.EditorLoop()
	// dowiedziec sie jak w micro dzialaja key bindy - jak wylaczany jest program.
	// czy kazdy pane/term/buff w/e to tam osobna goroutina?
	for {
		e := s.PollEvent()
		if e != nil {
			editor.Events <- e
		}
	}
}
