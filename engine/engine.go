package engine

import (
	"errors"
	_ "fmt"
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

func (p *ParquetEngine) Execute(stmt parser.Statement) (map[string][]value.Value, error) {
	switch stmt.(type) {
	case *parser.Select:
		return p.executeSelect(stmt.(*parser.Select))
	default:
		return nil, errors.New("unsupported statement")
	}
}

func (p *ParquetEngine) executeSelect(stmt *parser.Select) (map[string][]value.Value, error) {
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
		return nil, err
	}
	limit := total
	if stmt.Limit != nil {
		limit = int64(stmt.Limit.Rowcount)
	}
	result := make(map[string][]value.Value)
	if len(filterCols) > 0 {
		cols = append(cols, filterCols...)
		cols = filterDup(cols)
		rows, rowCnt, err := p.FetchRows(cols, limit)
		if err != nil {
			return nil, err
		}
		if len(rows) == 0 {
			return result, nil
		}
		row := make(map[string]value.Value)
		for i := 0; i < rowCnt; i++ {
			for n, col := range rows {
				row[n] = col[i]
			}
			if ok, err := stmt.Where.Expr.Evaluate(row); err != nil {
				return nil, err
			} else if ok {
				for n, v := range row {
					if _, ok := result[n]; ok {
						result[n] = append(result[n], v)
					} else {
						result[n] = []value.Value{v}
					}
				}
			}
		}
	} else {
		result, _, err = p.FetchRows(cols, limit)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
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

func (p *ParquetEngine) FetchRows(cols []string, limit int64) (rows map[string][]value.Value, count int, err error) {
	cr, err := p.GetColumnReader()
	if err != nil {
		return
	}
	rows = make(map[string][]value.Value)
	for _, col := range cols {
		var vals []interface{}
		// TODO loop fetch, if read count < limit
		vals, _, _, err = cr.ReadColumnByPath(p.schema.GetName()+"."+col, limit)
		if err != nil {
			return
		}
		count = len(vals)
		rows[col] = value.NewFromParquetValues(vals)
	}
	return
}

func filterDup(cols []string) []string {
	m := make(map[string]bool)
	for _, v := range cols {
		m[v] = true
	}
	result := make([]string, 0, len(m))
	for key, _ := range m {
		result = append(result, key)
	}
	return result
}
