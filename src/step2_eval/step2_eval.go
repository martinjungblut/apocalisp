package main

import (
	"fmt"
	"github.com/peterh/liner"
	"manalisp"
	"os"
	"path/filepath"
)

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

			output, err := manalisp.Rep(sexpr, manalisp.Step2Eval)
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
