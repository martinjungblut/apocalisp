package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Environment struct {
	outer *Environment
	table map[string]Type
}

func NewEnvironment(outer *Environment, symbols []string, nodes []Type) *Environment {
	table := make(map[string]Type)

	environment := &Environment{
		table: table,
		outer: outer,
	}

	for i := 0; i < len(symbols); i++ {
		if symbols[i] == "&" {
			rest := nodes[i:]
			if i+1 < len(symbols) {
				environment.Set(symbols[i+1], Type{List: &rest})
			} else {
				environment.Set("&", Type{List: &rest})
			}
			break
		} else {
			if i < len(nodes) {
				environment.Set(symbols[i], nodes[i])
			}
		}
	}

	return environment
}

func (env *Environment) Set(symbol string, node Type) {
	env.table[symbol] = node
}

func (env *Environment) SetCallable(symbol string, callable func(...Type) Type) {
	env.table[symbol] = Type{
		Callable: &callable,
		Symbol:   &symbol,
	}
}

func (env *Environment) Find(symbol string) *Environment {
	for key, _ := range env.table {
		if key == symbol {
			return env
		}
	}

	if env.outer != nil {
		return env.outer.Find(symbol)
	}

	return nil
}

func (env *Environment) Get(symbol string) (Type, error) {
	if e := env.Find(symbol); e != nil {
		for key, value := range e.table {
			if key == symbol {
				return value, nil
			}
		}
	}
	return Type{}, errors.New(fmt.Sprintf("Error: '%s' not found.", symbol))
}

