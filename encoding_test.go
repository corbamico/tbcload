package tbcload

import (
	"bytes"
	"fmt"
	"testing"
)

type testVector struct {
	src     string
	encoded string
}

var testData = []testVector{
	{"", ""},
	{"proc", ",CHr@"},
	{"button", "7YDEFTnw"},
	{"-text", "Kt(yG@v"},
	{"Hello TclPro", "RZ!iChROo@jZSfD"},
	{"cbk_clicked", "y+aY?hafq@VY|+"},
	{"tbcload::bcproc", "rpwhC;Z2b3<?<+EfqT+"},
}

func testEncode(t *testing.T, v []testVector) {
	dst := make([]byte, 1024)
	for index := 0; index < len(v); index++ {
		ndst := Encode(dst, []byte(v[index].src))

		if !bytes.Equal(dst[:ndst], []byte(v[index].encoded)) {
			t.Errorf("Encode Error,expected :{%s}, output: {%s}", v[index].encoded, dst[:ndst])
		}
	}
}
func testDecode(t *testing.T, v []testVector) {
	dst := make([]byte, 1024)
	for index := 0; index < len(v); index++ {
		ndst := Decode(dst, []byte(v[index].encoded))

		if !bytes.Equal(dst[:ndst], []byte(v[index].src)) {
			t.Errorf("Decode Error,expected :{%s}, output: {%s}", v[index].src, dst[:ndst])
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	testEncode(t, testData)
	testDecode(t, testData)
}

func ExampleEncode() {
	src := []byte("proc")
	dst := make([]byte, 150)
	length := Encode(dst, src)
	fmt.Printf("%s", dst[:length])
	// Output:
	// ,CHr@
}
