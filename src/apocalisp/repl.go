package apocalisp

import (
	"apocalisp/core"
	"fmt"
	"github.com/peterh/liner"
	"os"
	"path/filepath"
	"runtime/debug"
)

func Repl(eval func(*core.Type, *core.Environment) (*core.Type, error)) {
	// decrease max stack size to make TCO-related tests useful
	debug.SetMaxStack(1 * 1024 * 1024)

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
	environment := DefaultEnvironment()
	environment.SetCallable("eval", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			node := args[0]
			if r, err := eval(&node, environment); err == nil {
				return *r
			}
		}
		return *core.NewNil()
	})

	argv := core.NewList()
	for i := range os.Args {
		if i > 1 {
			argv.Append(core.Type{String: &os.Args[i]})
		}
	}
	environment.Set("*ARGV*", *argv)

	_, _ = Rep(`(def! not (fn* (a) (if a false true)))`, environment, eval)
	_, _ = Rep(`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))`, environment, eval)

	if len(os.Args) >= 2 {
		_, _ = Rep(fmt.Sprintf(`(load-file "%s")`, os.Args[1]), environment, eval)
	} else {
		// repl
		fmt.Println("This is apocaLISP.")
		for {
			if sexpr, err := line.Prompt("user> "); err == nil {
				line.AppendHistory(sexpr)

				output, err := Rep(sexpr, environment, eval)
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
}
