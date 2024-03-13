package main

import (
	"io"
	"testing"
)

func TestTapeWrite(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, 0)

	newFileContents, _ := io.ReadAll(file)

	expected := "abc"
	got := string(newFileContents)

	if expected != got {
		t.Errorf("expected %q, got %q", expected, got)
	}

}
