%{
package main

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
%token <str> SELECT FROM WHERE ORDER_BY LIMIT OFFSET 
%token <str> IDENT INTEGRAL

%type <expr> expr
%type <str> table_name col func_name compare 
%type <value> value
%type <stmt> command
%type <sel> select_stmt
%type <sel_field> sel_field
%type <sel_field_list> sel_field_list
%type <where> where_opt
%type <limit> limit_opt

%left <str> AND OR
%right <str> NOT '!'
%left <str> '=' '<' '>' LE GE NE IS LIKE IN
%start any_command

%%

any_command:
    command
    {
        setResult(yylex, $1)
    }

command:
    select_stmt
    {
        $$ = $1
    }

select_stmt:
  SELECT sel_field_list FROM table_name where_opt limit_opt
  {
    $$ = NewSelect($2, $4, $5, $6)
  }



sel_field_list:
    sel_field
    {
        $$ = SelectFieldList{$1}
    }
| sel_field_list ',' sel_field
    {
        $$ = append($$, $3)
    }


sel_field:
'*'
    {
        $$ = &StarExpr{}
    }
| col
    {
        $$ = &ColExpr{$1}
    }
| func_name '(' sel_field_list ')'
    {
        $$ = &FuncExpr{Name: $1, Fields: $3}
    }

func_name:
    IDENT
    {
        $$ = $1
    }

table_name:
    IDENT
    { 
      $$ = $1
    }


col:
  IDENT
  {
      $$ = $1
  }

where_opt:
    {
        $$ = nil
    }
| WHERE expr 
    {
    $$ = &Where{$2}
    } 

expr:
  col compare value 
  {
     $$ = &ComparisonExpr{Left: $1, Operator: $2, Right: $3}
  } 
| expr AND expr
  {
    $$ = &AndExpr{Left: $1, Right: $3}
  }
| expr OR  expr
  {
    $$ = &OrExpr{Left: $1, Right: $3}
  }
| NOT expr
  {
    $$ = &NotExpr{Expr: $2}
  }

limit_opt:
  {
    $$ = nil
  }
| LIMIT INTEGRAL
    {
        limit, _ := strconv.Atoi($2)
        $$ = &Limit{Rowcount: limit}
    }
| LIMIT INTEGRAL ',' INTEGRAL
    {
        offset, _ := strconv.Atoi($2)
        limit, _ := strconv.Atoi($4)
        $$ = &Limit{Offset: offset, Rowcount: limit}
    }
| LIMIT INTEGRAL OFFSET INTEGRAL
  {
        limit, _ := strconv.Atoi($2)
        offset, _ := strconv.Atoi($4)
        $$ = &Limit{Offset: offset, Rowcount: limit}
  }

value:
    IDENT
    {
      $$ = value.Str{Val: $1}
    }
|
    INTEGRAL
    {
      v, _ := strconv.Atoi($1)
      $$ = value.Int{Val: int64(v)}
    }

compare:
  '='
  {
    $$ = "="
  }
| '<'
  {
    $$ = "<"
  }
| '>'
  {
    $$ = ">"
  }
| LE
  {
    $$ = "<="
  }
| GE 
    {
    $$ = ">="
    }
| NE
   {
    $$ = "!="
    }
