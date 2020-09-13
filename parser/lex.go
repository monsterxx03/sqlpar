package parser

import (
	"errors"
)

var debug = false

func init() {
	yyErrorVerbose = true
}

type Lexer struct {
	s      *Scanner
	result Statement
	err    error
}

func NewLexer(sql string) *Lexer {
	s := NewScanner(sql)
	return &Lexer{s: s}
}

func (l *Lexer) Lex(yylval *yySymType) int {
	tok, err := l.s.Scan()
	if err != nil {
		l.Error(err.Error())
	}
	yylval.str = tok.Literal

	return tok.Token
}

func (l *Lexer) Error(s string) {
	l.err = errors.New(s)
}
