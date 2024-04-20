package linesbuff

import (
	"slices"
)

type ArrayBuffer struct {
	lines []Line
}

func (ab *ArrayBuffer) Buffer() []Line {
    return ab.lines
}

// alternatively reader could be used here, probably memory overhead would be lower?
func NewArrayBuffer(fileContent []byte) *ArrayBuffer {
	lines := []Line{}
	for {
		// TODO: \r\n support
		n := slices.Index(fileContent, '\n')
		if n == -1 {
			// this may be the case where there is need for empty line at the very end
			break
		}
		lines = append(lines, Line{
			// do i need \n character in line array?
			Content: fileContent[:n],
		})
		fileContent = fileContent[n+1:]
	}

	return &ArrayBuffer{
		lines: lines,
	}
}
