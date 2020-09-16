package value

import (
	"testing"
)

var cases = []struct{
	left Value
	right Value
	op string
	comparable bool
}{
	{left: Int{1}, right: Int{2}, op: "<", comparable: true},
	{left: Int{1}, right: Float{1.0}, op: "=", comparable: true},
	{left: Int{1}, right: Bool{false}, op: "=", comparable: false},
	{left: Int{1}, right: Str{"dummy"}, op: "=", comparable: false},
	{left: Int{1}, right: Null{}, op: "=", comparable: false},
	{left: Float{1}, right: Float{1.0}, op: "=", comparable: true},
	{left: Float{1}, right: Bool{}, op: "=", comparable: false},
	{left: Float{1}, right: Str{}, op: "=", comparable: false},
	{left: Float{1}, right: Null{}, op: "=", comparable: false},
	{left: Bool{false}, right: Str{}, op: "=", comparable: false},
	{left: Bool{false}, right: Null{}, op: "=", comparable: false},
	{left: Bool{true}, right: Bool{false}, op: ">", comparable: true},
	{left: Bool{false}, right: Bool{false}, op: "=", comparable: true},
	{left: Bool{true}, right: Bool{true}, op: "=", comparable: true},
	{left: Str{""}, right: Null{}, op: "=", comparable: false},
	{left: Str{"abc"}, right: Str{"aaa"}, op: ">", comparable: true},
	{left: Null{}, right: Null{}, op: "=", comparable: true},
}

func TestValueComparable(t *testing.T) {
	for _, _case := range cases {
		comparable := IsComparable(_case.left, _case.right)
		if comparable != _case.comparable {
			t.Errorf("%s and %s are not comparable", _case.left, _case.right)
		}
		if comparable {
			ok, err := Compare(_case.left, _case.right, _case.op)
			if err != nil {
				t.Error(err)
			}
			if !ok {
				t.Errorf("failed to compare %s and %s", _case.left, _case.right)
			}
		}
	}
}
