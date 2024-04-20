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

type offset struct {
	x int
	y int
}

type Editor struct {
	Lines  linesbuff.LinesBuffer
	Screen tcell.Screen
	Events chan tcell.Event

	cursorPos  position
	offset     offset
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
        offset: offset{
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

	for y := e.offset.y; y < e.size.height; y++ {
		if y >= len(e.Lines.Buffer()) {
			break
		}
		for x := e.offset.x; x < e.size.width; x++ {
            if x >= len(e.Lines.Buffer()[y].Content) {
                break
            }
			e.Screen.SetContent(x, y, rune(e.Lines.Get(x, y)), nil, tcell.StyleDefault)
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

	if !e.insertMode && event.Rune() == 'i' {
		e.insertMode = true
		return
	}

	if e.insertMode {
		e.insertChar(event.Rune())
        return
	}

	// cursor movement
	switch event.Rune() {
	case 'j':
		e.cursorPos.y += 1
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'k':
		e.cursorPos.y -= 1
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'h':
		e.cursorPos.x -= 1
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'l':
		e.cursorPos.x += 1
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	}

}

func (e *Editor) insertChar(c rune) {
	e.Lines.Insert(c, e.cursorPos.x, e.cursorPos.y)
	e.cursorPos.x += 1
	e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
}

func (e *Editor) quit() {
	e.Screen.Fini()
	os.Exit(0)
}
