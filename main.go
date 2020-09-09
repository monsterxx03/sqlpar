package main

import (
	"strings"
	"io"
	"os"
	"fmt"
	"text/tabwriter"
	"github.com/peterh/liner"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

//go:generate goyacc -o parser.go parser.y

func ExecuteSelect(pr *reader.ParquetReader, stmt *Select) error {
	allFields := false
	tgtFields := make([]string, 0, len(stmt.Fields)) 
	for _, field := range stmt.Fields {
		if _, ok := field.(*StarExpr); ok {
			allFields = true
			break
		}
		tgtFields = append(tgtFields, field.(*ColExpr).Name)
	}
	if allFields {
		// TODO query all cols
		fmt.Println("query all")
		return nil
	}
	limit := pr.GetNumRows()
	if stmt.Limit != nil {
		limit = int64(stmt.Limit.Rowcount)
	}
	result := make([][]interface{}, 0, len(tgtFields))
	for _, field := range tgtFields {
		val, _, _, err := pr.ReadColumnByPath(stmt.TableName + "." + field, limit)
		if err != nil { return  err }
		result = append(result, val)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, strings.Join(tgtFields, "\t") + "\t")
	for i :=int64(0); i<limit; i++ {
		line := make([]string, 0, len(tgtFields))
		for _, col := range result {
			line = append(line, fmt.Sprint(col[i]))
		}
		fmt.Fprintln(w, strings.Join(line, "\t") + "\t")
	}
	w.Flush()
	return nil
}


func parseSQL(sql string) Statement {
	if strings.HasSuffix(sql, ";") {
		sql = strings.TrimSuffix(sql, ";")
	}
	lex := NewLexer(sql)
	parser := yyNewParser()
	parser.Parse(lex)
	return lex.result
}

func main() {
	yyErrorVerbose = true
	if len(os.Args) == 1 {
		fmt.Println("Usage: sqlpar test.parquet")
		os.Exit(1)
	}
	ll := liner.NewLiner()
	defer ll.Close()
	ll.SetCtrlCAborts(true)

	var path = os.Args[1]
	fr, err := local.NewLocalFileReader(path)
	if err != nil {
		panic(err)
	}
	for {
		input, err := ll.Prompt(">> ")
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		ll.AppendHistory(input)
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			continue
		}
		stmt := parseSQL(input)
		if stmt == nil {
			continue
		}
		pr, err := reader.NewParquetColumnReader(fr, 2)
		if err != nil {
			panic(err)
		}
		switch v := stmt.(type) {
		case *Select:
			if err := ExecuteSelect(pr, v); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		default:
			continue
		}
	}
}
