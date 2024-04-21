package linesbuff

type Line struct {
	Content []byte
}

type LinesBuffer interface {
	Buffer() []Line
	Insert(r rune, x, y int)
	GetChar(x, y int) byte
	GetRow(y int) []byte
}
