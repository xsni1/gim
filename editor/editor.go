package editor

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/xsni1/gim/debug"
	linesbuff "github.com/xsni1/gim/lines_buff"
)

type position struct {
	x int
	y int
}

type Editor struct {
	Lines  linesbuff.LinesBuffer
	Screen tcell.Screen
	Events chan tcell.Event

	cursorPos  position
	insertMode bool
}

func NewEditor(s tcell.Screen) *Editor {
	e := &Editor{
		Screen: s,
		Events: make(chan tcell.Event),
		cursorPos: position{
			x: 0,
			y: 0,
		},
	}

	// a := fmt.Sprintf("NEW EDITOR METHOD - %p", &s)
	// debug.Debug(a)

	s.ShowCursor(e.cursorPos.x, e.cursorPos.y)

	return e
}

func (e *Editor) EditorLoop() {
	for {
		select {
		case ev := <-e.Events:
			switch event := ev.(type) {
			case *tcell.EventKey:
				e.handleKeyEvent(event)
			case *tcell.EventResize:
				// fmt.Println("resize event")
			}
		}

        e.Screen.Sync()
	}
}

func (e *Editor) handleKeyEvent(event *tcell.EventKey) {
	if event.Key() == tcell.KeyESC {
		e.quit()
	}

	switch event.Rune() {
	case 'i':
		e.insertMode = true
	case 'j':
		if e.insertMode {

		} else {
			e.cursorPos.y += 1
            debug.Debug(fmt.Sprint(e.cursorPos.y))
            debug.Debug("\n")
			e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
		}
	}

}

func (e *Editor) quit() {
	e.Screen.Fini()
	os.Exit(0)
}
