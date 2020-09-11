package main

import (
	"log"
	"strings"
	"text/scanner"
)

var debug = false

type Lexer struct {
	s      scanner.Scanner
	result Statement
}

func NewLexer(sql string) *Lexer {
	var s scanner.Scanner
	s.Init(strings.NewReader(sql))
	return &Lexer{s: s}
}

func (l *Lexer) Lex(yylval *yySymType) int {
	for tok := l.s.Scan(); tok != scanner.EOF; tok = l.s.Scan() {
		text := l.s.TokenText()
		yylval.str = text
		switch tok {
		case scanner.Int:
			return INTEGRAL
		case ',', '*', '(', ')':
			return int(tok)
		case scanner.Ident:
			text = strings.ToUpper(text)
			switch text {
			case "SELECT":
				return SELECT
			case "FROM":
				return FROM
			case "LIMIT":
				return LIMIT
			default:
				return IDENT
			}
		default:
			return ILLEGAL
		}
	}
	return 0
}

func (l *Lexer) Error(s string) {
	log.Printf("parse error: %s", s)
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentChar(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_'
}
