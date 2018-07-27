package main

import (
	"log"
	"os"

	"github.com/corbamico/tbcload"
)

//const testFile = "c:/project/go/src/github.com/ActiveState/teapot/lib/tbcload/tests/tbc10/proc.tbc"
const testFile = "c:/project/go/src/github.com/corbamico/tbcload/cmd/test/test.tcl"

func main() {
	fs, err := os.Open(testFile)
	if err != nil {
		log.Fatalln(err)
		return
	}
	p := tbcload.NewParser(fs, os.Stdout)
	if err = p.Parse(); err != nil {
		log.Fatalln(err)
	}
}
