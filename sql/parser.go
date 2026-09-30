package sql

import (
	"errors"
	"odydb/cell"
	"odydb/schema"
	"strconv"
	"strings"
)

type Parser struct {
	buf string
	pos int
}

var supportedDataTypes = map[string]cell.CellType{
	"string": cell.TypeStr,
	"int64":  cell.TypeI64,
}

func NewParser(s string) Parser {
	return Parser{buf: s, pos: 0}
}

func (p *Parser) ParseStmt() (out interface{}, err error) {
	if p.tryKeyword("SELECT") {
		stmt := &StmtSelect{}
		err = p.parseSelect(stmt)
		out = stmt
	} else if p.tryKeyword("CREATE", "TABLE") {
		stmt := &StmtCreateTable{}
		err = p.parseCreateTable(stmt)
		out = stmt
	} else if p.tryKeyword("INSERT", "INTO") {
		stmt := &StmtInsert{}
		err = p.parseInsert(stmt)
		out = stmt
	} else if p.tryKeyword("UPDATE") {
		stmt := &StmtUpdate{}
		err = p.parseUpdate(stmt)
		out = stmt
	} else if p.tryKeyword("DELETE", "FROM") {
		stmt := &StmtDelete{}
		err = p.parseDelete(stmt)
		out = stmt
	} else {
		err = errors.New("unknown statement")
	}

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (p *Parser) parseDelete(out *StmtDelete) error {
	var ok bool
	if out.Table, ok = p.tryName(); !ok {
		return errors.New("expect table name")
	}

	if err := p.parseWhere(&out.Keys); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseUpdate(out *StmtUpdate) error {
	var ok bool
	if out.Table, ok = p.tryName(); !ok {
		return errors.New("expect table name")
	}

	if ok = p.tryKeyword("set"); !ok {
		return errors.New("expect 'set'")
	}

	for !p.tryKeyword("where") {
		if len(out.Value) > 0 && !p.tryPunctuation(",") {
			return errors.New("expect comma")
		}

		if keyValue, err := p.parseEquality(); err != nil {
			return err
		} else {
			out.Value = append(out.Value, *keyValue)
		}
	}

	for !p.tryPunctuation(";") {
		if len(out.Keys) > 0 && !p.tryKeyword("and") {
			return errors.New("expect 'and'")
		}

		if keyValue, err := p.parseEquality(); err != nil {
			return err
		} else {
			out.Keys = append(out.Keys, *keyValue)
		}
	}

	return nil
}

func (p *Parser) parseEquality() (*NamedCell, error) {
	colName, ok := p.tryName()
	if !ok {
		return nil, errors.New("expect column name")
	}

	if !p.tryPunctuation("=") {
		return nil, errors.New("expect equal sign")
	}

	valCell := cell.Cell{}
	if err := p.parseValue(&valCell); err != nil {
		return nil, err
	}

	return &NamedCell{Column: colName, Value: valCell}, nil
}

func (p *Parser) parseInsert(out *StmtInsert) error {
	var ok bool
	if out.Table, ok = p.tryName(); !ok {
		return errors.New("expect table name")
	}

	if !p.tryKeyword("VALUES") {
		return errors.New("expect 'values'")
	}

	if !p.tryPunctuation("(") {
		return errors.New("expect opening bracket")
	}

	for !p.tryPunctuation(")") {
		if len(out.Value) > 0 && !p.tryPunctuation(",") {
			return errors.New("expect comma")
		}

		cell := cell.Cell{}

		if err := p.parseValue(&cell); err != nil {
			return err
		}

		out.Value = append(out.Value, cell)
	}

	if !p.tryPunctuation(";") {
		return errors.New("expect semicolon")
	}

	return nil
}

func (p *Parser) parseCreateTable(out *StmtCreateTable) error {
	var ok bool
	if out.Table, ok = p.tryName(); !ok {
		return errors.New("expect table name")
	}

	if ok = p.tryPunctuation("("); !ok {
		return errors.New("expect opening bracket")
	}

	// parse columns
	for {
		if len(out.Cols) > 0 && !p.tryPunctuation(",") {
			return errors.New("expect comma")
		}

		if p.tryKeyword("primary", "key") {
			break
		}

		colName, ok := p.tryName()
		if !ok {
			return errors.New("expect column name")
		}

		colType, ok := p.parseDataType()
		if !ok {
			return errors.New("unknown datatype")
		}

		out.Cols = append(out.Cols, schema.Column{Name: colName, Type: colType})
	}

	// parse primary key
	if ok = p.tryPunctuation("("); !ok {
		return errors.New("expect opening bracket")
	}

	for !p.tryPunctuation(")") {
		if len(out.Pkey) > 0 && !p.tryPunctuation(",") {
			return errors.New("expect comma")
		}

		pKey, ok := p.tryName()
		if !ok {
			return errors.New("expect column name")
		}

		out.Pkey = append(out.Pkey, pKey)
	}

	// end
	if ok = p.tryPunctuation(")"); !ok {
		return errors.New("expect closing bracket")
	}
	if ok = p.tryPunctuation(";"); !ok {
		return errors.New("expect semicolon")
	}

	return nil
}

func (p *Parser) parseDataType() (t cell.CellType, ok bool) {
	for k, v := range supportedDataTypes {
		if p.tryKeyword(k) {
			return v, true
		}
	}

	return
}

func (p *Parser) parseSelect(out *StmtSelect) error {
	for !p.tryKeyword("FROM") {
		if len(out.Cols) > 0 && !p.tryPunctuation(",") {
			return errors.New("expect comma")
		}
		if name, ok := p.tryName(); ok {
			out.Cols = append(out.Cols, name)
		} else {
			return errors.New("expect column")
		}
	}

	if len(out.Cols) == 0 {
		return errors.New("expect column list")
	}
	var ok bool
	if out.Table, ok = p.tryName(); !ok {
		return errors.New("expect table name")
	}

	return p.parseWhere(&out.Keys)
}

func (p *Parser) parseWhere(out *[]NamedCell) error {
	if !p.tryKeyword("WHERE") {
		return errors.New("expect WHERE")
	}

	for !p.tryPunctuation(";") {
		if len(*out) > 0 && !p.tryKeyword("AND") {
			return errors.New("expect AND")
		}

		if nCell, err := p.parseEquality(); err != nil {
			return err
		} else {
			*out = append(*out, *nCell)
		}
	}

	return nil
}

func (p *Parser) parseEqual(out *NamedCell) error {
	var ok bool
	out.Column, ok = p.tryName()
	if !ok {
		return errors.New("expect column")
	}
	if !p.tryPunctuation("=") {
		return errors.New("expect punctuation")
	}

	return p.parseValue(&out.Value)
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

		localPos += 1
	}

	intStartIdx := localPos

	for isDigit(p.buf[localPos]) {
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

func (p *Parser) parseString(out *cell.Cell) error {
	startQuote := p.buf[p.pos]
	localPos := p.pos + 1

	var items []byte

	for ; p.buf[localPos] != startQuote && localPos < len(p.buf); localPos += 1 {
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
	p.pos = localPos + 1 // skip closing quote

	return nil
}

func (p *Parser) tryPunctuation(token string) bool {
	p.skipSpaces()
	if !(p.pos+len(token) <= len(p.buf) && p.buf[p.pos:p.pos+len(token)] == token) {
		return false
	}
	p.pos += len(token)
	return true
}

func (p *Parser) tryName() (string, bool) {
	p.skipSpaces()

	if !isNameStart(p.buf[p.pos]) {
		return "", false
	}

	localPos := p.pos
	nameStartIdx := p.pos

	for localPos < len(p.buf) && isNameContinue(p.buf[localPos]) {
		localPos += 1
	}

	p.pos = localPos
	return p.buf[nameStartIdx:localPos], true
}

func (p *Parser) tryKeyword(kws ...string) bool {
	p.skipSpaces()
	localPos := p.pos

	for _, kw := range kws {
		for localPos < len(p.buf) && isSpace(p.buf[localPos]) {
			localPos += 1
		}

		kwStartIdx := localPos

		for localPos < len(p.buf) && !isSeparator(p.buf[localPos]) {
			localPos++
		}

		parsedKw := strings.ToLower(p.buf[kwStartIdx:localPos])

		if strings.ToLower(kw) != parsedKw {
			return false
		}
	}

	p.pos = localPos
	return true
}

func (p *Parser) isEnd() bool {
	p.skipSpaces()

	return p.pos == len(p.buf)
}

func (p *Parser) skipSpaces() {
	for p.pos < len(p.buf) && isSpace(p.buf[p.pos]) {
		p.pos += 1
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
