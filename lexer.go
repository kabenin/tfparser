/*
	This file contains simple lexical extractors for tfparser
*/

package tfparser

import (
	"fmt"
	"strings"
)

func (p *parser) skipComment() error {
	if p.i >= len(p.data) {
		return nil
	}
	if p.data[p.i] == '#' {
		p.skipTillEOL()
		return nil
	}
	if p.data[p.i] == '/' {
		if p.i+1 > len(p.data) {
			p.err = fmt.Errorf("Unexpected end of file after '/'")
			return p.err
		}
		// Single line comment starting with // ?
		if p.data[p.i+1] == '/' {
			p.skipTillEOL()
			return nil
		}
		return p.skipMulitlineComment()
	}
	return nil
}

func (p *parser) skipTillEOL() {
	for ; p.i < len(p.data) && p.data[p.i] != '\n'; p.i++ {
	}
	p.i++
	p.skipNewLine()
}

func (p *parser) skipMulitlineComment() error {
	if p.i > len(p.data) {
		return nil
	}
	if p.i+1 >= len(p.data) || (p.data[p.i] != '/' && p.data[p.i+1] != '*') {
		return nil
	}
	nestedCount := 0
	for ; p.i+1 < len(p.data); p.i++ {
		if p.data[p.i] == '/' && p.data[p.i+1] == '*' {
			nestedCount++
		}
		if p.data[p.i] == '*' && p.data[p.i+1] == '/' {
			p.i += 2
			nestedCount--
			if nestedCount == 0 {
				p.skipNewLine()
				return nil
			}
		}
	}
	p.err = fmt.Errorf("Unable to find closing multiline comment")
	return p.err
}

var newlines = "\n\r"

func (p *parser) skipNewLine() {
	if p.i >= len(p.data) {
		return
	}
	for ; p.i < len(p.data) && strings.Contains(newlines, string(p.data[p.i])); p.i++ {
	}
}

// Func skips block till the corresponding closing bracket
func (p *parser) skipBlock() error {
	if p.i >= len(p.data) || p.data[p.i] != '{' {
		return nil
	}
	bracesBalance := 0
	for ; p.i < len(p.data); p.i++ {
		if p.data[p.i] == '{' {
			bracesBalance++
		}
		if p.data[p.i] == '}' {
			bracesBalance--
		}
		if bracesBalance == 0 {
			p.i++
			p.skipNewLine()
			return nil
		}
		p.skipComment()
	}
	if bracesBalance != 0 {
		p.err = fmt.Errorf("Unable to find closing brace for block")
		return p.err
	}
	return nil
}

func (p *parser) peek() string {
	p.popWhitespaces()
	peeked, _ := p.peekWithLength()
	return peeked
}

func (p *parser) pop() string {
	p.popWhitespaces()
	peeked, len := p.peekWithLength()
	p.i += len
	p.popWhitespaces()
	return peeked
}

func (p *parser) popToken(tok string) error {
	p.popWhitespaces()
	t := p.pop()
	if t != tok {
		p.err = fmt.Errorf("Unexpected token %#v, expected %#q", t, tok)
		return p.err
	}
	return nil
}

// gets next 'word' from data and returns it along with its length
func (p *parser) peekWithLength() (string, int) {
	if p.i >= len(p.data) {
		return "", 0
	}
	if p.data[p.i] == '"' {
		return p.peekQuotedStringWithLength()
	}
	if strings.Contains(symbols, string(p.data[p.i])) {
		return string(p.data[p.i]), 1
	}
	return p.peekIdentifierWithLength()
}

func (p *parser) peekQuotedStringWithLength() (string, int) {
	if p.i >= len(p.data) || p.data[p.i] != '"' {
		return "", 0
	}
	for i := p.i + 1; i < len(p.data); i++ {
		if p.data[i] == '"' && p.data[i-1] != '\\' {
			return p.data[p.i+1 : i], i - p.i + 1 //len(p.data[p.i+1:i]) + 2
		}
	}
	return "", 0
}

var symbols = "{}=\""

func (p *parser) peekIdentifierWithLength() (string, int) {
	for i := p.i; i <= len(p.data); i++ {
		if strings.Contains(symbols+whitespaces, string(p.data[i])) {
			return p.data[p.i:i], i - p.i //len(p.data[p.i:i])
		}
	}
	return "", 0
}

var whitespaces = "\t\r\n "

func (p *parser) popWhitespaces() {
	i := p.i
	p.skipComment()
	for ; p.i < len(p.data) && strings.Contains(whitespaces, string(p.data[p.i])); p.i++ {
	}
	if i != p.i {
		p.popWhitespaces()
	}
}