func DefaultEnvironment(parser Parser) *Environment {
	env := NewEnvironment(nil, []string{}, []Type{})

	env.SetCallable("+", func(inputs ...Type) Type {
		floatFound := false
		var rfloat float64 = 0
		if inputs[0].IsInteger() {
			rfloat = float64(inputs[0].AsInteger())
		} else if inputs[0].IsFloat() {
			rfloat = inputs[0].AsFloat()
			floatFound = true
		}

		for _, input := range inputs[1:] {
			if input.IsInteger() {
				rfloat += float64(input.AsInteger())
			} else if input.IsFloat() {
				rfloat += input.AsFloat()
				floatFound = true
			}
		}

		if floatFound {
			return Type{Float: &rfloat}
		} else {
			var rint int64 = int64(rfloat)
			return Type{Integer: &rint}
		}
	})

	env.SetCallable("-", func(inputs ...Type) Type {
		floatFound := false
		var rfloat float64 = 0
		if inputs[0].IsInteger() {
			rfloat = float64(inputs[0].AsInteger())
		} else if inputs[0].IsFloat() {
			rfloat = inputs[0].AsFloat()
			floatFound = true
		}

		for _, input := range inputs[1:] {
			if input.IsInteger() {
				rfloat -= float64(input.AsInteger())
			} else if input.IsFloat() {
				rfloat -= input.AsFloat()
				floatFound = true
			}
		}

		if floatFound {
			return Type{Float: &rfloat}
		} else {
			var rint int64 = int64(rfloat)
			return Type{Integer: &rint}
		}
	})

	env.SetCallable("/", func(inputs ...Type) Type {
		var rfloat float64 = 0
		if inputs[0].IsInteger() {
			rfloat = float64(inputs[0].AsInteger())
		} else if inputs[0].IsFloat() {
			rfloat = inputs[0].AsFloat()
		}

		for _, input := range inputs[1:] {
			if input.IsInteger() {
				rfloat /= float64(input.AsInteger())
			} else if input.IsFloat() {
				rfloat /= input.AsFloat()
			}
		}

		return Type{Float: &rfloat}
	})

	env.SetCallable("*", func(inputs ...Type) Type {
		floatFound := false
		var rfloat float64 = 0
		if inputs[0].IsInteger() {
			rfloat = float64(inputs[0].AsInteger())
		} else if inputs[0].IsFloat() {
			rfloat = inputs[0].AsFloat()
			floatFound = true
		}

		for _, input := range inputs[1:] {
			if input.IsInteger() {
				rfloat *= float64(input.AsInteger())
			} else if input.IsFloat() {
				rfloat *= input.AsFloat()
				floatFound = true
			}
		}

		if floatFound {
			return Type{Float: &rfloat}
		} else {
			var rint int64 = int64(rfloat)
			return Type{Integer: &rint}
		}
	})

	env.SetCallable("list", func(args ...Type) Type {
		list := NewList()
		for _, arg := range args {
			list.Append(arg)
		}
		return *list
	})

	env.SetCallable("list?", func(args ...Type) Type {
		return *NewBoolean(args[0].IsList())
	})

	env.SetCallable("empty?", func(args ...Type) Type {
		var value int64 = int64(len(args[0].AsIterable()))
		return *NewBoolean(value == 0)
	})

	env.SetCallable("count", func(args ...Type) Type {
		var value int64 = int64(len(args[0].AsIterable()))
		return Type{Integer: &value}
	})

	env.SetCallable("=", func(args ...Type) Type {
		if len(args) == 2 {
			return *NewBoolean(args[0].Compare(args[1]))
		}
		return *NewBoolean(false)
	})

	env.SetCallable("<", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() < args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable("<=", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() <= args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable(">", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() > args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable(">=", func(args ...Type) Type {
		result := false
		if len(args) == 2 {
			if args[0].IsInteger() && args[1].IsInteger() {
				result = args[0].AsInteger() >= args[1].AsInteger()
			}
		}
		return *NewBoolean(result)
	})

	env.SetCallable("pr-str", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		concatenated := strings.Join(parts, " ")
		return Type{String: &concatenated}
	})

	env.SetCallable("str", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		concatenated := strings.Join(parts, "")
		return Type{String: &concatenated}
	})

	env.SetCallable("prn", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		fmt.Println(strings.Join(parts, " "))
		return *NewNil()
	})

	env.SetCallable("println", func(args ...Type) Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		fmt.Println(strings.Join(parts, " "))
		return *NewNil()
	})

	env.SetCallable("read-string", func(args ...Type) Type {
		sexpr := args[0].AsString()
		if node, err := parser.Parse(sexpr); err == nil && node != nil {
			return *node
		}
		return *NewNil()
	})

	env.SetCallable("slurp", func(args ...Type) Type {
		if len(args) >= 1 {
			filepath := args[0].AsString()
			if contents, err := ioutil.ReadFile(filepath); err == nil {
				scontents := string(contents)
				return Type{String: &scontents}
			}
		}
		return *NewNil()
	})

	env.SetCallable("atom", func(args ...Type) Type {
		if len(args) >= 1 {
			return *NewAtom(args[0])
		}
		return *NewNil()
	})

	env.SetCallable("atom?", func(args ...Type) Type {
		if len(args) >= 1 {
			return *NewBoolean(args[0].IsAtom())
		}
		return *NewBoolean(false)
	})

	env.SetCallable("deref", func(args ...Type) Type {
		if len(args) >= 1 {
			return args[0].AsAtom()
		}
		return *NewNil()
	})

	env.SetCallable("reset!", func(args ...Type) Type {
		if len(args) >= 2 {
			if args[0].IsAtom() {
				args[0].SetAtom(args[1])
				return args[1]
			}
		}
		return *NewNil()
	})

	env.SetCallable("swap!", func(args ...Type) Type {
		if len(args) >= 2 {
			node, callable := args[0], args[1]
			fargs := append([]Type{node.AsAtom()}, args[2:]...)

			if node.IsAtom() && callable.IsCallable() {
				result := callable.CallCallable(fargs...)
				node.SetAtom(result)
				return result
			}

			if node.IsAtom() && callable.IsFunction() {
				result := callable.CallFunction(fargs...)
				node.SetAtom(result)
				return result
			}
		}
		return *NewNil()
	})

	env.SetCallable("cons", func(args ...Type) Type {
		list := *NewList()
		if len(args) >= 2 {
			list.Append(args[0])
			for _, node := range args[1].AsIterable() {
				list.Append(node)
			}
		}
		return list
	})

	env.SetCallable("concat", func(args ...Type) Type {
		list := *NewList()
		for _, arg := range args {
			for _, node := range arg.AsIterable() {
				list.Append(node)
			}
		}
		return list
	})

	return env
}
