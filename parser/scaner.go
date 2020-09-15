package parser

import (
	"bytes"
	"strings"
	"unicode"
)

const (
	EOF = -1
)

type Scanner struct {
	src []rune
	pos int
	buf bytes.Buffer
}

var keywords =  map[string]int{
	"SELECT": SELECT,
	"WHERE": WHERE,
	"FROM": FROM,
	"LIMIT": LIMIT,
	"OFFSET": OFFSET,
	"AND": AND, "OR": OR, "NOT": NOT,
}

func NewScanner(sql string) *Scanner {
	return &Scanner{src: []rune(sql)}
}

func (s *Scanner) Scan() (*Token, error) {
	s.ignoreSpace()
	ch := s.next()
	token := ch
	lit := string(ch)
	switch {
	case isDigitChar(ch):
		token = s.scanNumber(ch)
		lit = s.buf.String()
	case isIdentChar(ch):
		token = rune(s.scanIdent(ch))
		lit = s.buf.String()
		_t, ok := keywords[strings.ToUpper(lit)]
		if ok {
			token = rune(_t)
		}
	case isOperatorChar(ch):
		char := s.scanOperator(ch)
		if s.peek() == '=' {
			s.next()
			if char == '>' {
				lit = ">="
				token = GE
			} else if char == '<' {
				lit = "<="
				token = LE
			} else if char == '!' {
				lit = "!="
				token = NE
			}
		}
	}
	return &Token{Token: int(token), Literal: lit}, nil
}

func (s *Scanner) scanNumber(head rune) rune {
	s.buf.Reset()
	s.buf.WriteRune(head)
	for isDigitChar(s.peek()) {
		s.buf.WriteRune(s.next())
	}

	if s.peek() == '.' {
		s.buf.WriteRune(s.next())
		for isDigitChar(s.peek()) {
			s.buf.WriteRune(s.next())
		}
		return FLOAT
	}

	return INTEGER
}

func (s *Scanner) scanIdent(head rune) int {
	s.buf.Reset()
	s.buf.WriteRune(head)
	for isIdentChar(s.peek()) {
		s.buf.WriteRune(s.next())
	}
	return IDENT
}

func (s *Scanner) scanOperator(head rune) rune {
	s.buf.Reset()

	char := head
	s.buf.WriteRune(head)
	for isOperatorChar(s.peek()) {
		char = s.next()
		s.buf.WriteRune(char)
	}
	return char
}

func (s *Scanner) ignoreSpace() {
	for unicode.IsSpace(s.peek()) {
		s.next()
	}
}

func (s *Scanner) next() rune {
	n := s.peek()
	if n != EOF {
		s.pos++
	}
	return n
}

func (s *Scanner) peek() rune {
	if len(s.src) <= s.pos {
		return -1
	}
	return s.src[s.pos]
}

type Token struct {
	Token   int
	Literal string
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigitChar(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentChar(ch rune) bool {
	return isLetter(ch) || isDigitChar(ch) || ch == '_'
}

func isOperatorChar(ch rune) bool {
	return ch == '=' || ch == '>' || ch == '<' || ch == '!'
}

