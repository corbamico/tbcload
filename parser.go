package tbcload

import (
	"bufio"
	"io"
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

//Parse from io.Reader
func (p *Parser) Parse() (err error) {
	return
}
func (p *Parser) parseIntLine() (res int64, err error) {
	var buf [maxCharsOneLine]byte
	nRead, err := p.r.ReadRaw(buf[:])
	parseInt()
	return
}
func (p *Parser) parseIntList() (err error) {
	return
}
func (p *Parser) parseByteCode() (err error) {
	return
}
func (p *Parser) parseObjectArray() (err error) {
	return
}
func (p *Parser) parseObject() (err error) {
	return
}
func (p *Parser) parseSimpleObject() (err error) {
	return
}
func (p *Parser) parseXStringObject() (err error) {
	return
}
func (p *Parser) parseProcedureObject() (err error) {
	return
}
func (p *Parser) parseCompiledLocal() (err error) {
	return
}
func (p *Parser) parseExcRangeArray() (err error) {
	return
}
func (p *Parser) parseAuxDataArray() (err error) {
	return
}

//only conver asci85 to hex printing.
func (p *Parser) parseHex() (err error) {
	return
}

//ascii85 decode ,and then disassemble code
func (p *Parser) parseCode() (err error) {
	return
}
