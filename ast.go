package main

type (
	Statement interface {
	}

	Select struct {
		Func   *FuncExpr
		Fields SelectFieldList
		TableName string
		Where *Where
		OrderBy OrderBy
		Limit *Limit
	}

	Where struct {
		Expr Expr
	}

	SelectExpr interface {}

	SelectField interface {
	}
	SelectFieldList []SelectField

	ColExpr struct {
		Name  string
	}
	StarExpr struct {
	}
	FuncExpr struct {
		Name string
		Fields SelectFieldList
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


func NewSelect(selExpr SelectExpr, table string, limit *Limit) *Select {
	switch v := selExpr.(type) {
	case *FuncExpr:
		return &Select{Func: v, TableName: table, Limit: limit}
	case SelectFieldList:
		return &Select{Fields: v, TableName: table, Limit: limit}
	}
	return nil
}

func NewWhere() *Where {
	return nil
}
