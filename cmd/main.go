package main

import (
	"log"
	"os"

	"github.com/corbamico/tbcload"
)

const testFile = "c:/project/go/src/github.com/ActiveState/teapot/lib/tbcload/tests/tbc10/proc.tbc"

func main() {
	fs, err := os.Open(testFile)
	if err != nil {
		log.Fatalln(err)
		return
	}
	p := tbcload.NewParser(fs, os.Stderr)
	if err = p.Parse(); err != nil {
		log.Fatalln(err)
	}
}
