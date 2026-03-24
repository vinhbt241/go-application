package main

import (
	"io"
	"testing"
)

func TestTate_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, io.SeekStart)
	newFileContens, _ := io.ReadAll(file)

	got := string(newFileContens)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
