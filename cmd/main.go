package main

import (
	"log"
	"os"

	"github.com/corbamico/tbcload"
)

const testFile = "D:\\Project\\go\\src\\github.com\\activestate\\teapot\\lib\\tbcload\\tests\\tbc10\\proc.tbc"

func main() {
	fs, err := os.Open(testFile)
	if err != nil {
		log.Fatalln(err)
		return
	}
	p := tbcload.NewParser(fs, os.Stderr)
	err = p.Parse()
	log.Fatalln(err)
}
