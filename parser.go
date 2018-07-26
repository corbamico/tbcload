package tbcload

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//Parser read tbc file and write 'dissemble' to w
type Parser struct {
	r Decoder
	w bufio.Writer
}

//NewParser create Parser
func NewParser(r io.Reader, w io.Writer) *Parser {
	return &Parser{r: *NewDecoder(r), w: *bufio.NewWriter(w)}
}

const tbcFileBeginWith = "TclPro ByteCode "

//Parse from io.Reader
func (p *Parser) Parse() (err error) {
	if err = p.skipUntil(tbcFileBeginWith); err != nil {
		return
	}
	err = p.parseByteCode()
	p.w.Flush()
	return
}

func (p *Parser) skipUntil(prefix string) (err error) {
	var buf [maxCharsOneLine]byte
	var nRead int
	for {
		if nRead, err = p.r.ReadRaw(buf[:]); err != nil {
			return
		}
		if strings.HasPrefix(string(buf[:nRead]), prefix) {
			return nil
		}
	}
}

func (p *Parser) parseIntLine() (res int64, err error) {
	var buf [maxCharsOneLine]byte
	var nRead int
	if nRead, err = p.r.ReadRaw(buf[:]); err == nil {
		res, err = strconv.ParseInt(string(buf[:nRead]), 10, 32)
	}
	return
}
func (p *Parser) parseIntList() (res []int64, err error) {
	var buf [maxCharsOneLine]byte
	var nRead int
	var i64 int64
	if nRead, err = p.r.ReadRaw(buf[:]); err != nil {
		return nil, err
	}
	fields := strings.Fields(string(buf[:nRead]))
	result := make([]int64, len(fields))
	for i, s := range fields {
		if i64, err = strconv.ParseInt(s, 10, 32); err != nil {
			return nil, err
		}
		result[i] = i64
	}
	return result, nil
}
func (p *Parser) parseByteCode() (err error) {
	//1. procedure struct info

	if _, err = p.parseRawStringLine(); err != nil {
		return
	}
	//2. ByteCode
	if err = p.parseHex(); err != nil {
		return
	}
	//3. CodeDelta
	if err = p.parseHex(); err != nil {
		return
	}
	//4. CodeLength
	if err = p.parseHex(); err != nil {
		return
	}
	//5. ObjectArray
	if err = p.parseObjectArray(); err != nil {
		return
	}
	//6. ExcRangeArray
	if err = p.parseExcRangeArray(); err != nil {
		return
	}
	//6. AuxDataArray
	err = p.parseAuxDataArray()
	return
}
func (p *Parser) parseObjectArray() (err error) {
	var num int64
	if num, err = p.parseIntLine(); err != nil {
		return err
	}
	for index := 0; index < int(num); index++ {
		//output
		p.w.WriteString(fmt.Sprintf("[lit-%04d]", index))
		p.parseObject()
		p.w.WriteByte('\n')
	}
	return
}
func (p *Parser) parseObjectType() (c byte, err error) {
	var buf [maxCharsOneLine]byte
	if _, err = p.r.ReadRaw(buf[:]); err != nil {
		return 0, err
	}
	return buf[0], err
}

// ErrUnsupoortedObjectType means object type is not correct
var ErrUnsupoortedObjectType = errors.New("object type is not supported")

func (p *Parser) parseObject() (err error) {
	var objType byte
	if objType, err = p.parseObjectType(); err != nil {
		return err
	}
	switch objType {
	case 'i', 'd', 's':
		err = p.parseSimpleObject()
	case 'x':
		err = p.parseXStringObject()
	case 'p':
		err = p.parseProcedureObject()
	default:
		err = ErrUnsupoortedObjectType
	}
	return
}
func (p *Parser) parseSimpleObject() (err error) {
	var buf [maxCharsOneLine]byte
	var nRead int
	if nRead, err = p.r.ReadRaw(buf[:]); err != nil {
		return err
	}

	//output
	p.w.Write(buf[:nRead])
	return
}
func (p *Parser) parseXStringObject() (err error) {
	var buf [20480]byte
	var nRead int
	if _, err = p.parseIntLine(); err != nil {
		return err
	}
	if nRead, err = p.r.Read(buf[:]); err != nil {
		return err
	}

	//output
	p.w.Write(buf[:nRead])
	return
}
func (p *Parser) parseProcedureObject() (err error) {
	var lengths []int64
	//1. ByteCode
	if err = p.parseByteCode(); err != nil {
		return
	}
	//2. numArgs numCompiledLocal
	if lengths, err = p.parseIntList(); err != nil || len(lengths) != 2 {
		return
	}
	//3. for-loop {CompiledLocal}
	for index := 0; index < int(lengths[1]); index++ {
		if err = p.parseCompiledLocal(); err != nil {
			return
		}
	}
	return
}
func (p *Parser) parseCompiledLocal() (err error) {
	var sName string
	var ints []int64
	//1. name
	if _, err = p.parseIntLine(); err != nil {
		return
	}
	if sName, err = p.parseASCII85StringLine(); err != nil {
		return
	}
	//2. index hasDef mask
	if ints, err = p.parseIntList(); err != nil || len(ints) != 3 {
		return
	}
	ss := fmt.Sprintf("[local-%d]%s,hasDefault:%d", ints[0], sName, ints[1])
	p.w.WriteString(ss)

	//3. if (hasDef) Object
	err = p.parseObject()
	return
}
func (p *Parser) parseExcRangeArray() (err error) {
	var nLen int64
	if nLen, err = p.parseIntLine(); err != nil {
		return
	}
	for index := 0; index < int(nLen); index++ {
		p.parseRawStringLine()
	}
	return
}
func (p *Parser) parseAuxDataArray() (err error) {
	//TODO we dont support AuxData parser. later.
	var num int64
	if num, err = p.parseIntLine(); err != nil {
		return
	}
	for index := 0; index < int(num); index++ {
		//we only support CMP_FOREACH_INFO('F')
		//and only for numLists=1,numVars=1
		//F
		//numLists firstValueTemp loopCtTemp
		//numVars
		//*varIndexesPtr
		_, err = p.parseRawStringLine()
		_, err = p.parseRawStringLine()
		_, err = p.parseRawStringLine()
		_, err = p.parseRawStringLine()
	}
	return err
}

//only conver asci85 to hex printing.
func (p *Parser) parseHex() (err error) {
	var buf [20480]byte
	var nRead int
	if _, err = p.parseIntLine(); err != nil {
		return
	}
	if nRead, err = p.r.Read(buf[:]); err != nil {
		return
	}
	s := hex.EncodeToString(buf[:nRead])
	_, err = p.w.WriteString(s)
	err = p.w.WriteByte('\n')
	return
}

//ascii85 decode ,and then disassemble code
func (p *Parser) parseCode() (err error) {
	//TODO disassemble bytecode
	return p.parseHex()
}

func (p *Parser) parseRawStringLine() (str string, err error) {
	var buf [maxCharsOneLine]byte
	var nRead int
	if nRead, err = p.r.ReadRaw(buf[:]); err != nil {
		return "", err
	}
	return string(buf[:nRead]), err
}

func (p *Parser) parseASCII85StringLine() (str string, err error) {
	var buf [maxCharsOneLine]byte
	var nRead int
	if nRead, err = p.r.Read(buf[:]); err != nil {
		return "", err
	}
	return string(buf[:nRead]), err
}
