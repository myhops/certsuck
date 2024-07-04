package keytool

import "io"

type yesReader struct {
	line      string
	written   []byte
}

var _ io.Reader = &yesReader{}

// Read returns y followed by a newline.
func (y *yesReader) Read(b []byte) (int, error) {
	if len(y.written) == 0 {
		y.written = append(y.written, append([]byte(y.line), []byte("\n")...)...)
	}
	n := copy(b, y.written)
	y.written = y.written[n:]
	return n, nil
}

func NewYesReader(line ...string) *yesReader {
	var l = "y"
	if len(line) > 0 &&  line[0] != "" {
		l = line[0]
	}
	return &yesReader{
		line: l,
	}
}