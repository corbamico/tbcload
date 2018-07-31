package tbcload

import (
	"bufio"
	"bytes"
	"encoding/ascii85"
	"errors"
	"io"
)

/*
 * Encoder
 */

// Encode encodes src into at most MaxEncodedLen(len(src))
// bytes of dst, returning the actual number of bytes written.
//
// The encoding handles 4-byte chunks, using a special encoding
// for the last fragment, so Encode is not appropriate for use on
// individual blocks of a large data stream. Use NewEncoder() instead.
func Encode(dst, src []byte) int {
	if len(src) == 0 {
		return 0
	}
	//step 1 align to 4 bytes,padding as 0
	srcCopy := []byte(src)
	srcCopy = align4Bytes(srcCopy, 0)

	//step 2 Big-Endian to Little-Endian
	exchangeEvery4(srcCopy)

	//step 3 ascii85 encode
	encodeLen := ascii85.Encode(dst, srcCopy)

	//step 4 reverse string every 5 bytes
	exchangeEvery5(dst[:encodeLen])

	//step 5 map special char
	for index := 0; index < encodeLen; index++ {
		//there is 'z'(122) special for 0x00000000
		if n := int(dst[index] - '!'); n < len(encodeMap) {
			dst[index] = encodeMap[n]
		}
	}
	//if 0x00 00 00 00->'z', we cannot known how many '!' drop
	return encodeLen
	//return encodeLen - (len(srcCopy) - len(src))
}

/*
 * Decoder
 */

// Decode encodes src into at most len(src)
// bytes of dst, returning the actual number of bytes written.
//
// The encoding handles 5-byte chunks, using a special encoding
// for the last fragment, so Encode is not appropriate for use on
// individual blocks of a large data stream. Use NewEncoder() instead.
func Decode(dst, src []byte) (ndst int) {
	ndst = 0
	//step 0 map special char
	srcCopy := []byte(src)
	for index := 0; index < len(src); index++ {
		srcCopy[index] = decodeMap[src[index]] + '!'
	}

	//step 1 align to 5 bytes,padding as 0 (but not for z)
	//no need do this
	//BUGBUG, if there is 'z' in middle
	srcCopy = align5Bytes(srcCopy, '!')

	//step 2 Big-Endian to Little-Endian(but not for 'z')
	exchangeEvery5(srcCopy)

	//step 3 ascii85 decode
	ndst, _, _ = ascii85.Decode(dst, srcCopy, true)

	//step 4 reverse string every 4 bytes
	exchangeEvery4(dst[:ndst])

	//step 5 drop padding
	if padding := len(srcCopy) - len(src); padding == 4 {
		ndst = ndst - 2
	} else {
		ndst = ndst - padding
	}
	return
}

//Decoder wrap decode for stream reader
type Decoder struct {
	wrapped io.Reader
	srcbuf  [20480]byte
	dstbuf  [20480]byte
	src     []byte
	dst     []byte
	//we did not record error, is that ok?
}

const maxCharsOneLine = 72

//NewDecoder return Decoder which wrap Decode for stream reader
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{wrapped: &eatLastNewLineReader{wrapped: newLineReader(r, maxCharsOneLine)}}
}

// ErrDecodeErr mean error while decoding from bytes
var ErrDecodeErr = errors.New("error decoding from bytes")

func (d *Decoder) Read(p []byte) (nRead int, err error) {
	var nsrc int
	for {
		//we ask more buffer
		if len(d.dst) > len(p) {
			return 0, io.ErrShortBuffer
		}
		//we have data, read from it
		if len(d.dst) > 0 {
			nRead = copy(p, d.dst)
			d.dst = d.dst[nRead:]
			return
		}
		//we dont have data, so read it frist
		nsrc, err = d.wrapped.Read(d.srcbuf[:])
		//we read zero byte
		if nsrc == 0 {
			return nsrc, err
		}

		d.src = d.srcbuf[:nsrc]
		nRead = Decode(d.dstbuf[:], d.src)
		d.dst = d.dstbuf[:nRead]
		if nRead == 0 {
			return 0, ErrDecodeErr
		}
		//we read&decode into d.dst ,so back to for-loop return
	}
}

