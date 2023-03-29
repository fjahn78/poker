package main

import (
	"io"
	"log"
)

type tape struct {
	file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
	if _, err = t.file.Seek(0, 0); err != nil {
    log.Fatal(err)
  }
	return t.file.Write(p)
}
