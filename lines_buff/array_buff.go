package linesbuff

import (
	"slices"
)

type ArrayBuffer struct {
	lines []Line
}

// alternatively reader could be used here, probably memory overhead would be lower?
func NewArrayBuffer(fileContent []byte) *ArrayBuffer {
	lines := []Line{}
	for {
		// TODO: \r\n support
		n := slices.Index(fileContent, '\n')
		if n == -1 {
			// does linux add \n to eof?
			// this may be the case where there is need for empty line at the very end
			break
		}
		// this has to be copied - otherwise the subslice will share the same underlying array (as long as it fits in the capacit ,
		// because new underlying array will get allocated then)
		cp := make([]byte, len(fileContent[:n+1]))
		copy(cp, fileContent[:n+1])
		lines = append(lines, Line{
			// do i need \n character in line array?
			Content: cp,
		})
		fileContent = fileContent[n+1:]
	}

	return &ArrayBuffer{
		lines: lines,
	}
}

// look for ways to optimize - over-allocate capacity?
// also check how to profile such things
func (ab *ArrayBuffer) Insert(r rune, x, y int) {
	line := &ab.lines[y]
	line.Content = append(line.Content[:x], append([]byte{byte(r)}, line.Content[x:]...)...)
}

func (ab *ArrayBuffer) LinesNum() int {
	return len(ab.lines)
}

func (ab *ArrayBuffer) GetChar(x, y int) byte {
	return ab.lines[y].Content[x]
}

func (ab *ArrayBuffer) GetRow(y int) []byte {
	return ab.lines[y].Content
}

func (ab *ArrayBuffer) NewLine(x, y int) {
	newline := ab.lines[y].Content[x:]
	ab.lines = append(ab.lines[:y+1], append([]Line{{Content: newline}}, ab.lines[y+1:]...)...)
	curline := make([]byte, len(ab.lines[y].Content[:x]), len(ab.lines[y].Content[:x])+1)
	copy(curline, ab.lines[y].Content[:x])
	curline = append(curline, '\n')
	ab.lines[y].Content = curline
}

func (ab *ArrayBuffer) Buffer() []byte {
	// c := ab.LinesNum()
	// for _, l := range ab.lines {
	//     c *= len(l.Content)
	// }
	buf := make([]byte, 0, 0)
	for _, l := range ab.lines {
		buf = append(buf, l.Content...)
	}
	return buf
}

func (ab *ArrayBuffer) RemoveChar(x, y int) {
	// remove line
	if len(ab.lines[y].Content) == 1 {
		ab.lines = append(ab.lines[:y], ab.lines[y+1:]...)
		return
	}

	ab.lines[y].Content = append(ab.lines[y].Content[:x], ab.lines[y].Content[x+1:]...)
}
