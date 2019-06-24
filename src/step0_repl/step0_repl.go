package main

import (
	"fmt"
	"github.com/peterh/liner"
	"os"
	"path/filepath"
)

func READ(sexp string) string {
	return sexp
}

func PRINT(sexp string) string {
	return sexp
}

func EVAL(sexp string) string {
	return sexp
}

func rep(sexp string) string {
	return PRINT(EVAL(READ(sexp)))
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Print("Error while calling 'os.Getwd()'.")
		os.Exit(1)
	}
	historyFilePath := filepath.Join(cwd, ".manalisp_history")

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	// read/write history
	if f, err := os.Open(historyFilePath); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	defer func() {
		if f, err := os.Create(historyFilePath); err == nil {
			line.WriteHistory(f)
			f.Close()
		}
	}()

	// repl
	fmt.Print("This is manaLISP.\n")
	for {
		if sexp, err := line.Prompt("user> "); err == nil {
			line.AppendHistory(sexp)
			fmt.Printf("%s\n", rep(sexp))
		} else {
			fmt.Print("\nFarewell!\n")
			break
		}
	}
}
