package parser

import (
	"fmt"
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
		Evaluate(row map[string]value.Value) (bool, error)
		GetTargetCols() []string
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

	NotExpr struct {
		Expr Expr
	}
)

func (e *ComparisonExpr) Evaluate(row map[string]value.Value) (bool, error) {
	val, ok := row[e.Left]
	if !ok {
		return true, nil
	}
	if !value.IsComparable(val, e.Right) {
		return false, fmt.Errorf("%s: %t and %t are not comparable", e.Left, val, e.Right)
	}
	if ok, err := value.Compare(val, e.Right, e.Operator); err != nil {
		return false, err
	} else {
		return ok, nil
	}
}

func (e *ComparisonExpr) GetTargetCols() []string {
	return []string{e.Left}
}

func (e *AndExpr) Evaluate(row map[string]value.Value) (bool, error) {
	leftOk, err := e.Left.Evaluate(row)
	if err != nil {
		return false, err
	}
	rightOk, err := e.Right.Evaluate(row)
	if err != nil {
		return false, err
	}
	if leftOk && rightOk {
		return true, nil
	}
	return false, nil
}

func (e *AndExpr) GetTargetCols() []string {
	return filterDup(append(e.Left.GetTargetCols(), e.Right.GetTargetCols()...))
}

func (e *OrExpr) Evaluate(row map[string]value.Value) (bool, error) {
	leftOk, err := e.Left.Evaluate(row)
	if err != nil {
		return false, err
	}
	if leftOk {
		return true, nil
	}
	rightOk, err := e.Right.Evaluate(row)
	if err != nil {
		return false, err
	}
	return rightOk, nil
}

func (e *OrExpr) GetTargetCols() []string {
	return filterDup(append(e.Left.GetTargetCols(), e.Right.GetTargetCols()...))
}

func (e *NotExpr) Evaluate(row map[string]value.Value) (bool, error) {
	ok, err := e.Expr.Evaluate(row)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

func (e *NotExpr) GetTargetCols() []string {
	return filterDup(e.Expr.GetTargetCols())
}

type OrderBy []*Order

type Order struct {
	Cols      []string
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

func NewWhere(expr Expr) *Where {
	return &Where{expr}
}

func filterDup(cols []string) []string {
	m := make(map[string]bool)
	for _, v := range cols {
		m[v] = true
	}
	result := make([]string, 0, len(m))
	for key, _ := range m {
		result = append(result, key)
	}
	return result
}
