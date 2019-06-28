package main

import (
	"fmt"
	"github.com/peterh/liner"
	"manalispcore"
	"os"
	"path/filepath"
)

func READ(sexpr string) manalispcore.MalType {
	return manalispcore.ReadStr(sexpr)
}

func PRINT(_type manalispcore.MalType) string {
	return manalispcore.PrintStr(_type)
}

func EVAL(_type manalispcore.MalType) manalispcore.MalType {
	return _type
}

func rep(sexpr string) string {
	return PRINT(EVAL(READ(sexpr)))
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
		if sexpr, err := line.Prompt("user> "); err == nil {
			line.AppendHistory(sexpr)
			fmt.Printf("%s\n", rep(sexpr))
		} else {
			fmt.Print("\nFarewell!\n")
			break
		}
	}
}
