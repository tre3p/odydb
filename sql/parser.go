package sql

import "strings"

type Parser struct {
	buf string
	pos int
}

func NewParser(s string) Parser {
	return Parser{buf: s, pos: 0}
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