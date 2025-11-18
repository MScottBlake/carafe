package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pborman/ansi"
)

type prependingWriter struct {
	prefix         string
	originalWriter io.Writer
}

// NewPrependingWriter returns a new writer that prepends a given prefix at the start of each line.
func NewPrependingWriter(originalWriter io.Writer, prefix string) io.Writer {
	return &prependingWriter{
		originalWriter: originalWriter,
		prefix:         prefix,
	}
}

func (w *prependingWriter) Write(p []byte) (n int, err error) {
	p = bytes.ReplaceAll(p, []byte("\r"), []byte("\n"))
	buf := bytes.NewBuffer(p)
	scanner := bufio.NewScanner(buf)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		line = stripLine(line)
		if line == "" {
			continue
		}

		formattedLine := fmt.Sprintf("%s%s\n", w.prefix, line)
		if bytesWritten, err := w.originalWriter.Write([]byte(formattedLine)); err != nil {
			return bytesWritten, err
		}
	}

	return len(p), nil
}

func stripLine(line string) string {
	strippedLine, err := ansi.Strip([]byte(line))
	if err != nil {
		// ansi.Strip returns an error if one of the escape codes is invalid. Return line so we don't swallow anything
		return line
	}

	line = string(strippedLine)
	line = strings.TrimSpace(line)

	return line
}
