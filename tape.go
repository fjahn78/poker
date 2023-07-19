package poker

import (
	"log"
	"os"
)

type Tape struct {
	File *os.File
}

func (t *Tape) Write(p []byte) (n int, err error) {
	if err = t.File.Truncate(0); err != nil {
		log.Fatal(err)
	}
	if _, err = t.File.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	return t.File.Write(p)
}
