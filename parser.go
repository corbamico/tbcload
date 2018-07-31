package tbcload

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

//Parser read tbc file and write 'dissemble' to w
type Parser struct {
	r      Decoder
	w      bufio.Writer
	Detail bool //true: disassemble bytecode

	codeBytes  bytes.Buffer
	codeDelta  bytes.Buffer
	codeLength bytes.Buffer
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
	if err = p.parseCode(); err != nil {
		return
	}
	//3. CodeDelta
	if err = p.parseCodeDelta(); err != nil {
		return
	}
	//4. CodeLength
	if err = p.parseCodeLength(); err != nil {
		return
	}

	//if dump all instruction
	if p.Detail {
		if err = p.parseDecompile(); err != nil {
			return
		}
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
		if err = p.parseObject(); err != nil {
			return err
		}
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
	var nLen int64
	if nLen, err = p.parseIntLine(); err != nil {
		return err
	}
	if nLen == 0 {
		p.r.ReadRaw(buf[:])
		return nil
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

	p.w.WriteString("\n---procedure begin---\n")
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
	p.w.WriteString("\n---procedure end  ---")
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
	ss := fmt.Sprintf("[local-%02d]name=%s,hasDefault=%d ", ints[0], sName, ints[1])
	p.w.WriteString(ss)

	//3. if (hasDef) Object
	if ints[1] == 1 {
		err = p.parseObject()
	}
	p.w.WriteByte('\n')
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
		if _, err = p.parseRawStringLine(); err != nil {
			return
		}
		if _, err = p.parseRawStringLine(); err != nil {
			return
		}
		if _, err = p.parseRawStringLine(); err != nil {
			return
		}
		if _, err = p.parseRawStringLine(); err != nil {
			return
		}
	}
	return err
}

func (p *Parser) parseCodeDelta() (err error) {
	var buf [20480]byte
	var nRead int
	var nRes int64
	if nRes, err = p.parseIntLine(); err != nil {
		return
	}
	if nRead, err = p.r.Read(buf[:]); err != nil {
		return
	}
	p.codeDelta.Reset()

	if nRes > 0 && nRead > 0 {
		s := hex.EncodeToString(buf[:nRead])
		_, err = p.w.WriteString(s)
		err = p.w.WriteByte('\n')

		if nRead > int(nRes) {
			nRead = int(nRes)
		}
		p.codeDelta.Write(buf[:nRead])
	}
	return
}
func (p *Parser) parseCodeLength() (err error) {
	var buf [20480]byte
	var nRead int
	var nRes int64
	if nRes, err = p.parseIntLine(); err != nil {
		return
	}
	if nRead, err = p.r.Read(buf[:]); err != nil {
		return
	}

	p.codeLength.Reset()

	if nRes > 0 && nRead > 0 {
		s := hex.EncodeToString(buf[:nRead])
		_, err = p.w.WriteString(s)
		err = p.w.WriteByte('\n')

		if nRead > int(nRes) {
			nRead = int(nRes)
		}
		p.codeLength.Write(buf[:nRead])
	}
	return
}

func (p *Parser) parseCode() (err error) {
	var buf [20480]byte
	var nRead int
	var nRes int64
	if nRes, err = p.parseIntLine(); err != nil {
		return
	}
	if nRead, err = p.r.Read(buf[:]); err != nil {
		return
	}
	p.codeBytes.Reset()

	if nRes > 0 && nRead > 0 {
		s := hex.EncodeToString(buf[:nRead])
		_, err = p.w.WriteString(s)
		err = p.w.WriteByte('\n')

		if nRead > int(nRes) {
			nRead = int(nRes)
		}
		p.codeBytes.Write(buf[:nRead])
	}
	return
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
	// if isCode {
	// 	err = p.parseDisassembleCode(buf[:nRead])
	// }
	return
}

func parseOperand(src []byte, operandType byte) (res string, nextSrc []byte, err error) {
	var nLen = 0
	switch operandType {
	case OPERAND_NONE:
		nLen = 0
	case OPERAND_INT1, OPERAND_LVT1, OPERAND_OFFSET1:
		/* One byte signed integer. */
		res = strconv.Itoa(int(int8(src[0])))
		nLen = 1
	case OPERAND_UINT1, OPERAND_LIT1, OPERAND_SCLS1:
		/* One byte unsigned integer. */
		res = strconv.Itoa(int(uint8(src[0])))
		nLen = 1

	case OPERAND_INT4, OPERAND_IDX4, OPERAND_OFFSET4:
		/* Four byte signed integer. */
		var i int32
		nLen = 4
		buf := bytes.NewBuffer(src[:nLen])
		err = binary.Read(buf, binary.BigEndian, &i)
		res = strconv.Itoa(int(i))

	case OPERAND_UINT4, OPERAND_LVT4, OPERAND_AUX4, OPERAND_LIT4:
		/* Four byte unsigned integer. */
		var i uint32
		nLen = 4
		buf := bytes.NewBuffer(src[:nLen])
		err = binary.Read(buf, binary.BigEndian, &i)
		res = strconv.Itoa(int(i))
		//BUG? int(uint32(i))?
	}

	return res, src[nLen:], err
}

func paresOneOp(src []byte) (res string, numBytes int, err error) {
	var b strings.Builder
	var str string
	opInt := int(src[0])
	if opInt >= len(tclOpTable) {
		return "", 1, os.ErrNotExist
	}
	opName := tclOpTable[opInt].name
	bytes := tclOpTable[opInt].numBytes
	numOperands := tclOpTable[opInt].numOperands
	op1 := tclOpTable[opInt].opTypes[0]
	op2 := tclOpTable[opInt].opTypes[1]
	b.WriteString(opName)
	src = src[1:]

	if numOperands > 0 {
		if str, src, err = parseOperand(src, op1); err != nil {
			return b.String(), bytes, err
		}
		b.WriteByte(' ')
		b.WriteString(str)
	}

	if numOperands > 1 {
		if str, _, err = parseOperand(src, op2); err != nil {
			return b.String(), bytes, err
		}
		b.WriteByte(' ')
		b.WriteString(str)
	}
	return b.String(), bytes, err

}

func (p *Parser) parseDecompile() (err error) {
	var str string
	var bytes int
	var totalBytes int
	src := p.codeBytes.Bytes()
	codeDelta := p.codeDelta.Bytes()
	codeLength := p.codeLength.Bytes()

	numCmds := len(codeDelta)
	indexCmds := 0
	cmdBegin := true
	var cmdBytes, cmdDelta int

	if len(src) == 0 {
		return
	}
	for len(src) > 0 {
		if str, bytes, err = paresOneOp(src); err != nil {
			return err
		}
		//1. print command title: command %d,pc=xx-xx
		if numCmds > 0 && indexCmds < (numCmds) {

			if cmdBegin {
				cmdBytes = bytes
				cmdBegin = false
				//BUG,FIXME, we dont consider codeDelta = 0xFF 4bytes case
				samePCforCmds := true

				for samePCforCmds {

					//cmdDelta = int(p.codeDelta[indexCmds+1])

					p.w.WriteString(fmt.Sprintf("\tCommand %d", indexCmds))
					if indexCmds < len(codeLength) {
						//BUG,FIXME, we dont consider codeLength = 0xFF 4bytes case
						p.w.WriteString(fmt.Sprintf(",pc= %d-%d", totalBytes, totalBytes+int(codeLength[indexCmds])-1))
					}
					p.w.WriteByte('\n')

					indexCmds++
					//cmdDelta 是下一条Command 相对于当前Comand的bytes偏移量
					//如果是最后一条命令，cmdDelta赋值MAX
					if indexCmds < len(codeDelta) {
						cmdDelta = int(codeDelta[indexCmds])
					} else {
						cmdDelta = math.MaxInt32
					}

					if cmdDelta != 0 {
						samePCforCmds = false
					}
				}

				if cmdBytes >= cmdDelta {
					cmdBegin = true
				}

			} else {
				cmdBytes += bytes
				if cmdBytes >= cmdDelta {
					cmdBegin = true
				}
			}
		}

		//2. print command instruction
		p.w.WriteString(fmt.Sprintf("\t(%d)", totalBytes))
		p.w.WriteString(str)
		p.w.WriteByte('\n')
		src = src[bytes:]
		totalBytes += bytes
	}
	return
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
