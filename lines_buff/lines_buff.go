package linesbuff

type Line struct {
	Content []byte
}

type LinesBuffer interface {
    Buffer() []Line
}
