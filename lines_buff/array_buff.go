package linesbuff

import (
	"slices"
)

type ArrayBuffer struct {
	lines []Line
}

// func (ab *ArrayBuffer) Buffer() []Line {
// 	return ab.lines
// }

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

// func (ab *ArrayBuffer) GetCol(x int) []byte {
// 	return ab.lines[y].Content
// }
