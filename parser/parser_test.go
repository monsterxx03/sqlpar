package parser

import (
	"github.com/monsterxx03/sqlpar/value"
	"reflect"
	"testing"
)

var cases = []struct {
	sql    string
	parsed *Select
	err    string
}{
	{
		sql: "select name,age from user",
		parsed: &Select{
			Fields:    SelectFieldList{&ColExpr{Name: "name"}, &ColExpr{Name: "age"}},
			TableName: "user",
		},
		err: "",
	},
	{
		sql: "select * from user",
		parsed: &Select{
			Fields:    SelectFieldList{&StarExpr{}},
			TableName: "user",
		},
	},
	{
		sql: "select * from user where age > 10",
		parsed: &Select{
			Fields:    SelectFieldList{&StarExpr{}},
			TableName: "user",
			Where:     &Where{&ComparisonExpr{Left: "age", Operator: ">", Right: value.Int{10}}},
		},
	},
	{
		sql: "select * from user where age = 10 and name=\"haha\"",
		parsed: &Select{
			Fields:    SelectFieldList{&StarExpr{}},
			TableName: "user",
			Where: &Where{
				Expr: &AndExpr{
					Left:  &ComparisonExpr{Left: "age", Operator: "=", Right: value.Int{10}},
					Right: &ComparisonExpr{Left: "name", Operator: "=", Right: value.Str{"haha"}},
				}},
		},
	},
	{
		sql: "select * from user where age >= 10 and name!=\"haha\"",
		parsed: &Select{
			Fields:    SelectFieldList{&StarExpr{}},
			TableName: "user",
			Where: &Where{
				Expr: &AndExpr{
					Left:  &ComparisonExpr{Left: "age", Operator: ">=", Right: value.Int{10}},
					Right: &ComparisonExpr{Left: "name", Operator: "!=", Right: value.Str{"haha"}},
				}},
		},
	},
	{
		sql: "select * from user where age >= 10 or name=\"haha\" and (age < 100)",
		parsed: &Select{
			Fields: SelectFieldList{&StarExpr{}},
			TableName: "user",
			Where: &Where{
				Expr: &AndExpr{
					Left: &OrExpr{
						Left: &ComparisonExpr{Left: "age", Operator: ">=", Right: value.Int{10}},
						Right: &ComparisonExpr{Left: "name", Operator: "=", Right: value.Str{"haha"}},
					},
					Right: &ComparisonExpr{Left: "age", Operator: "<", Right: value.Int{100}},
				},
			},
		},
	},
}

func TestParser(t *testing.T) {
	for _, _case := range cases {
		stmt, err := Parse(_case.sql)
		if _case.err == "" && err != nil {
			t.Error(err)
		}
		if _case.err != "" && _case.err != err.Error() {
			t.Errorf("error not match, expected: %s, get: %s", _case.err, err.Error())
		}
		if !reflect.DeepEqual(stmt, _case.parsed) {
			t.Errorf("parsed not match, expected: %+v, get: %+v", _case.parsed, stmt)
		}
	}
}
