package parser

import (
	"github.com/monsterxx03/sqlpar/value"
)

type (
	Statement interface {
	}

	Desc struct {
		TableName string
	}

	Select struct {
		Fields    SelectFieldList
		TableName string
		Where     *Where
		OrderBy   OrderBy
		GroupBy   GroupBy
		Limit     *Limit
	}

	Where struct {
		Expr Expr
	}

	SelectField interface {
	}
	SelectFieldList []SelectField

	ColExpr struct {
		Name string
	}
	StarExpr struct {
	}
	FuncExpr struct {
		Name   string
		Fields SelectFieldList
	}
)

type (
	Expr interface {
	}
	ComparisonExpr struct {
		Left     string
		Operator string
		Right    value.Value
	}
	AndExpr struct {
		Left, Right Expr
	}

	OrExpr struct {
		Left, Right Expr
	}

	XorExpr struct {
		Left, Right Expr
	}

	NotExpr struct {
		Expr Expr
	}
)

type OrderBy []*Order

type Order struct {
	Expr      Expr
	Direction string
}

type GroupBy []string

type Limit struct {
	Offset, Rowcount int
}

func NewSelect(fields SelectFieldList, table string, where *Where, limit *Limit) *Select {
	return &Select{Fields: fields, TableName: table,
		Where: where, Limit: limit}
}

func NewWhere() *Where {
	return nil
}
