package apocalisp

import (
	"apocalisp/core"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/peterh/liner"
)

func withLiner(handler func(*liner.State)) {
	state := liner.NewLiner()
	defer state.Close()

	state.SetCtrlCAborts(false)
	handler(state)
}

func Repl(eval func(*core.Type, *core.Environment, bool) (*core.Type, error), parser core.Parser) {
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
	environment := DefaultEnvironment(parser, eval)

	argv := core.NewList()
	for i := range os.Args {
		if i > 1 {
			argv.Append(core.Type{String: &os.Args[i]})
		}
	}
	environment.Set("*ARGV*", *argv)
	environment.Set("*host-language*", *core.NewString("apocalisp"))

	_, _ = Rep(`(def! not (fn* (a) (if a false true)))`, environment, eval, parser)
	_, _ = Rep(`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))`, environment, eval, parser)
	_, _ = Rep(`(defmacro! cond (fn* (& xs) (if (> (count xs) 0) (list 'if (first xs) (if (> (count xs) 1) (nth xs 1) (throw "odd number of forms to cond")) (cons 'cond (rest (rest xs)))))))`, environment, eval, parser)

	if len(os.Args) >= 2 {
		_, _ = Rep(fmt.Sprintf(`(load-file "%s")`, os.Args[1]), environment, eval, parser)
	} else {
		_, _ = Rep(`(println (str "Mal [" *host-language* "]"))`, environment, eval, parser)
		for {
			if sexpr, err := line.Prompt("user> "); err == nil {
				line.AppendHistory(sexpr)

				if output, err := Rep(sexpr, environment, eval, parser); err == nil {
					if len(output) > 0 {
						fmt.Println(output)
					}
				} else {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("\nFarewell!")
				break
			}
		}
	}
}
