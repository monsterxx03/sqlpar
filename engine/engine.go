package engine

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/monsterxx03/sqlpar/parser"
	"github.com/monsterxx03/sqlpar/value"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
)

type RecordSet struct {
	Cols []string
	Rows []value.Row
}

func (rs *RecordSet) Fit(cols []string) {
	m := make(map[string]int)
	for i, col := range rs.Cols {
		m[col] = i
	}
	rRows := make([]value.Row, len(rs.Rows))
	for i := range rRows {
		rRows[i] = make([]value.Value, len(cols))
	}
	for k, col := range cols {
		j := m[col]
		for i, row := range rs.Rows {
			rRows[i][k] = row[j]
		}
	}
	rs.Cols = cols
	rs.Rows = rRows
}

func (rs *RecordSet) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, strings.Join(rs.Cols, "\t"))
	for _, row := range rs.Rows {
		fmt.Fprintln(w, row.String())
	}
	w.Flush()
	return buf.String()
}

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

func (p *ParquetEngine) Execute(stmt parser.Statement) (*RecordSet, error) {
	switch stmt.(type) {
	case *parser.Select:
		return p.executeSelect(stmt.(*parser.Select))
	default:
		return nil, errors.New("unsupported statement")
	}
}

func (p *ParquetEngine) executeSelect(stmt *parser.Select) (*RecordSet, error) {
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
	limit := int(total)
	if stmt.Limit != nil {
		limit = stmt.Limit.Rowcount
	}
	cr, err := p.GetColumnReader()
	if err != nil {
		return nil, err
	}
	var result *RecordSet
	if len(filterCols) > 0 {
		_cols := append(cols, filterCols...)
		_cols = filterDup(_cols)
		result = &RecordSet{Cols: _cols, Rows: make([]value.Row, 0)}
		for {
			res, err := p.FetchRows(cr, _cols, limit-len(result.Rows))
			if err != nil {
				return nil, err
			}
			if len(res.Rows) == 0 {
				break
			}
			res, err = filter(stmt.Where.Expr, res)
			if err != nil {
				return nil, err
			}
			result.Rows = append(result.Rows, res.Rows...)
			if len(result.Rows) >= int(limit) {
				result.Rows = result.Rows[:limit]
				break
			}
		}
		result.Fit(cols)
	} else {
		result, err = p.FetchRows(cr, cols, limit)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func filter(expr parser.Expr, result *RecordSet) (*RecordSet, error) {
	rows := make([]value.Row, 0)
	for _, row := range result.Rows {
		if ok, err := expr.Evaluate(result.Cols, row); err != nil {
			return nil, err
		} else if ok {
			rows = append(rows, row)
		}
	}
	result.Rows = rows
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

func (p *ParquetEngine) FetchRows(cr *reader.ParquetReader, cols []string, limit int) (result *RecordSet, err error) {
	result = &RecordSet{Cols: cols, Rows: make([]value.Row, 0)}
	cidx := -1
	for _, col := range cols {
		cidx++
		var vals []interface{}
		vals, _, _, err = cr.ReadColumnByPath(p.schema.GetName()+"."+col, int64(limit))
		if err != nil {
			return
		}
		if len(vals) == 0 {
			return result, nil
		}
		// init Rows
		if len(result.Rows) == 0 {
			result.Rows = make([]value.Row, len(vals))
			for i := 0; i < len(vals); i++ {
				result.Rows[i] = make(value.Row, len(cols))
			}
		}
		for ridx, val := range vals {
			result.Rows[ridx][cidx] = value.NewFromParquetValue(val)
		}
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