//ReadRaw only read  from wrapped io without Decoding
func (d *Decoder) ReadRaw(p []byte) (nRead int, err error) {
	nRead, err = d.wrapped.Read(p)
	return
}

type eatLastNewLineReader struct {
	wrapped io.Reader
}

func (r *eatLastNewLineReader) Read(p []byte) (nRead int, err error) {
	nRead, err = r.wrapped.Read(p)
	if nRead > 0 && p[nRead-1] == '\n' {
		nRead--
	}
	if nRead > 0 && p[nRead-1] == '\r' {
		nRead--
	}
	return
}

//numCharsLineReader Read
//implement continue read if size of line eq numChars, until
//size of line less than numChars
type numCharsLineReader struct {
	wrapped   bufio.Reader
	numChars  int //number of each line
	readError error
	lastStr   string
}

func newLineReader(r io.Reader, numChars int) io.Reader {
	return &numCharsLineReader{wrapped: *bufio.NewReader(r), numChars: numChars}
}
func (r *numCharsLineReader) Read(p []byte) (nRead int, err error) {
	var line string

	//len(p)至少需要FixedSize of line大小
	if len(p) < r.numChars {
		return 0, io.ErrShortBuffer
	}
	for {
		//从buffered string中直接返回
		if len(r.lastStr) > 0 {
			nRead = copy(p, r.lastStr)
			if nRead >= len(r.lastStr) {
				//返回了整个buffered string，清空buffer和readError
				r.lastStr = ""
				r.readError = nil
				return nRead, r.readError
			}
			//buffered string后移，以便下次调用Read时再读
			r.lastStr = r.lastStr[nRead:]
			return nRead, nil
		}

		//buffer string为空，需要从wrapped里面Read到buffer中
		var b bytes.Buffer
		var bRead = true

		for bRead {
			line, err = r.wrapped.ReadString('\n') //includes '\n'
			nLen := len(line)
			//if char[72]+"\r\n" || char[72] + "\n" ,则继续读下一行
			if nLen == (r.numChars+2) && err == nil {
				//bRead = true
			} else if nLen == (r.numChars+1) && line[nLen-1] == '\n' && line[nLen-2] != '\r' && err == nil {
				//bRead = true
			} else {
				bRead = false
			}

			//如果line之间的'\n'，则删除；保留String末尾的'\r' '\n'
			if nLen >= r.numChars {
				b.WriteString(line[:r.numChars])
			} else if nLen > 0 {
				b.WriteString(line)
			}
		}

		//已经连续读完，如果读到空，直接返回
		if b.Len() == 0 {
			return 0, err
		}

		r.readError = err
		r.lastStr = b.String()
		//保存到buffer中，返回for循环
	}
}

func align4Bytes(src []byte, padding byte) []byte {
	switch len(src) % 4 {
	case 1:
		src = append(src, padding)
		fallthrough
	case 2:
		src = append(src, padding)
		fallthrough
	case 3:
		src = append(src, padding)
		fallthrough
	case 0:
		return src
	}
	return src
}

func align5Bytes(src []byte, padding byte) []byte {
	//'z' dont need align to 5 bytes
	numZ := bytes.Count(src, []byte{'z'})
	nLen := len(src) - numZ

	switch nLen % 5 {
	case 1:
		src = append(src, padding)
		fallthrough
	case 2:
		src = append(src, padding)
		fallthrough
	case 3:
		src = append(src, padding)
		fallthrough
	case 4:
		src = append(src, padding)
		fallthrough
	case 0:
		return src
	}
	return src
}

