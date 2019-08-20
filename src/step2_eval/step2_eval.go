package main

import (
	"errors"
	"fmt"
	"github.com/peterh/liner"
	"manalispcore"
	"os"
	"path/filepath"
)

func READ(sexpr string) (*manalispcore.MalType, error) {
	return manalispcore.ReadStr(sexpr)
}

func PRINT(malType *manalispcore.MalType) string {
	return manalispcore.PrintStr(malType)
}

func EVAL(malType *manalispcore.MalType, environment *manalispcore.Environment) (*manalispcore.MalType, error) {
	if !malType.IsList() {
		if t, err := evalAst(malType, environment); err == nil {
			return t, nil
		} else {
			return nil, err
		}
	}

	if malType.IsEmptyList() {
		return malType, nil
	} else {
		if container, err := evalAst(malType, environment); err == nil {
			objects := container.AsList()
			function := objects[1]
			parameters := objects[2:]
			result := (*function.NativeFunction)(parameters...)
			return result, nil
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Unexpected behavior.")
}

func evalAst(ast *manalispcore.MalType, environment *manalispcore.Environment) (*manalispcore.MalType, error) {
	if ast.IsSymbol() {
		if f, err := environment.Find(ast.AsSymbol()); err == nil {
			return f, nil
		} else {
			return nil, err
		}
	}

	if ast.IsList() {
		l := manalispcore.NewList()

		var e *error
		e = nil
		ast.EachInList(func(k manalispcore.MalType) {
			if t, err := EVAL(&k, environment); err == nil {
				l.AddToList(*t)
			} else {
				e = &err
			}
		})

		if e == nil {
			return l, nil
		} else {
			return nil, *e
		}
	}

	return ast, nil
}

func rep(sexpr string) (string, error) {
	environment := manalispcore.NewEnvironment()

	environment.DefineFunction("+", func(inputs ...manalispcore.MalType) *manalispcore.MalType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			if input.IsInteger() {
				r += *input.Integer
			}
		}
		return &manalispcore.MalType{Integer: &r}
	})

	environment.DefineFunction("-", func(inputs ...manalispcore.MalType) *manalispcore.MalType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r -= *input.Integer
		}
		return &manalispcore.MalType{Integer: &r}
	})

	environment.DefineFunction("/", func(inputs ...manalispcore.MalType) *manalispcore.MalType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r /= *input.Integer
		}
		return &manalispcore.MalType{Integer: &r}
	})

	environment.DefineFunction("*", func(inputs ...manalispcore.MalType) *manalispcore.MalType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r *= *input.Integer
		}
		return &manalispcore.MalType{Integer: &r}
	})

	t, err := READ(sexpr)
	if err != nil {
		return "", err
	}

	evaluated, err := EVAL(t, environment)
	if err != nil {
		return "", err
	}

	return PRINT(evaluated), nil
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

			output, err := rep(sexpr)
			if err == nil {
				fmt.Printf("%s\n", output)
			} else {
				fmt.Printf("%s\n", err.Error())
			}
		} else {
			fmt.Print("\nFarewell!\n")
			break
		}
	}
}
