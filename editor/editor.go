package editor

import (
	"os"

	"github.com/gdamore/tcell"
	linesbuff "github.com/xsni1/gim/lines_buff"
)

type position struct {
	x int
	y int
}

type size struct {
	width  int
	height int
}

type Editor struct {
	Lines  linesbuff.LinesBuffer
	Screen tcell.Screen
	Events chan tcell.Event

	cursorPos  position
	insertMode bool
	size       size
}

func NewEditor(s tcell.Screen, fileContent []byte) *Editor {
	e := &Editor{
		Screen: s,
		Events: make(chan tcell.Event),
		cursorPos: position{
			x: 0,
			y: 0,
		},
		Lines: linesbuff.NewArrayBuffer(fileContent),
	}

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
				width, height := event.Size()
				e.size.width = width
				e.size.height = height
			}
		}

		e.Display()
		e.Screen.Sync()
	}
}

func (e *Editor) Display() {
	pos := position{
		x: 0,
		y: 0,
	}
	for _, l := range e.Lines.Buffer() {
		for _, c := range l.Content {
            e.Screen.SetContent(pos.x, pos.y, rune(c), nil, tcell.StyleDefault)
            pos.x++
		}
        pos.x = 0
        pos.y++
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
			e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
		}
	}

}

func (e *Editor) quit() {
	e.Screen.Fini()
	os.Exit(0)
}