func exchangeEvery4(src []byte) {
	//BUG fixme,there is some len(src)%4!=0
	// if len(src)%4 != 0 {
	// 	return
	// }

	for len(src) >= 4 {
		src[0], src[1], src[2], src[3] = src[3], src[2], src[1], src[0]
		src = src[4:]
	}
}
func exchangeEvery5(src []byte) {
	//special 'z' for 0x00000000, dont exchange 'z'
	for len(src) > 0 {
		if src[0] == 'z' {
			src = src[1:]
		} else if len(src) >= 5 {
			src[0], src[1], src[3], src[4] = src[4], src[3], src[1], src[0]
			src = src[5:]
		} else {
			return
		}
	}
	return
}

var encodeMap = [...]byte{
	'!',  /*  0: ! */
	'v',  /*  1: was ", is now v (and this is for hilit:") */
	'#',  /*  2: # */
	'w',  /*  3: was $, is now w */
	'%',  /*  4: % */
	'&',  /*  5: & */
	'\'', /*  6: ' */
	'(',  /*  7: ( */
	')',  /*  8: ) */
	'*',  /*  9: * */
	'+',  /* 10: + */
	',',  /* 11: , */
	'-',  /* 12: - */
	'.',  /* 13: . */
	'/',  /* 14: / */
	'0',  /* 15: 0 */
	'1',  /* 16: 1 */
	'2',  /* 17: 2 */
	'3',  /* 18: 3 */
	'4',  /* 19: 4 */
	'5',  /* 20: 5 */
	'6',  /* 21: 6 */
	'7',  /* 22: 7 */
	'8',  /* 23: 8 */
	'9',  /* 24: 9 */
	':',  /* 25: : */
	';',  /* 26: ; */
	'<',  /* 27: < */
	'=',  /* 28: = */
	'>',  /* 29: > */
	'?',  /* 30: ? */
	'@',  /* 31: @ */
	'A',  /* 32: A */
	'B',  /* 33: B */
	'C',  /* 34: C */
	'D',  /* 35: D */
	'E',  /* 36: E */
	'F',  /* 37: F */
	'G',  /* 38: G */
	'H',  /* 39: H */
	'I',  /* 40: I */
	'J',  /* 41: J */
	'K',  /* 42: K */
	'L',  /* 43: L */
	'M',  /* 44: M */
	'N',  /* 45: N */
	'O',  /* 46: O */
	'P',  /* 47: P */
	'Q',  /* 48: Q */
	'R',  /* 49: R */
	'S',  /* 50: S */
	'T',  /* 51: T */
	'U',  /* 52: U */
	'V',  /* 53: V */
	'W',  /* 54: W */
	'X',  /* 55: X */
	'Y',  /* 56: Y */
	'Z',  /* 57: Z */
	'x',  /* 58: was [, is now x */
	'y',  /* 59: was \, is now y */
	'|',  /* 60: was ], is now | */
	'^',  /* 61: ^ */
	'_',  /* 62: _ */
	'`',  /* 63: ` */
	'a',  /* 64: a */
	'b',  /* 65: b */
	'c',  /* 66: c */
	'd',  /* 67: d */
	'e',  /* 68: e */
	'f',  /* 69: f */
	'g',  /* 70: g */
	'h',  /* 71: h */
	'i',  /* 72: i */
	'j',  /* 73: j */
	'k',  /* 74: k */
	'l',  /* 75: l */
	'm',  /* 76: m */
	'n',  /* 77: n */
	'o',  /* 78: o */
	'p',  /* 79: p */
	'q',  /* 80: q */
	'r',  /* 81: r */
	's',  /* 82: s */
	't',  /* 83: t */
	'u',  /* 84: u */
}

const a85Whitespace = byte(0xff)  //-1
const a85IllegalChar = byte(0xfe) //-2
const a85Z = byte(0xfd)           //-3

