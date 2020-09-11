package main

import (
	"github.com/monsterxx03/sqlpar/value"
	"github.com/xitongsys/parquet-go/reader"
)

type RecordSet struct {
	Cols   []string
	ColRes []value.Value
}

type Engine interface {
	FetchColumn(col ColExpr, compareTo value.Value, rowCount int64) (*RecordSet, error)
}

type ParquetEngine struct {
	schemaName string
	r          *reader.ParquetReader
}

func (p *ParquetEngine) FetchColumn(col string, n int64, op string, compareTo value.Value) (*RecordSet, error) {
	vals, err := p.fetch(col, n, op, compareTo)
	if err != nil {
		return nil, err
	}
	return &RecordSet{Cols: []string{col}, ColRes: vals}, nil
}

func (p *ParquetEngine) fetch(col string, n int64, op string, compareTo value.Value) ([]value.Value, error) {
	vals, _, _, err := p.r.ReadColumnByPath(p.schemaName+"."+col, n)
	if err != nil {
		return nil, err
	}
	rs := make([]value.Value, 0)
	for _, _v := range vals {
		v := value.NewFromParquetValue(_v)
		ok := true
		if compareTo != nil {
			ok, err = value.Compare(v, compareTo, op)
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
