package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/corbamico/tbcload"
)

//const testFile = "c:/project/go/src/github.com/ActiveState/teapot/lib/tbcload/tests/tbc10/proc.tbc"
//const testFile = "c:/project/go/src/github.com/corbamico/tbcload/cmd/test/test.tcl"

func usage() {
	message := `
Usage: distclbytecode [-h|--help] [file|url]
disassemble tcl bytecode file (usally .tbc file).

Example:
    distclbytecode  test.tbc  #disassemble a file named test.tbc
    distclbytecode  https://github.com/ActiveState/teapot/raw/master/lib/tbcload/tests/tbc10/proc.tbc
                	#disassemble from a url
	`
	fmt.Println(message)
	os.Exit(1)
}

func parseFile(uri string) {
	r, err := os.Open(uri)
	if err != nil {
		fmt.Printf("failed read from file (%s), error as (%s)\n", uri, err)
		return
	}
	p := tbcload.NewParser(r, os.Stdout)
	if err = p.Parse(); err != nil {
		fmt.Printf("failed parse file (%s), error as (%s)\n", uri, err)
		return
	}
}
func parseURL(uri string) {
	r, err := http.Get(uri)
	if err != nil {
		fmt.Printf("failed read from uri (%s), error as (%s)\n", uri, err)
		return
	}
	p := tbcload.NewParser(r.Body, os.Stdout)
	if err = p.Parse(); err != nil {
		fmt.Printf("failed parse uri (%s), error as (%s)\n", uri, err)
		return
	}
	r.Body.Close()
}

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
	}
	uri := os.Args[1]

	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
		parseURL(uri)
	} else {
		parseFile(uri)
	}
}
