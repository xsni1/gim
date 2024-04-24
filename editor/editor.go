package editor

import (
	"fmt"
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

	cursorPos     position
	offset        offset
	insertMode    bool
	size          size
	desiredCol    int
	gutterWidth   int
	infoBarHeight int
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
		Lines:       linesbuff.NewArrayBuffer(fileContent),
		gutterWidth: 3,
	}
	e.cursorPos.x = e.gutterWidth

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
		x: e.gutterWidth,
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
		pos.x = e.gutterWidth
		pos.y++
	}
	e.drawGutter()
	e.drawInfoBar()
}

func (e *Editor) drawGutter() {
	for y := 0; y < e.Lines.LinesNum(); y++ {
		for i, d := range fmt.Sprint(y + 1) {
			e.Screen.SetContent(i, y, rune(d), nil, tcell.StyleDefault)
		}
	}

	for y := e.Lines.LinesNum(); y < e.size.height; y++ {
		e.Screen.SetContent(0, y, rune('~'), nil, tcell.StyleDefault)
	}
}

func (e *Editor) drawInfoBar() {
	mode := "NORMAL"
	if e.insertMode {
		mode = "INSERT"
	}
	for x, l := range mode {
		e.Screen.SetContent(x, e.size.height-1, l, nil, tcell.StyleDefault)
		x++
	}

	// idk if this loop is even needed tbh
	for x := len(mode); x < e.size.width; x++ {
		e.Screen.SetContent(x, e.size.height-1, ' ', nil, tcell.StyleDefault)
	}

	pos := fmt.Sprintf("%d/%d", e.cursorPos.x, e.cursorPos.y)
	for x := e.size.width - len(pos); x < e.size.width; x++ {
		e.Screen.SetContent(x, e.size.height-1, rune(pos[x-(e.size.width-len(pos))]), nil, tcell.StyleDefault)
	}
}

func (e *Editor) clampPosX() {
	if len(e.Lines.GetRow(e.absPos().y))-1+e.gutterWidth <= e.absPos().x {
		if e.offset.x+e.size.width > len(e.Lines.GetRow(e.absPos().y))-1 {
			if e.size.width > len(e.Lines.GetRow(e.absPos().y))-1 {
				e.offset.x = 0
			} else {
				// move view to the center:
				// TODO: move this to separate method
				// i want to center it only if the line we are going to is not visible on the screen
				if e.offset.x >= len(e.Lines.GetRow(e.absPos().y))-1 {
					e.offset.x = len(e.Lines.GetRow(e.absPos().y)) - 1 - (e.size.width / 2)
				}
			}
		}
		e.cursorPos.x = len(e.Lines.GetRow(e.absPos().y)) - 1 - e.offset.x + e.gutterWidth
	}
}

func (e *Editor) absPos() position {
	return position{
		y: e.cursorPos.y + e.offset.y,
		x: e.cursorPos.x + e.offset.x,
	}
}

// TODO: add methods for getting current position relative and absolute
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
		e.insertChar(event.Rune())
		return
	}

	// cursor movement
	switch event.Rune() {
	case 'j':
		if e.Lines.LinesNum()-1 <= e.absPos().y {
			return
		}

		e.cursorPos.y++
		if e.cursorPos.y > e.size.height-1 {
			e.offset.y++
			e.cursorPos.y--
		}

		e.clampPosX()

		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'k':
		if e.absPos().y <= 0 {
			return
		}

		if e.cursorPos.y == 0 && e.offset.y > 0 {
			e.cursorPos.y++
			e.offset.y--
		}
		e.cursorPos.y--

		e.clampPosX()

		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'h':
		if e.absPos().x-e.gutterWidth <= 0 {
			return
		}
		e.cursorPos.x -= 1
		if e.cursorPos.x-e.gutterWidth+1 == 0 && e.offset.x > 0 {
			e.offset.x--
			e.cursorPos.x += 1
		}
		e.Screen.ShowCursor(e.cursorPos.x, e.cursorPos.y)
	case 'l':
		if len(e.Lines.GetRow(e.absPos().y))-1+e.gutterWidth <= e.absPos().x {
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
