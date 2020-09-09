package main

type (
	Statement interface {
	}

	Select struct {
		Fields SelectFieldList
		TableName string
		Where *Where
		OrderBy OrderBy
		Limit *Limit
	}

	Where struct {
		Expr Expr
	}

	SelectField interface {
		
	}
	SelectFieldList []SelectField

	ColExpr struct {
		Name  string
	}
	StarExpr struct {
	}
)

type (
	Expr interface {
		
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
	Expr Expr
	Direction string
}

type Limit struct {
	Offset, Rowcount int
}


func NewSelect(fields SelectFieldList, table string, limit *Limit) *Select {
	return &Select{Fields: fields, TableName: table, Limit: limit}
}

func NewWhere() *Where {
	return nil
}
