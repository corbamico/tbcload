package tbcload

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
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

func ExampleChainReader() {
	r1 := strings.NewReader("1234\n5678\n90\n12\n345")
	r2 := newLineReader(r1, 4)
	//r3 := &eatLastNewLineReader{wrapped: r2}
	r4 := bufio.NewReader(r2)
	s1, _, _ := r4.ReadLine()
	fmt.Printf("%s\n", s1)
	s1, _, _ = r4.ReadLine()
	fmt.Printf("%s\n", s1)
	s1, _, _ = r4.ReadLine()
	fmt.Printf("%s\n", s1)
	// Output:
	// 1234567890
	// 12
	// 345
}

func ExampleChainReader2() {
	r1 := strings.NewReader("1234\n5678\n90\n12\n345")
	r2 := newLineReader(r1, 4)
	r3 := &eatLastNewLineReader{wrapped: r2}
	//r4 := bufio.NewReader(r2)
	var buf [128]byte
	s1, _ := r3.Read(buf[:])
	fmt.Printf("%s\n", buf[:s1])
	s1, _ = r3.Read(buf[:])
	fmt.Printf("%s\n", buf[:s1])
	s1, _ = r3.Read(buf[:])
	fmt.Printf("%s\n", buf[:s1])
	// Output:
	// 1234567890
	// 12
	// 345
}
func ExampleEncode() {
	//src := []byte("testing aliases for non-existent targets")
	src := []byte{0, 0, 0, 0}
	dst := make([]byte, 280)
	length := Encode(dst, src)
	fmt.Printf("%s", dst[:length])
	// Output:
	// ,CHr@
}

func ExampleDecode() {
	src := []byte("4;,>!?.EH&ih-(!e2xi6zyiE<!22>:v35>:v22Ppv2j:U!*|yTv0#>6#5cSs!)'!")
	//src := []byte("z")
	dst := make([]byte, 280)
	length := Decode(dst, src)
	fmt.Printf("%s", string(dst[:length]))
	// Output:
	// ,CHr@
}
