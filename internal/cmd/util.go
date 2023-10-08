package cmd

import (
	"bytes"
	"io"
	"strings"
)

func readAll(r io.Reader) (string, error) {
	text, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	if text[0] == '#' {
		text = text[bytes.IndexByte(text, '\n'):]
	}

	return strings.TrimLeft(string(text), "\r\n"), nil
}
