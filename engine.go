package main

import (
	"github.com/xitongsys/parquet-go/reader"
)

type Value interface{
	Eq (val Value) bool
	Gt (val Value) bool
	Ge (val Value) bool
	Lt (val Value) bool
	Le (val Value) bool
	Match (pattern string) bool
}

type ResultSet struct {
	Cols []string
	ColRes []Value
}

type Engine interface {
	FetchColumn(col ColExpr, compareTo Value, rowCount int64) (*ResultSet, error)
}


type ParquetEngine struct {
	schemaName string
	r *reader.ParquetReader
}

func (p *ParquetEngine) FetchColumn(col string, n int64, compareTo Value) (*ResultSet, error) {
	vals, err := p.fetch(col, n, compareTo)
	if err != nil {
		return nil, err
	}
	return &ResultSet{Cols: []string{col}, ColRes: vals}, nil
}

func (p *ParquetEngine) fetch(col string, n int64, compareTo Value) ([]Value, error){
	vals, _, _, err := p.r.ReadColumnByPath(p.schemaName + "." + col, n)
	if err != nil { return nil, err }
	rs := make([]Value, 0)
	for _, v := range vals {
		ok := true
		if compareTo != nil {
			ok, err = eval.Eval(v)
			if err != nil {
				return nil, err
			}
		} 
		if ok {
			rs = append(rs, v)	
		}
	}
	return rs, nil
}
