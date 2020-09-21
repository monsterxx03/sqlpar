package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/monsterxx03/sqlpar/engine"
	"github.com/monsterxx03/sqlpar/parser"
	"github.com/peterh/liner"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//go:generate goyacc -o parser/parser.go parser/parser.y

var historyFile = "~/.sqlpar_history"

var (
	sqlQuery    = flag.String("sql", "", "run sql directly")
	parquetFile = flag.String("file", "", "parquet file to query")
)

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

func runShell() {
	ll := liner.NewLiner()
	defer ll.Close()
	ll.SetCtrlCAborts(true)

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

	en, err := engine.NewParquetEngine(*parquetFile)
	if err != nil {
		panic(err)
	}
	for {
		input, err := ll.Prompt("sqlpar> ")
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		ll.AppendHistory(input)
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
		if result != nil {
			fmt.Println(result)
		}
	}
}

func runSQL(sql string) (*engine.RecordSet, error) {
	en, err := engine.NewParquetEngine(*parquetFile)
	if err != nil {
		return nil, err
	}
	stmt, err := parser.Parse(sql)
	if err != nil {
		return nil, err
	}
	result, err := en.Execute(stmt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	flag.Parse()
	if *parquetFile == "" {
		fmt.Fprintf(os.Stderr, "missing -file")
		os.Exit(1)
	}
	if *sqlQuery == "" {
		runShell()
	} else {
		result, err := runSQL(*sqlQuery)
		if err != nil {
			panic(err)
		}
		if result != nil {
			fmt.Println(result)
		}
	}
}
