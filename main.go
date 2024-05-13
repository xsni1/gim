package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell"
	"github.com/xsni1/gim/editor"
)

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

// Key bindings:
//   - define a map with default keybindings where key is an action name and value is pointer to a function
//      - "CursorUp": *editor.CursorUp (so this is something that lives in code and is compiled)
//   - another map that binds action name with a keypress is needed
//      - "CursorUp": "K" (in-code map that can be overwritten using file configuration)
//   - user can provide some external configuration to override default keybindings by providing a key with action name and value of keypress/es
//      - "CursorUp": "Z" (file configuration)
//   - we should be able to define more complex keybindings where user needs to provide a sequence of keypresses - use trie data structure
//   - also support for ctrl/alt mod keypresses (tcell probably supports detecting it)

func main() {
	var fileContent []byte
	var file *os.File
	if len(os.Args) >= 2 {
		fileName := os.Args[1]
		if fileName != "" {
			f, err := os.OpenFile(fileName, os.O_RDWR, os.ModeAppend)
			if err != nil {
				fmt.Println("err opening file", err)
				os.Exit(1)
			}
			defer f.Close()
			bytes, err := io.ReadAll(f)
			if err != nil {
				fmt.Println("err reading file", err)
				os.Exit(1)
			}
			fileContent = bytes
			file = f
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

	editor := editor.NewEditor(s, fileContent, file)
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
