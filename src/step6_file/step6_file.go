package main

import (
	"apocalisp"
	"apocalisp/core"
	"fmt"
	"github.com/peterh/liner"
	"os"
	"path/filepath"
	"runtime/debug"
)

func main() {
	// reference to evaluation function
	EVAL := apocalisp.Step6Eval

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

	// environment
	environment := apocalisp.DefaultEnvironment()
	environment.SetCallable("eval", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			node := args[0]
			if r, err := EVAL(&node, environment); err == nil {
				return *r
			}
		}
		return *core.NewNil()
	})

	_, _ = apocalisp.Rep(`(def! not (fn* (a) (if a false true)))`, environment, EVAL)
	_, _ = apocalisp.Rep(`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))`, environment, EVAL)

	// decrease max stack size to make TCO-related tests useful
	debug.SetMaxStack(1 * 1024 * 1024)

	// repl
	fmt.Println("This is apocaLISP.")
	for {
		if sexpr, err := line.Prompt("user> "); err == nil {
			line.AppendHistory(sexpr)

			output, err := apocalisp.Rep(sexpr, environment, EVAL)
			if err == nil {
				if len(output) > 0 {
					fmt.Printf("%s\n", output)
				}
			} else {
				fmt.Printf("%s\n", err.Error())
			}
		} else {
			fmt.Println("\nFarewell!")
			break
		}
	}
}
