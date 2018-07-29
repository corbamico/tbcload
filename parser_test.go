package tbcload

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

const uriPath = "https://github.com/ActiveState/teapot/raw/master/lib/tbcload/tests/tbc10"

var fileNames = []string{
	"aux1.tbc",
	"break.tbc",
	"break1.tbc",
	"break2.tbc",
	"catch.tbc",
	"catch1.tbc",
	"cont.tbc",
	"cont1.tbc",
	"expr.tbc",
	"expr1.tbc",
	"expr2.tbc",
	"for.tbc",
	"foreach.tbc",
	"interp.tbc",
	"override.tbc",
	"proc.tbc",
	"procbod1.tbc",
	"procbod2.tbc",
	"procbod3.tbc",
	"procbreak1.tbc",
	"proccatch1.tbc",
	"proccatch2.tbc",
	"proccontinue1.tbc",
	"procepc1.tbc",
	"procepc2.tbc",
	"procshd1.tbc",
	"procshd2.tbc",
	"procshd3.tbc",
	"procshd4.tbc",
	"procshd5.tbc",
	"procshd6.tbc",
	"procshd7.tbc",
	"procshd8.tbc",
	"procvar1.tbc",
	"procvar2.tbc",
	"while.tbc",
}

func testURLTBC(t *testing.T, uriPath string, fileName string) {
	r, err := http.Get(fmt.Sprintf("%s/%s", uriPath, fileName))
	if err != nil {
		t.Errorf("failed read uri:%s", fileName)
		return
	}
	p := NewParser(r.Body, ioutil.Discard)
	err = p.Parse()
	if err != nil {
		t.Errorf("failed parse uri:%s;err=%s", fileName, err)
	}
	r.Body.Close()
	t.Logf("success uri:%s", fileName)
}

func testURLs(t *testing.T, uriPath string, fileNames []string) {
	for _, s := range fileNames {
		testURLTBC(t, uriPath, s)
	}
}
func TestParser(t *testing.T) {
	testURLs(t, uriPath, fileNames)
}

func TestSingleFile(t *testing.T) {
	//sFile := "D:\\Program Files (x86)\\TclPro1.4\\win32-ix86\\bin\\simple.tbc"
	sFile := "D:\\Project\\go\\src\\github.com\\corbamico\\tbcload\\test\\test.tcl"
	fs, err := os.Open(sFile)
	if err != nil {
		t.Error(err)
		return
	}
	p := NewParser(fs, os.Stdout)
	if err = p.Parse(); err != nil {
		t.Error(err)
	}
}
