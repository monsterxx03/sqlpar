package value

import (
	"fmt"
	"reflect"
)

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

var validIntType = map[reflect.Kind]bool{reflect.Uint: true, reflect.Uint8: true, reflect.Uint16: true, reflect.Uint32: true, reflect.Uint64: true,
	reflect.Int: true, reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int64: true}

func IsValidIntType(t reflect.Kind) bool {
	_, ok := validIntType[t]
	return ok
}

func IsValidNumType(t reflect.Kind) bool {
	if IsValidIntType(t) || t == reflect.Float32 || t == reflect.Float64 {
		return true
	}
	return false
}

func IsComparable(v1, v2 Value) bool {
	t1, t2 := reflect.TypeOf(v1).Kind(), reflect.TypeOf(v2).Kind()
	fmt.Println(t1, t2)
	if t1 == t2 {
		return true
	}
	if IsValidNumType(t1) && IsValidNumType(t2) {
		return true
	}
	return false
}

func Compare(v1, v2 Value, op string) (bool, error) {
	if !IsComparable(v1, v2) {
		return false, fmt.Errorf("%s and %s are not comparable", reflect.TypeOf(v1), reflect.TypeOf(v2))
	}
	switch op {
	case "==":
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