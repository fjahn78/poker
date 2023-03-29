package poker

import (
	"log"
	"os"
)

type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	if err = t.file.Truncate(0); err != nil {
		log.Fatal(err)
	}
	if _, err = t.file.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	return t.file.Write(p)
}
