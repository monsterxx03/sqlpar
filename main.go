package main

import (
	"bufio"
	"fmt"
	"github.com/monsterxx03/sqlpar/engine"
	"github.com/monsterxx03/sqlpar/parser"
	"github.com/peterh/liner"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/tool/parquet-tools/schematool"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

//go:generate goyacc -o parser/parser.go parser/parser.y

var historyFile = "/home/will/.sqlpar_history"

func ExecuteSelect(pr *reader.ParquetReader, stmt *parser.Select) error {
	allFields := false
	tgtFields := make([]string, 0, len(stmt.Fields))
	for _, field := range stmt.Fields {
		switch v := field.(type) {
		case *parser.StarExpr:
			allFields = true
			break
		case *parser.ColExpr:
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

func runSelect(en *engine.ParquetEngine, stmt *parser.Select) error {
	fmt.Printf("%+v", stmt)
	return nil
}

func showTable(pr *reader.ParquetReader) {
	tree := schematool.CreateSchemaTree(pr.SchemaHandler.SchemaElements)
	fmt.Println(tree.OutputJsonSchema())
}

func RunShell() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: sqlpar test.parquet")
		os.Exit(1)
	}
	ll := liner.NewLiner()
	defer ll.Close()
	ll.SetCtrlCAborts(true)

	var path = os.Args[1]

	var file *os.File
	_, err := os.Stat(historyFile)
	if os.IsNotExist(err) {
		file, err = os.Create(historyFile)
		if err != nil {
			panic(err)
		}
	} else {
		file, err = os.OpenFile(historyFile, os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		ll.ReadHistory(file)
	}
	defer func() {
		_, err := ll.WriteHistory(file)
		if err != nil {
			panic(err)
		}
		file.Close()
	}()
	s := bufio.NewScanner(file)
	for s.Scan() {
		ll.AppendHistory(s.Text())
	}

	en, err := engine.NewParquetEngine(path)
	if err != nil {
		panic(err)
	}
	for {
		pr, err := en.GetColumnReader()
		if err != nil {
			panic(err)
		}
		input, err := ll.Prompt("sqlpar> ")
		if err == io.EOF {
			return
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
		stmt, err := parser.Parse(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result, err := en.Execute(stmt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(result)
	}
}

func test() {
}

func main() {
	RunShell()
	// test()
}
