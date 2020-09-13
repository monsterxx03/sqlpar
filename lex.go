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
	preTok rune
}

func NewLexer(sql string) *Lexer {
	var s scanner.Scanner
	s.Init(strings.NewReader(sql))
	return &Lexer{s: s}
}

func (l *Lexer) Lex(yylval *yySymType) int {
	for tok := l.s.Scan(); tok != scanner.EOF; tok = l.s.Scan() {
		TOKEN := ILLEGAL
		text := l.s.TokenText()
		yylval.str = text
		switch tok {
		case scanner.Int:
			TOKEN = INTEGRAL
		case '"', ',', '*', '(', ')':
			TOKEN = int(tok)
		// TODO handle > < = >= <= !=
		case scanner.Ident:
			text = strings.ToUpper(text)
			switch text {
			case "SELECT":
				TOKEN = SELECT
			case "FROM":
				TOKEN = FROM
			case "LIMIT":
				TOKEN = LIMIT
			case "WHERE":
				TOKEN = WHERE
			default:
				TOKEN = IDENT
			}
		default:
			l.preTok = tok
			return TOKEN
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
