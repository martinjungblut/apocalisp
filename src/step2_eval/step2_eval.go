package main

import (
	"errors"
	"fmt"
	"github.com/peterh/liner"
	"manalisp"
	"os"
	"path/filepath"
)

func READ(sexpr string) (*manalisp.ManalispType, error) {
	return manalisp.Parse(sexpr)
}

func PRINT(malType *manalisp.ManalispType) string {
	return malType.ToString()
}

func EVAL(malType *manalisp.ManalispType, environment *manalisp.Environment) (*manalisp.ManalispType, error) {
	if !malType.IsList() {
		if t, err := evalAst(malType, environment); err == nil {
			return t, nil
		} else {
			return nil, err
		}
	} else if malType.IsEmptyList() {
		return malType, nil
	} else if malType.IsList() {
		if container, err := evalAst(malType, environment); err == nil {
			objects := container.AsList()
			function := objects[1]
			parameters := objects[2:]

			if function.IsNativeFunction() {
				result := function.CallNativeFunction(parameters...)
				return &result, nil
			} else {
				return nil, errors.New(fmt.Sprintf("Symbol is not a function: %s", function.ToString()))
			}
		} else {
			return nil, err
		}
	}

	return nil, errors.New("Unexpected behavior.")
}

func evalAst(node *manalisp.ManalispType, environment *manalisp.Environment) (*manalisp.ManalispType, error) {
	if node.IsSymbol() {
		f := environment.Find(node.AsSymbol())
		return &f, nil
	}

	if node.IsList() {
		all := manalisp.NewList()
		for _, element := range node.AsList() {
			if evaluated, err := EVAL(&element, environment); err == nil {
				all.AddToList(*evaluated)
			} else {
				return nil, err
			}
		}
		return all, nil
	}

	if node.IsVector() {
		all := manalisp.NewVector()
		for _, element := range node.AsVector() {
			if evaluated, err := EVAL(&element, environment); err == nil {
				all.AddToVector(*evaluated)
			} else {
				return nil, err
			}
		}
		return all, nil
	}

	if node.IsHashmap() {
		all := manalisp.NewHashmap()
		for _, element := range node.AsHashmap() {
			if evaluated, err := EVAL(&element, environment); err == nil {
				all.AddToHashmap(*evaluated)
			} else {
				return nil, err
			}
		}
		return all, nil
	}

	return node, nil
}

func rep(sexpr string) (string, error) {
	environment := manalisp.NewEnvironment()

	environment.DefineFunction("+", func(inputs ...manalisp.ManalispType) manalisp.ManalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			if input.IsInteger() {
				r += *input.Integer
			}
		}
		return manalisp.ManalispType{Integer: &r}
	})

	environment.DefineFunction("-", func(inputs ...manalisp.ManalispType) manalisp.ManalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r -= *input.Integer
		}
		return manalisp.ManalispType{Integer: &r}
	})

	environment.DefineFunction("/", func(inputs ...manalisp.ManalispType) manalisp.ManalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r /= *input.Integer
		}
		return manalisp.ManalispType{Integer: &r}
	})

	environment.DefineFunction("*", func(inputs ...manalisp.ManalispType) manalisp.ManalispType {
		r := *inputs[0].Integer
		for _, input := range inputs[1:] {
			r *= *input.Integer
		}
		return manalisp.ManalispType{Integer: &r}
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
