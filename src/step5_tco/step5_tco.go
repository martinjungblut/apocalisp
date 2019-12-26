package main

import (
	"apocalisp"
	"fmt"
	"github.com/peterh/liner"
	"os"
	"path/filepath"
	"runtime/debug"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Print("Error while calling 'os.Getwd()'.")
		os.Exit(1)
	}
	historyFilePath := filepath.Join(cwd, ".apocalisp_history")

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

	environment := apocalisp.DefaultEnvironment()

	_, _ = apocalisp.Rep("(def! not (fn* (a) (if a false true)))", environment, apocalisp.Step5Eval)

	// decrease max stack size to make TCO-related tests useful
	debug.SetMaxStack(1 * 1024 * 1024)

	// repl
	fmt.Print("This is apocaLISP.\n")

	for {
		if sexpr, err := line.Prompt("user> "); err == nil {
			line.AppendHistory(sexpr)

			output, err := apocalisp.Rep(sexpr, environment, apocalisp.Step5Eval)
			if err == nil {
				if len(output) > 0 {
					fmt.Printf("%s\n", output)
				}
			} else {
				fmt.Printf("%s\n", err.Error())
			}
		} else {
			fmt.Print("\nFarewell!\n")
			break
		}
	}
}
