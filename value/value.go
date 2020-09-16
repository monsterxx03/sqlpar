package value

import (
	"fmt"
	"strings"
)

type Row []Value

func (r *Row) String() string {
	result := make([]string, len(*r))
	for i, val := range *r {
		result[i] = val.String()
	}
	return strings.Join(result, "\t")
}

type Value interface {
	Gt(Value) bool
	Ge(Value) bool
	Lt(Value) bool
	Le(Value) bool
	Eq(Value) bool
	String() string
}

type Int struct {
	Val int64
}

func (i Int) String() string { return fmt.Sprint(i.Val) }

func (i Int) Gt(v Value) bool {
	switch v.(type) {
	case Int:
		return i.Val > v.(Int).Val
	case Float:
		return float64(i.Val) > v.(Float).Val
	default:
		return false
	}
}

func (i Int) Ge(v Value) bool {
	switch v.(type) {
	case Int:
		return i.Val >= v.(Int).Val
	case Float:
		return float64(i.Val) >= v.(Float).Val
	default:
		return false
	}
}

func (i Int) Lt(v Value) bool {
	switch v.(type) {
	case Int:
		return i.Val < v.(Int).Val
	case Float:
		return float64(i.Val) < v.(Float).Val
	default:
		return false
	}
}

func (i Int) Le(v Value) bool {
	switch v.(type) {
	case Int:
		return i.Val <= v.(Int).Val
	case Float:
		return float64(i.Val) <= v.(Float).Val
	default:
		return false
	}
}

func (i Int) Eq(v Value) bool {
	switch v.(type) {
	case Int:
		return i.Val == v.(Int).Val
	case Float:
		return float64(i.Val) == v.(Float).Val
	default:
		return false
	}
}

type Float struct {
	Val float64
}

func (i Float) String() string { return fmt.Sprint(i.Val) }

func (i Float) Gt(v Value) bool {
	switch v.(type) {
	case Float:
		return i.Val > v.(Float).Val
	case Int:
		return i.Val > float64(v.(Int).Val)
	default:
		return false
	}
}

func (i Float) Ge(v Value) bool {
	switch v.(type) {
	case Float:
		return i.Val >= v.(Float).Val
	case Int:
		return i.Val >= float64(v.(Int).Val)
	default:
		return false
	}
}

func (i Float) Lt(v Value) bool {
	switch v.(type) {
	case Float:
		return i.Val < v.(Float).Val
	case Int:
		return i.Val < float64(v.(Int).Val)
	default:
		return false
	}
}

func (i Float) Le(v Value) bool {
	switch v.(type) {
	case Float:
		return i.Val <= v.(Float).Val
	case Int:
		return i.Val <= float64(v.(Int).Val)
	default:
		return false
	}
}

func (i Float) Eq(v Value) bool {
	switch v.(type) {
	case Float:
		return i.Val == v.(Float).Val
	case Int:
		return i.Val == float64(v.(Int).Val)
	default:
		return false
	}
}

type Str struct {
	Val string
}

func (i Str) String() string { return i.Val }

func (i Str) Gt(v Value) bool { return i.Val > v.(Str).Val }

func (i Str) Ge(v Value) bool { return i.Val >= v.(Str).Val }

func (i Str) Lt(v Value) bool { return i.Val < v.(Str).Val }

func (i Str) Le(v Value) bool { return i.Val <= v.(Str).Val }

func (i Str) Eq(v Value) bool { return i.Val == v.(Str).Val }

type Bool struct {
	Val bool
}

func (i Bool) String() string { return fmt.Sprint(i.Val) }

func (i Bool) Gt(v Value) bool { return i.Val == true && v.(Bool).Val == false }

func (i Bool) Ge(v Value) bool { return i.Val == v.(Bool).Val || i.Val == true }

func (i Bool) Lt(v Value) bool { return i.Val == false && v.(Bool).Val == true }

func (i Bool) Le(v Value) bool { return i.Val == v.(Bool).Val || v.(Bool).Val == true }

func (i Bool) Eq(v Value) bool { return i.Val == v.(Bool).Val }

type Null struct{}

func (i Null) String() string  { return "null" }
func (i Null) Gt(v Value) bool { return false }
func (i Null) Ge(v Value) bool { return false }
func (i Null) Lt(v Value) bool { return false }
func (i Null) Le(v Value) bool { return false }
func (i Null) Eq(v Value) bool { return true }

type Alien struct {
	Val interface{}
}

func (i Alien) Gt(v Value) bool { return false }
func (i Alien) Ge(v Value) bool { return false }
func (i Alien) Lt(v Value) bool { return false }
func (i Alien) Le(v Value) bool { return false }
func (i Alien) Eq(v Value) bool { return false }

func (i Alien) String() string { return fmt.Sprint(i.Val) }

type Int96 struct {
}

func NewFromParquetValue(v interface{}) Value {
	switch v.(type) {
	case nil:
		return Null{}
	case int:
		return Int{Val: int64(v.(int))}
	case int8:
		return Int{Val: int64(v.(int8))}
	case int16:
		return Int{Val: int64(v.(int16))}
	case int32:
		return Int{Val: int64(v.(int32))}
	case int64:
		return Int{Val: v.(int64)}
	case uint:
		return Int{Val: int64(v.(uint))}
	case uint8:
		return Int{Val: int64(v.(uint8))}
	case uint16:
		return Int{Val: int64(v.(uint16))}
	case uint32:
		return Int{Val: int64(v.(uint32))}
	case uint64:
		return Int{Val: int64(v.(uint64))}
	case float32:
		return Float{Val: float64(v.(float32))}
	case float64:
		return Float{Val: v.(float64)}
	case bool:
		return Bool{Val: v.(bool)}
	case string:
		return Str{Val: v.(string)}
	default:
		return Alien{Val: v}
	}
}

func IsComparable(v1, v2 Value) bool {
	switch v1.(type) {
	case Int, Float:
		switch v2.(type) {
		case Int:
			return true
		case Float:
			return true
		default:
			return false
		}
	case Bool:
		if _, ok := v2.(Bool); ok {
			return true
		}
	case Str:
		if _, ok := v2.(Str); ok {
			return true
		}
	case Null:
		if _, ok := v2.(Null); ok {
			return true
		}
	default:
		return false
	}
	return false
}

func Compare(v1, v2 Value, op string) (bool, error) {
	if !IsComparable(v1, v2) {
		return false, fmt.Errorf("%t and %t are not comparable", v1, v2)
	}
	switch op {
	case "=":
		return v1.Eq(v2), nil
	case "!=", "<>":
		return !v1.Eq(v2), nil
	case ">":
		return v1.Gt(v2), nil
	case ">=":
		return v1.Ge(v2), nil
	case "<":
		return v1.Lt(v2), nil
	case "<=":
		return v1.Le(v2), nil
	}
	return false, fmt.Errorf("unknow operation %s", op)
}
