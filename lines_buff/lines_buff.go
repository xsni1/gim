package linesbuff

type Line struct {
	Content []byte
}

type LinesBuffer interface {
	Insert(r rune, x, y int)
	GetChar(x, y int) byte
	GetRow(y int) []byte
	LinesNum() int
	NewLine(x, y int)
	Buffer() []byte
	RemoveChar(x, y int)
}