var decodeMap = [...]byte{
	a85IllegalChar, /* ^@ */
	a85IllegalChar, /* ^A */
	a85IllegalChar, /* ^B */
	a85IllegalChar, /* ^C */
	a85IllegalChar, /* ^D */
	a85IllegalChar, /* ^E */
	a85IllegalChar, /* ^F */
	a85IllegalChar, /* ^G */
	a85IllegalChar, /* ^H */
	a85Whitespace,  /* \t */
	a85Whitespace,  /* \n */
	a85IllegalChar, /* ^K */
	a85IllegalChar, /* ^L */
	a85IllegalChar, /* ^M */
	a85IllegalChar, /* ^N */
	a85IllegalChar, /* ^O */
	a85IllegalChar, /* ^P */
	a85IllegalChar, /* ^Q */
	a85IllegalChar, /* ^R */
	a85IllegalChar, /* ^S */
	a85IllegalChar, /* ^T */
	a85IllegalChar, /* ^U */
	a85IllegalChar, /* ^V */
	a85IllegalChar, /* ^W */
	a85IllegalChar, /* ^X */
	a85IllegalChar, /* ^Y */
	a85IllegalChar, /* ^Z */
	a85IllegalChar, /* ^[ */
	a85IllegalChar, /* ^\ */
	a85IllegalChar, /* ^] */
	a85IllegalChar, /* ^^ */
	a85IllegalChar, /* ^_ */
	a85Whitespace,  /*   */
	0,              /* ! */
	a85IllegalChar, /* " (for hilit: ") */
	2,              /* # */
	a85IllegalChar, /* $ */
	4,              /* % */
	5,              /* & */
	6,              /* ' */
	7,              /* ( */
	8,              /* ) */
	9,              /* * */
	10,             /* + */
	11,             /* , */
	12,             /* - */
	13,             /* . */
	14,             /* / */
	15,             /* 0 */
	16,             /* 1 */
	17,             /* 2 */
	18,             /* 3 */
	19,             /* 4 */
	20,             /* 5 */
	21,             /* 6 */
	22,             /* 7 */
	23,             /* 8 */
	24,             /* 9 */
	25,             /* : */
	26,             /* ; */
	27,             /* < */
	28,             /* = */
	29,             /* > */
	30,             /* ? */
	31,             /* @ */
	32,             /* A */
	33,             /* B */
	34,             /* C */
	35,             /* D */
	36,             /* E */
	37,             /* F */
	38,             /* G */
	39,             /* H */
	40,             /* I */
	41,             /* J */
	42,             /* K */
	43,             /* L */
	44,             /* M */
	45,             /* N */
	46,             /* O */
	47,             /* P */
	48,             /* Q */
	49,             /* R */
	50,             /* S */
	51,             /* T */
	52,             /* U */
	53,             /* V */
	54,             /* W */
	55,             /* X */
	56,             /* Y */
	57,             /* Z */
	a85IllegalChar, /* [ */
	a85IllegalChar, /* \ */
	a85IllegalChar, /* ] */
	61,             /* ^ */
	62,             /* _ */
	63,             /* ` */
	64,             /* a */
	65,             /* b */
	66,             /* c */
	67,             /* d */
	68,             /* e */
	69,             /* f */
	70,             /* g */
	71,             /* h */
	72,             /* i */
	73,             /* j */
	74,             /* k */
	75,             /* l */
	76,             /* m */
	77,             /* n */
	78,             /* o */
	79,             /* p */
	80,             /* q */
	81,             /* r */
	82,             /* s */
	83,             /* t */
	84,             /* u */
	1,              /* v (replaces ") " */
	3,              /* w (replaces $) */
	58,             /* x (replaces [) */
	59,             /* y (replaces \) */
	89,             /* z->z,special for 0x00000000 */
	//	a85Z,       /* z */
	a85IllegalChar, /* { */
	60,             /* | (replaces ]) */
	a85IllegalChar, /* } */
	a85IllegalChar, /* ~ */
}
