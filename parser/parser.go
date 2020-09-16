// Code generated by goyacc -o parser/parser.go parser/parser.y. DO NOT EDIT.

//line parser/parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser/parser.y:2

import (
	"github.com/monsterxx03/sqlpar/value"
	"strconv"
)

func setResult(yylex interface{}, stmt Statement) {
	yylex.(*Lexer).result = stmt
}

//line parser/parser.y:15
type yySymType struct {
	yys            int
	str            string
	value          value.Value
	expr           Expr
	stmt           Statement
	sel            *Select
	sel_field      SelectField
	sel_field_list SelectFieldList
	orderBy        OrderBy
	where          *Where
	limit          *Limit
}

const ILLEGAL = 57346
const SELECT = 57347
const FROM = 57348
const WHERE = 57349
const ORDER_BY = 57350
const LIMIT = 57351
const OFFSET = 57352
const IDENT = 57353
const INTEGER = 57354
const FLOAT = 57355
const TRUE = 57356
const FALSE = 57357
const NULL = 57358
const AND = 57359
const OR = 57360
const NOT = 57361
const LE = 57362
const GE = 57363
const NE = 57364
const IS = 57365
const LIKE = 57366
const IN = 57367

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ILLEGAL",
	"SELECT",
	"FROM",
	"WHERE",
	"ORDER_BY",
	"LIMIT",
	"OFFSET",
	"IDENT",
	"INTEGER",
	"FLOAT",
	"TRUE",
	"FALSE",
	"NULL",
	"AND",
	"OR",
	"NOT",
	"'!'",
	"'='",
	"'<'",
	"'>'",
	"LE",
	"GE",
	"NE",
	"IS",
	"LIKE",
	"IN",
	"','",
	"'*'",
	"'('",
	"')'",
	"'\"'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser/parser.y:138

func Parse(s string) (Statement, error) {
	l := NewLexer(s)
	yyParse(l)
	return l.result, l.err
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 10,
	32, 9,
	-2, 11,
}

const yyPrivate = 57344

const yyLast = 64

