package main

import (
	"fmt"
	"github.com/peterh/liner"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/tool/parquet-tools/schematool"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

//go:generate goyacc -o parser.go parser.y

func ExecuteSelect(pr *reader.ParquetReader, stmt *Select) error {
	allFields := false
	tgtFields := make([]string, 0, len(stmt.Fields))
	for _, field := range stmt.Fields {
		switch v := field.(type) {
		case *StarExpr:
			allFields = true
			break
		case *ColExpr:
			tgtFields = append(tgtFields, v.Name)
		default:
			return fmt.Errorf("don't support %+v", v)
		}
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
		val, _, _, err := pr.ReadColumnByPath(stmt.TableName+"."+field, limit)
		if err != nil {
			return err
		}
		result = append(result, val)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, strings.Join(tgtFields, "\t")+"\t")
	for i := int64(0); i < limit; i++ {
		line := make([]string, 0, len(tgtFields))
		for _, col := range result {
			line = append(line, fmt.Sprint(col[i]))
		}
		fmt.Fprintln(w, strings.Join(line, "\t")+"\t")
	}
	w.Flush()
	return nil
}

func showTable(pr *reader.ParquetReader) {
	tree := schematool.CreateSchemaTree(pr.SchemaHandler.SchemaElements)
	fmt.Println(tree.OutputJsonSchema())
}

func parseSQL(sql string) Statement {
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
		pr, err := reader.NewParquetColumnReader(fr, 2)
		if err != nil {
			panic(err)
		}
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
		if strings.HasSuffix(input, ";") {
			input = strings.TrimSuffix(input, ";")
		}
		if strings.ToLower(input) == "show table" {
			showTable(pr)
			continue
		}
		stmt := parseSQL(input)
		if stmt == nil {
			continue
		}
		switch v := stmt.(type) {
		case *Select:
			if err := ExecuteSelect(pr, v); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		case *Desc:
			fmt.Println(v)
		default:
			continue
		}
	}
}
