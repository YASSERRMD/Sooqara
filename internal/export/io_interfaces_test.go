package export

import (
	"io"
	"testing"
)

func TestIOReaderInterface(t *testing.T) {
	var r io.Reader = nil
	_ = r
}

func TestIOWriterInterface(t *testing.T) {
	var w io.Writer = nil
	_ = w
}

func TestIOCloserInterface(t *testing.T) {
	var c io.Closer = nil
	_ = c
}
