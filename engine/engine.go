package engine

import (
	"errors"
	"github.com/monsterxx03/sqlpar/parser"
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

func (p *ParquetEngine) Execute(stmt parser.Statement) error {
	switch stmt.(type) {
	case *parser.Select:
		p.executeSelect(stmt.(*parser.Select))
	default:
		return errors.New("unsupported statement")
	}
	return nil
}

func (p *ParquetEngine) executeSelect(stmt *parser.Select) error {
	cols := []string{}
	for _, field := range stmt.Fields {
		switch field.(type) {
		case *parser.StarExpr:
			cols = p.schema.GetAllFieldNames()
			break
		case *parser.ColExpr:
			cols = append(cols, field.(*parser.ColExpr).Name)
		}
	}
	filterCols := make([]string, 0)
	if stmt.Where != nil {
		filterCols = stmt.Where.Expr.GetTargetCols()
	}
	total, err := p.GetTotalRowCount()
	if err != nil {
		return err
	}
	limit := total
	if stmt.Limit != nil {
		limit = int64(stmt.Limit.Rowcount)
	}
	if len(filterCols) > 0 {
		res, err := p.FetchRows(filterCols, limit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ParquetEngine) GetTotalRowCount() (int64, error) {
	cr, err := p.GetColumnReader()
	if err != nil {
		return 0, err
	}
	return cr.GetNumRows(), nil
}

func (p *ParquetEngine) GetColumnReader() (*reader.ParquetReader, error) {
	return reader.NewParquetColumnReader(p.fr, 2)
}

func (p *ParquetEngine) FetchRows(cols []string, limit int64) (map[string][]value.Value, error) {
	cr, err := p.GetColumnReader()
	if err != nil {
		return nil, err
	}
	result := make(map[string][]value.Value)
	for _, col := range cols {
		vals, _, _, err := cr.ReadColumnByPath(p.schema.GetName() + "."+col, limit)
		if err != nil {
			return nil, err
		}
		result[col] = value.NewFromParquetValues(vals)
	}
	return result, nil
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
