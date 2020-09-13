package engine

import (
	"github.com/monsterxx03/sqlpar/value"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
)

type ParquetEngine struct {
	schema *ParquetSchema
	r      *reader.ParquetReader
	fr     source.ParquetFile
}

func NewParquetEngine(fileName string) (*ParquetEngine, error) {
	fr, err := local.NewLocalFileReader(fileName)
	if err != nil {
		return nil, err
	}
	r, err := reader.NewParquetColumnReader(fr, 2)
	if err != nil {
		return nil, err
	}
	return &ParquetEngine{schema: NewParquetSchema(r.SchemaHandler.SchemaElements), fr: fr}, nil
}

func (p *ParquetEngine) GetColumnReader() (*reader.ParquetReader, error) {
	return reader.NewParquetColumnReader(p.fr, 2)
}

func (p *ParquetEngine) FetchColumn(col string, n int64, op string, compareTo value.Value) ([]value.Value, error) {
	vals, err := p.fetch(col, n, op, compareTo)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func (p *ParquetEngine) fetch(col string, n int64, op string, compareTo value.Value) ([]value.Value, error) {
	cr, err := p.GetColumnReader()
	if err != nil {
		return nil, err
	}
	vals, _, _, err := cr.ReadColumnByPath(p.schema.GetName()+"."+col, n)
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
