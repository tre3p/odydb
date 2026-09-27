package sql

import (
	"errors"
	"odydb/cell"
	"strconv"
	"strings"
)

type Parser struct {
	buf string
	pos int
}

func NewParser(s string) Parser {
	return Parser{buf: s, pos: 0}
}

func (p *Parser) parseValue(out *cell.Cell) error {
	p.skipSpaces()
	if p.pos >= len(p.buf) {
		return errors.New("expect value")
	}

	ch := p.buf[p.pos]
	if ch == '"' || ch == '\'' {
		return p.parseString(out)
	} else if isDigit(ch) || ch == '-' || ch == '+' {
		return p.parseInt(out)
	} else {
		return errors.New("expect value")
	}
}

func (p *Parser) parseInt(out *cell.Cell) error {
	localPos := p.pos
	isNegative := false

	if !isDigit(p.buf[localPos]) {
		if p.buf[localPos] == '-' {
			isNegative = true
		}

		localPos+=1
	}

	intStartIdx := localPos

	for ;isDigit(p.buf[localPos]); {
		localPos += 1
	}

	if intStartIdx == localPos {
		return errors.New("expect numeric value")
	}

	parsedInt, err := strconv.Atoi(p.buf[intStartIdx:localPos])
	if err != nil {
		return err
	}

	if isNegative {
		parsedInt = -parsedInt
	}

	out.Type = cell.TypeI64
	out.I64 = int64(parsedInt)
	p.pos = localPos

	return nil
}

func (p *Parser) parseString (out *cell.Cell) error {
	startQuote := p.buf[p.pos]
	localPos := p.pos+1

	var items []byte

	for ;p.buf[localPos] != startQuote && localPos < len(p.buf); localPos += 1 {
		if p.buf[localPos] == '\\' {
			items = append(items, p.buf[localPos+1])
			localPos += 1
		} else {
			items = append(items, p.buf[localPos])
		}
	}

	if p.buf[localPos] != startQuote {
		return errors.New("expected closing quote")
	}

	out.Type = cell.TypeStr
	out.Str = items
	p.pos = localPos+1 // skip closing quote

	return nil
}

func (p *Parser) tryName() (string, bool) {
	p.skipSpaces()

	if !isNameStart(p.buf[p.pos]) {
		return "", false
	}

	localPos := p.pos
	nameStartIdx := p.pos

	for localPos < len(p.buf) && isNameContinue(p.buf[localPos]) {
		localPos+=1
	}

	p.pos = localPos
	return p.buf[nameStartIdx:localPos], true
}

func (p *Parser) tryKeyword(kw string) bool {
	p.skipSpaces()

	localPos := p.pos
	kwStartIdx := p.pos

	for localPos < len(p.buf) && !isSeparator(p.buf[localPos])  {
		localPos++
	}

	parsedKw := strings.ToLower(p.buf[kwStartIdx:localPos])

	if strings.ToLower(kw) == parsedKw {
		p.pos = localPos
		return true
	} else {
		return false
	}
}

func (p *Parser) isEnd() bool {
	p.skipSpaces()

	return p.pos == len(p.buf)
}

func (p *Parser) skipSpaces() {
	for p.pos < len(p.buf) && isSpace(p.buf[p.pos]) {
		p.pos+=1
	}
}

func isSeparator(ch byte) bool {
	return ch < 128 && !isNameContinue(ch)
}

func isSpace(ch byte) bool {
	switch ch {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}

	return false
}

func isAlpha(ch byte) bool {
	return 'a' <= (ch|32) && (ch|32) <= 'z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isNameStart(ch byte) bool {
	return isAlpha(ch) || ch == '_'
}

func isNameContinue(ch byte) bool {
	return isAlpha(ch) || isDigit(ch) || ch == '_'
}