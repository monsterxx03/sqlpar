package main

import (
	"bufio"
	"fmt"
	"github.com/monsterxx03/sqlpar/engine"
	"github.com/monsterxx03/sqlpar/parser"
	"github.com/monsterxx03/sqlpar/value"
	"github.com/peterh/liner"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/tool/parquet-tools/schematool"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//go:generate goyacc -o parser/parser.go parser/parser.y

var historyFile = "~/.sqlpar_history"

func showTable(pr *reader.ParquetReader) {
	tree := schematool.CreateSchemaTree(pr.SchemaHandler.SchemaElements)
	fmt.Println(tree.OutputJsonSchema())
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		parts := strings.SplitN(path, "/", 2)
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, parts[1]), nil
	}
	return path, nil
}

func loadHistory() (*os.File, error) {
	var file *os.File
	path, err := expandPath(historyFile)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
	}
	return file, nil
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
	file, err := loadHistory()
	if err != nil {
		panic(err)
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

func printResult(rows map[string][]value.Value) {

}

func main() {
	RunShell()
}
