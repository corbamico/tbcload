package main

import (
	"encoding/ascii85"
	"fmt"
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
	exchangeEvery5(dst)

	//step 5 map special char
	for index := 0; index < encodeLen; index++ {
		dst[index] = encodeMap[dst[index]-'!']
	}
	return encodeLen
}

func Decode(dst, src []byte) (ndst int) {
	ndst = 0
	//step 0 map special char
	srcCopy := []byte(src)
	for index := 0; index < len(src); index++ {
		srcCopy[index] = decodeMap[src[index]] + '!'
	}

	//step 1 align to 4 bytes,padding as 0
	srcCopy = align5Bytes(srcCopy, '!')

	//step 2 Big-Endian to Little-Endian
	exchangeEvery5(srcCopy)

	//step 3 ascii85 decode
	ndst, _, _ = ascii85.Decode(dst, srcCopy, true)

	//step 4 reverse string every 5 bytes
	exchangeEvery4(dst[:ndst])

	return
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
	switch len(src) % 5 {
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
	if len(src)%4 != 0 {
		return
	}

	for len(src) > 0 {
		src[0], src[1], src[2], src[3] = src[3], src[2], src[1], src[0]
		src = src[4:]
	}
}
func exchangeEvery5(src []byte) {
	if len(src)%5 != 0 {
		return
	}

	for len(src) > 0 {
		src[0], src[1], src[3], src[4] = src[4], src[3], src[1], src[0]
		src = src[5:]
	}
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

const A85_WHITESPACE = 0
const A85_ILLEGAL_CHAR = 0
const A85_Z = 0

var decodeMap = [...]byte{
	A85_ILLEGAL_CHAR, /* ^@ */
	A85_ILLEGAL_CHAR, /* ^A */
	A85_ILLEGAL_CHAR, /* ^B */
	A85_ILLEGAL_CHAR, /* ^C */
	A85_ILLEGAL_CHAR, /* ^D */
	A85_ILLEGAL_CHAR, /* ^E */
	A85_ILLEGAL_CHAR, /* ^F */
	A85_ILLEGAL_CHAR, /* ^G */
	A85_ILLEGAL_CHAR, /* ^H */
	A85_WHITESPACE,   /* \t */
	A85_WHITESPACE,   /* \n */
	A85_ILLEGAL_CHAR, /* ^K */
	A85_ILLEGAL_CHAR, /* ^L */
	A85_ILLEGAL_CHAR, /* ^M */
	A85_ILLEGAL_CHAR, /* ^N */
	A85_ILLEGAL_CHAR, /* ^O */
	A85_ILLEGAL_CHAR, /* ^P */
	A85_ILLEGAL_CHAR, /* ^Q */
	A85_ILLEGAL_CHAR, /* ^R */
	A85_ILLEGAL_CHAR, /* ^S */
	A85_ILLEGAL_CHAR, /* ^T */
	A85_ILLEGAL_CHAR, /* ^U */
	A85_ILLEGAL_CHAR, /* ^V */
	A85_ILLEGAL_CHAR, /* ^W */
	A85_ILLEGAL_CHAR, /* ^X */
	A85_ILLEGAL_CHAR, /* ^Y */
	A85_ILLEGAL_CHAR, /* ^Z */
	A85_ILLEGAL_CHAR, /* ^[ */
	A85_ILLEGAL_CHAR, /* ^\ */
	A85_ILLEGAL_CHAR, /* ^] */
	A85_ILLEGAL_CHAR, /* ^^ */
	A85_ILLEGAL_CHAR, /* ^_ */
	A85_WHITESPACE,   /*   */
	0,                /* ! */
	A85_ILLEGAL_CHAR, /* " (for hilit: ") */
	2,                /* # */
	A85_ILLEGAL_CHAR, /* $ */
	4,                /* % */
	5,                /* & */
	6,                /* ' */
	7,                /* ( */
	8,                /* ) */
	9,                /* * */
	10,               /* + */
	11,               /* , */
	12,               /* - */
	13,               /* . */
	14,               /* / */
	15,               /* 0 */
	16,               /* 1 */
	17,               /* 2 */
	18,               /* 3 */
	19,               /* 4 */
	20,               /* 5 */
	21,               /* 6 */
	22,               /* 7 */
	23,               /* 8 */
	24,               /* 9 */
	25,               /* : */
	26,               /* ; */
	27,               /* < */
	28,               /* = */
	29,               /* > */
	30,               /* ? */
	31,               /* @ */
	32,               /* A */
	33,               /* B */
	34,               /* C */
	35,               /* D */
	36,               /* E */
	37,               /* F */
	38,               /* G */
	39,               /* H */
	40,               /* I */
	41,               /* J */
	42,               /* K */
	43,               /* L */
	44,               /* M */
	45,               /* N */
	46,               /* O */
	47,               /* P */
	48,               /* Q */
	49,               /* R */
	50,               /* S */
	51,               /* T */
	52,               /* U */
	53,               /* V */
	54,               /* W */
	55,               /* X */
	56,               /* Y */
	57,               /* Z */
	A85_ILLEGAL_CHAR, /* [ */
	A85_ILLEGAL_CHAR, /* \ */
	A85_ILLEGAL_CHAR, /* ] */
	61,               /* ^ */
	62,               /* _ */
	63,               /* ` */
	64,               /* a */
	65,               /* b */
	66,               /* c */
	67,               /* d */
	68,               /* e */
	69,               /* f */
	70,               /* g */
	71,               /* h */
	72,               /* i */
	73,               /* j */
	74,               /* k */
	75,               /* l */
	76,               /* m */
	77,               /* n */
	78,               /* o */
	79,               /* p */
	80,               /* q */
	81,               /* r */
	82,               /* s */
	83,               /* t */
	84,               /* u */
	1,                /* v (replaces ") " */
	3,                /* w (replaces $) */
	58,               /* x (replaces [) */
	59,               /* y (replaces \) */
	A85_Z,            /* z */
	A85_ILLEGAL_CHAR, /* { */
	60,               /* | (replaces ]) */
	A85_ILLEGAL_CHAR, /* } */
	A85_ILLEGAL_CHAR, /* ~ */
}

func main() {
	// src := []byte("aBHr@2u|fD!D9(")
	// dst := make([]byte, 104)
	// Decode(dst, src)
	// fmt.Println(dst)
	// fmt.Println(string(dst))

	//00-00-00-00-00
	//src := []byte("已启用临时授权，请尽快测试并反馈。截止2020930将无法再使用临时授权！！！")
	src := []byte("恭喜，您已启用永久授权。感谢您对我们工作成果/版权保护的支持。")
	fmt.Println(len(src))
	//src := []byte("告警")
	dst := make([]byte, 150)
	Encode(dst, src)
	fmt.Println(dst)
	fmt.Println(string(dst))
}
