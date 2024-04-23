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

		// e.Screen.Clear()
		e.Display()
		e.Screen.Sync()
	}
}

func (e *Editor) Display() {
	pos := position{
		x: 0,
		y: 0,
	}
	for y := e.offset.y; y < e.size.height+e.offset.y; y++ {
		if y >= e.Lines.LinesNum() {
			break
		}
		for x := e.offset.x; x < e.size.width+e.offset.x; x++ {
			if x >= len(e.Lines.GetRow(y)) {
				e.Screen.SetContent(pos.x, pos.y, rune(' '), nil, tcell.StyleDefault)
				pos.x++
				continue
				// break
			}
			e.Screen.SetContent(pos.x, pos.y, rune(e.Lines.GetChar(x, y)), nil, tcell.StyleDefault)
			pos.x++
		}
		pos.x = 0
		pos.y++
	}
}

// TODO: add key to center the view
func (e *Editor) handleKeyEvent(event *tcell.EventKey) {
	if event.Key() == tcell.KeyESC {
		e.quit()
	}

	if !e.insertMode && event.Rune() == 'i' {
		e.insertMode = true
		return
	}

	if e.insertMode {
		// if event.Key() == tcell.KeyEnter {
		// 	e.Lines.Insert(c, e.cursorPos.x, e.cursorPos.y)
		// 	e.cursorPos.x += 1
		// 	e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
		// }
		e.insertChar(event.Rune())
		return
	}

	// TODO: CHECK IF EOL CHAR IS VISIBLE IN EDITOR
	// cursor movement
	switch event.Rune() {
	case 'j':
		if e.Lines.LinesNum()-1 <= e.cursorPos.y {
			return
		}
		e.cursorPos.y += 1
		if e.cursorPos.y > e.size.height-1 {
			e.offset.y++
			e.cursorPos.y -= 1
		}
		if len(e.Lines.GetRow(e.cursorPos.y))-1 <= e.cursorPos.x+e.offset.x {
			if e.offset.x+e.size.width > len(e.Lines.GetRow(e.cursorPos.y))-1 {
				if e.size.width > len(e.Lines.GetRow(e.cursorPos.y))-1 {
					e.offset.x = 0
				} else {
					// move view to the center:
					// TODO: move this to separate method
					// i want to center it only if the line we are going to is not visible on the screen
					if e.offset.x >= len(e.Lines.GetRow(e.cursorPos.y))-1 {
						e.offset.x = len(e.Lines.GetRow(e.cursorPos.y)) - 1 - (e.size.width / 2)
					}
                    
				}
			}

			e.cursorPos.x = len(e.Lines.GetRow(e.cursorPos.y)) - 1 - e.offset.x
		}
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'k':
		if e.cursorPos.y <= 0 {
			return
		}
		e.cursorPos.y -= 1
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'h':
		if e.cursorPos.x+e.offset.x <= 0 {
			return
		}
		e.cursorPos.x -= 1
		if e.cursorPos.x == 0 && e.offset.x > 0 {
			e.offset.x--
			e.cursorPos.x += 1
		}
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'l':
		if len(e.Lines.GetRow(e.cursorPos.y))-1 <= e.cursorPos.x+e.offset.x {
			return
		}
		e.cursorPos.x += 1
		if e.cursorPos.x >= e.size.width {
			e.offset.x++
			e.cursorPos.x -= 1
		}
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