var yyAct = [...]int{

	47, 48, 49, 50, 51, 55, 56, 29, 30, 27,
	12, 23, 10, 20, 13, 41, 11, 26, 6, 29,
	30, 53, 46, 44, 52, 5, 28, 15, 54, 22,
	24, 16, 7, 19, 4, 40, 31, 1, 39, 17,
	12, 42, 43, 33, 34, 35, 36, 37, 38, 25,
	21, 18, 3, 2, 8, 45, 32, 9, 14, 0,
	0, 0, 8, 8,
}
var yyPact = [...]int{

	29, -1000, -1000, -1000, 1, 10, -1000, -1000, -1000, -18,
	-1000, 16, 1, 1, 26, -1000, -1000, -20, 20, -2,
	-1000, -1000, 14, 2, -2, 22, -2, -1000, 5, -2,
	-2, -10, -12, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	12, 9, -1000, -1000, -1000, -1000, -6, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -28, -1000,
}
var yyPgo = [...]int{

	0, 11, 58, 49, 57, 56, 55, 53, 52, 18,
	25, 51, 50, 37,
}
var yyR1 = [...]int{

	0, 13, 7, 8, 10, 10, 9, 9, 9, 4,
	2, 3, 11, 11, 1, 1, 1, 1, 1, 12,
	12, 12, 12, 6, 6, 6, 6, 6, 6, 6,
	5, 5, 5, 5, 5, 5,
}
var yyR2 = [...]int{

	0, 1, 1, 6, 1, 3, 1, 1, 4, 1,
	1, 1, 0, 2, 3, 3, 3, 3, 2, 0,
	2, 4, 4, 2, 3, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1,
}
var yyChk = [...]int{

	-1000, -13, -7, -8, 5, -10, -9, 31, -3, -4,
	11, 6, 30, 32, -2, 11, -9, -10, -11, 7,
	33, -12, 9, -1, 32, -3, 19, 11, 12, 17,
	18, -1, -5, 21, 22, 23, 24, 25, 26, -1,
	30, 10, -1, -1, 33, -6, 34, 12, 13, 14,
	15, 16, 12, 12, 34, 11, 34,
}
var yyDef = [...]int{

	0, -2, 1, 2, 0, 0, 4, 6, 7, 0,
	-2, 0, 0, 0, 12, 10, 5, 0, 19, 0,
	8, 3, 0, 13, 0, 0, 0, 11, 20, 0,
	0, 0, 0, 30, 31, 32, 33, 34, 35, 18,
	0, 0, 16, 17, 14, 15, 0, 25, 26, 27,
	28, 29, 21, 22, 23, 0, 24,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 20, 34, 3, 3, 3, 3, 3,
	32, 33, 31, 3, 30, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	22, 21, 23,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 24, 25,
	26, 27, 28, 29,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:50
		{
			setResult(yylex, yyDollar[1].stmt)
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:52
		{
			yyVAL.stmt = yyDollar[1].sel
		}
	case 3:
		yyDollar = yyS[yypt-6 : yypt+1]
//line parser/parser.y:56
		{
			yyVAL.sel = NewSelect(yyDollar[2].sel_field_list, yyDollar[4].str, yyDollar[5].where, yyDollar[6].limit)
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:59
		{
			yyVAL.sel_field_list = SelectFieldList{yyDollar[1].sel_field}
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser/parser.y:60
		{
			yyVAL.sel_field_list = append(yyVAL.sel_field_list, yyDollar[3].sel_field)
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:63
		{
			yyVAL.sel_field = &StarExpr{}
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:64
		{
			yyVAL.sel_field = &ColExpr{yyDollar[1].str}
		}
	case 8:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser/parser.y:65
		{
			yyVAL.sel_field = &FuncExpr{Name: yyDollar[1].str, Fields: yyDollar[3].sel_field_list}
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:68
		{
			yyVAL.str = yyDollar[1].str
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:71
		{
			yyVAL.str = yyDollar[1].str
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:74
		{
			yyVAL.str = yyDollar[1].str
		}
	case 12:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser/parser.y:77
		{
			yyVAL.where = nil
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser/parser.y:79
		{
			yyVAL.where = NewWhere(yyDollar[2].expr)
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser/parser.y:83
		{
			yyVAL.expr = yyDollar[2].expr
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser/parser.y:86
		{
			yyVAL.expr = &ComparisonExpr{Left: yyDollar[1].str, Operator: yyDollar[2].str, Right: yyDollar[3].value}
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser/parser.y:88
		{
			yyVAL.expr = &AndExpr{Left: yyDollar[1].expr, Right: yyDollar[3].expr}
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser/parser.y:90
		{
			yyVAL.expr = &OrExpr{Left: yyDollar[1].expr, Right: yyDollar[3].expr}
		}
	case 18:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser/parser.y:92
		{
			yyVAL.expr = &NotExpr{Expr: yyDollar[2].expr}
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
//line parser/parser.y:95
		{
			yyVAL.limit = nil
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser/parser.y:97
		{
			limit, _ := strconv.Atoi(yyDollar[2].str)
			yyVAL.limit = &Limit{Rowcount: limit}
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser/parser.y:102
		{
			offset, _ := strconv.Atoi(yyDollar[2].str)
			limit, _ := strconv.Atoi(yyDollar[4].str)
			yyVAL.limit = &Limit{Offset: offset, Rowcount: limit}
		}
	case 22:
		yyDollar = yyS[yypt-4 : yypt+1]
//line parser/parser.y:108
		{
			limit, _ := strconv.Atoi(yyDollar[2].str)
			offset, _ := strconv.Atoi(yyDollar[4].str)
			yyVAL.limit = &Limit{Offset: offset, Rowcount: limit}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
//line parser/parser.y:115
		{
			yyVAL.value = value.Str{Val: ""}
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
//line parser/parser.y:116
		{
			yyVAL.value = value.Str{Val: yyDollar[2].str}
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:118
		{
			v, _ := strconv.Atoi(yyDollar[1].str)
			yyVAL.value = value.Int{Val: int64(v)}
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:123
		{
			v, _ := strconv.ParseFloat(yyDollar[1].str, 64)
			yyVAL.value = value.Float{Val: v}
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:127
		{
			yyVAL.value = value.Bool{true}
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:128
		{
			yyVAL.value = value.Bool{false}
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:129
		{
			yyVAL.value = value.Null{}
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:132
		{
			yyVAL.str = "="
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:133
		{
			yyVAL.str = "<"
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:134
		{
			yyVAL.str = ">"
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:135
		{
			yyVAL.str = "<="
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:136
		{
			yyVAL.str = ">="
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
//line parser/parser.y:137
		{
			yyVAL.str = "!="
		}
	}
	goto yystack /* stack new state and value */
}
