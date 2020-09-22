%{
package parser

import (
    "strconv"
    "github.com/monsterxx03/sqlpar/value"
)

func setResult(yylex interface{}, stmt Statement) {
    yylex.(*Lexer).result = stmt
}

%}

%union {
    str             string
    value           value.Value
    expr            Expr
    stmt            Statement
    sel             *Select
    sel_field       SelectField
    sel_field_list  SelectFieldList
    orderBy         OrderBy
    where           *Where
    limit           *Limit
}

%token ILLEGAL
%token <str> SELECT FROM WHERE ORDER_BY LIMIT OFFSET SHOW TABLE
%token <str> IDENT INTEGER FLOAT TRUE FALSE NULL

%type <expr> expr
%type <str> table_name col func_name compare
%type <value> value
%type <stmt> command show_table_stmt
%type <sel> select_stmt
%type <sel_field> sel_field nest_col
%type <sel_field_list> sel_field_list
%type <where> where_opt
%type <limit> limit_opt

%left <str> AND OR
%right <str> NOT '!'
%left <str> '=' '<' '>' LE GE NE IS LIKE IN
%start any_command

%%

any_command:
command { setResult(yylex, $1) }
| command ';' { setResult(yylex, $1) }

command:
select_stmt { $$ = $1 }
| show_table_stmt { $$ = $1}

show_table_stmt:
SHOW TABLE { $$ = NewShowTable() }

select_stmt:
SELECT sel_field_list FROM table_name where_opt limit_opt
{  $$ = NewSelect($2, $4, $5, $6) }

sel_field_list:
sel_field { $$ = SelectFieldList{$1} }
| sel_field_list ',' sel_field { $$ = append($$, $3) }

sel_field:
'*' { $$ = &StarExpr{} }
| '"' col '"' { $$ = &ColExpr{$2} }
| col { $$ = &ColExpr{$1} }
| nest_col { $$ = $1 }
| func_name '(' sel_field_list ')' { $$ = &FuncExpr{Name: $1, Fields: $3} }

func_name:
IDENT { $$ = $1 }

table_name:
IDENT { $$ = $1 }

col:
IDENT { $$ = $1 }

nest_col:
IDENT '[' INTEGER ']' {
subs:= []string{$1, $3}
$$ = &NestColExpr{Subs: subs}
}
| IDENT '[' '"' IDENT '"' ']' {
subs := []string{$1, $4}
$$ = &NestColExpr{Subs: subs}
}
| IDENT '.' IDENT {
subs:= []string{$1, $3}
$$ = &NestColExpr{Subs: subs}
}
| nest_col '[' '"' IDENT '"' ']' {
col := $1.(*NestColExpr)
col.Subs = append(col.Subs, $4)
}
| nest_col '[' INTEGER ']' {
col := $1.(*NestColExpr)
col.Subs = append(col.Subs, $3)
$$ = col
}
| nest_col '.' IDENT {
col := $1.(*NestColExpr)
col.Subs = append(col.Subs, $3)
$$ = col
}

where_opt:
{ $$ = nil }
| WHERE expr { $$ = NewWhere($2) }

expr:
'(' expr ')' { $$ = $2 }
| col compare value { $$ = &ComparisonExpr{Left: $1, Operator: $2, Right: $3} }
| expr AND expr { $$ = &AndExpr{Left: $1, Right: $3} }
| expr OR  expr { $$ = &OrExpr{Left: $1, Right: $3} }
| NOT expr { $$ = &NotExpr{Expr: $2} }

limit_opt: { $$ = nil }
| LIMIT INTEGER
{
	limit, _ := strconv.Atoi($2)
	$$ = &Limit{Rowcount: limit}
}
| LIMIT INTEGER ',' INTEGER
{
	offset, _ := strconv.Atoi($2)
	limit, _ := strconv.Atoi($4)
	$$ = &Limit{Offset: offset, Rowcount: limit}
}
| LIMIT INTEGER OFFSET INTEGER
{
	limit, _ := strconv.Atoi($2)
	offset, _ := strconv.Atoi($4)
	$$ = &Limit{Offset: offset, Rowcount: limit}
}

value:
'"' '"' { $$ = value.Str{Val: ""} }
| '"' IDENT '"' { $$ = value.Str{Val: $2} }
| INTEGER
{
	v, _ := strconv.Atoi($1)
	$$ = value.Int{Val: int64(v)}
}
| FLOAT
{
	v, _ := strconv.ParseFloat($1, 64)
	$$ = value.Float{Val: v}
}
| TRUE { $$ = value.Bool{true} }
| FALSE { $$ = value.Bool{false} }
| NULL { $$ = value.Null{} }

compare:
'=' { $$ = "=" }
| '<' { $$ = "<" }
| '>' { $$ = ">" }
| LE { $$ = "<=" }
| GE { $$ = ">=" }
| NE { $$ = "!=" }
%%

func Parse(s string) (Statement, error) {
    l:= NewLexer(s)
    yyParse(l)
    return l.result, l.err
}
