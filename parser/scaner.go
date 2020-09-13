package parser

import (
	"bytes"
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
		switch lit {
		case "SELECT":
			token = SELECT
		}
	case isOperatorChar(ch):
		s.scanOperator(ch)
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

func (s *Scanner) scanOperator(head rune) {
	s.buf.Reset()

	s.buf.WriteRune(head)
	for isOperatorChar(s.peek()) {
		s.buf.WriteRune(s.next())
	}
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
