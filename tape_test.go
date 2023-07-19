package poker_test

import (
	"io"
	"testing"

	poker "github.com/fjahn78/poker"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &poker.Tape{File: file}

	_, _ = tape.Write([]byte("ABC"))

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "ABC"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
