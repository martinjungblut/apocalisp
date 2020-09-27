package apocalisp

import (
	"apocalisp/core"
	"fmt"
	"io/ioutil"
	"strings"
)

func DefaultEnvironment(parser core.Parser, eval func(*core.Type, *core.Environment, bool) (*core.Type, error)) *core.Environment {
	environment := core.NewEnvironment(nil, []string{}, []core.Type{})

	environment.SetCallable("+", func(inputs ...core.Type) core.Type {
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
			return core.Type{Float: &rfloat}
		} else {
			var rint int64 = int64(rfloat)
			return core.Type{Integer: &rint}
		}
	})

	environment.SetCallable("-", func(inputs ...core.Type) core.Type {
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
			return core.Type{Float: &rfloat}
		} else {
			var rint int64 = int64(rfloat)
			return core.Type{Integer: &rint}
		}
	})

	environment.SetCallable("/", func(inputs ...core.Type) core.Type {
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

		return core.Type{Float: &rfloat}
	})

	environment.SetCallable("*", func(inputs ...core.Type) core.Type {
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
			return core.Type{Float: &rfloat}
		} else {
			var rint int64 = int64(rfloat)
			return core.Type{Integer: &rint}
		}
	})

	environment.SetCallable("list", func(args ...core.Type) core.Type {
		list := core.NewList()
		for _, arg := range args {
			list.Append(arg)
		}
		return *list
	})

	environment.SetCallable("list?", func(args ...core.Type) core.Type {
		return *core.NewBoolean(args[0].IsList())
	})

	environment.SetCallable("empty?", func(args ...core.Type) core.Type {
		return *core.NewBoolean(len(args[0].AsIterable()) == 0)
	})

	environment.SetCallable("count", func(args ...core.Type) core.Type {
		value := int64(len(args[0].AsIterable()))
		return core.Type{Integer: &value}
	})

	environment.SetCallable("=", func(args ...core.Type) core.Type {
		if len(args) == 2 {
			return *core.NewBoolean(args[0].Compare(args[1]))
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("<", func(args ...core.Type) core.Type {
		result := false
		if len(args) == 2 {
			result = args[0].AsNumber() < args[1].AsNumber()
		}
		return *core.NewBoolean(result)
	})

	environment.SetCallable("<=", func(args ...core.Type) core.Type {
		result := false
		if len(args) == 2 {
			result = args[0].AsNumber() <= args[1].AsNumber()
		}
		return *core.NewBoolean(result)
	})

	environment.SetCallable(">", func(args ...core.Type) core.Type {
		result := false
		if len(args) == 2 {
			result = args[0].AsNumber() > args[1].AsNumber()
		}
		return *core.NewBoolean(result)
	})

	environment.SetCallable(">=", func(args ...core.Type) core.Type {
		result := false
		if len(args) == 2 {
			result = args[0].AsNumber() >= args[1].AsNumber()
		}
		return *core.NewBoolean(result)
	})

	environment.SetCallable("pr-str", func(args ...core.Type) core.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		concatenated := strings.Join(parts, " ")
		return core.Type{String: &concatenated}
	})

	environment.SetCallable("str", func(args ...core.Type) core.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		concatenated := strings.Join(parts, "")
		return core.Type{String: &concatenated}
	})

	environment.SetCallable("prn", func(args ...core.Type) core.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(true))
		}
		fmt.Println(strings.Join(parts, " "))
		return *core.NewNil()
	})

	environment.SetCallable("println", func(args ...core.Type) core.Type {
		parts := make([]string, 0)
		for _, arg := range args {
			parts = append(parts, arg.ToString(false))
		}
		fmt.Println(strings.Join(parts, " "))
		return *core.NewNil()
	})

	environment.SetCallable("read-string", func(args ...core.Type) core.Type {
		sexpr := args[0].AsString()
		if node, err := parser.Parse(sexpr); err == nil && node != nil {
			return *node
		}
		return *core.NewNil()
	})

	environment.SetCallable("slurp", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			filepath := args[0].AsString()
			if contents, err := ioutil.ReadFile(filepath); err == nil {
				scontents := string(contents)
				return core.Type{String: &scontents}
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("atom", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewAtom(args[0])
		}
		return *core.NewNil()
	})

	environment.SetCallable("atom?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsAtom())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("deref", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return args[0].AsAtom()
		}
		return *core.NewNil()
	})

	environment.SetCallable("reset!", func(args ...core.Type) core.Type {
		if len(args) >= 2 {
			if args[0].IsAtom() {
				args[0].SetAtom(args[1])
				return args[1]
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("swap!", func(args ...core.Type) core.Type {
		if len(args) >= 2 {
			node, callable := args[0], args[1]
			fargs := append([]core.Type{node.AsAtom()}, args[2:]...)

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
		return *core.NewNil()
	})

	environment.SetCallable("cons", func(args ...core.Type) core.Type {
		list := *core.NewList()
		if len(args) >= 2 {
			list.Append(args[0])
			for _, node := range args[1].AsIterable() {
				list.Append(node)
			}
		}
		return list
	})

	environment.SetCallable("concat", func(args ...core.Type) core.Type {
		list := *core.NewList()
		for _, arg := range args {
			for _, node := range arg.AsIterable() {
				list.Append(node)
			}
		}
		return list
	})

	environment.SetCallable("vec", func(args ...core.Type) core.Type {
		vector := *core.NewVector()
		if len(args) >= 1 {
			for _, node := range args[0].AsIterable() {
				vector.Append(node)
			}
		}
		return vector
	})

	environment.SetCallable("first", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			if it := args[0].AsIterable(); len(it) >= 1 {
				return it[0]
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("rest", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			if it := args[0].AsIterable(); len(it) >= 2 {
				return *core.NewList(it[1:]...)
			}
		}
		return *core.NewList()
	})

	environment.SetCallable("nth", func(args ...core.Type) core.Type {
		if len(args) >= 2 {
			if it, nth := args[0].AsIterable(), args[1].AsInteger(); args[1].IsInteger() {
				// TODO: add test to ensure nth requires positive indexes
				if nth < 0 || nth >= int64(len(it)) {
					return *core.NewStringException(fmt.Sprintf("Invalid index '%d' for iterable of length '%d'.", nth, len(it)))
				} else {
					return it[nth]
				}
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("throw", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewException(args[0])
		}
		return *core.NewNil()
	})

	environment.SetCallable("map", func(args ...core.Type) core.Type {
		result := core.NewList()

		if len(args) >= 2 && args[1].IsIterable() {
			first, iterable := args[0], args[1].AsIterable()
			if first.IsFunction() {
				for _, e := range iterable {
					result.Append(first.CallFunction(e))
				}
			} else if first.IsCallable() {
				for _, e := range iterable {
					result.Append(first.CallCallable(e))
				}
			}
		}

		return *result
	})

	environment.SetCallable("apply", func(args ...core.Type) core.Type {
		if len(args) >= 2 {
			lastIndex := len(args) - 1
			first, middle, last := args[0], args[1:lastIndex], args[lastIndex]

			if (first.IsFunction() || first.IsCallable()) && last.IsIterable() {
				for _, e := range last.AsIterable() {
					middle = append(middle, e)
				}

				if first.IsFunction() {
					return first.CallFunction(middle...)
				} else if first.IsCallable() {
					return first.CallCallable(middle...)
				}
			}
		}

		return *core.NewList()
	})

	environment.SetCallable("hash-map", func(args ...core.Type) core.Type {
		return *core.NewHashmapFromSequence(args)
	})

	environment.SetCallable("eval", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			node := args[0]
			if r, err := eval(&node, environment, true); err == nil {
				return *r
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("nil?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsNil())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("true?", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsBoolean() {
			return *core.NewBoolean(args[0].AsBoolean())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("false?", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsBoolean() {
			return *core.NewBoolean(!args[0].AsBoolean())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("symbol?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsSymbol() && !strings.HasPrefix(args[0].AsSymbol(), ":"))
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("sequential?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsIterable())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("map?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsHashmap())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("symbol", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsString() {
			return *core.NewSymbol(args[0].AsString())
		}
		return *core.NewStringException("Provided value must be a symbol.")
	})

	environment.SetCallable("vector", func(args ...core.Type) core.Type {
		vec := *core.NewVector()
		for _, e := range args {
			vec.Append(e)
		}
		return vec
	})

	environment.SetCallable("vector?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsVector())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("keyword", func(args ...core.Type) core.Type {
		if len(args) >= 1 && (args[0].IsSymbol() || args[0].IsString()) {
			if strings.HasPrefix(args[0].AsSymbol(), ":") {
				return args[0]
			} else if args[0].IsSymbol() {
				return *core.NewSymbol(fmt.Sprintf(":%s", args[0].AsSymbol()))
			}

			if strings.HasPrefix(args[0].AsString(), ":") {
				return *core.NewSymbol(args[0].AsString())
			} else if args[0].IsString() {
				return *core.NewSymbol(fmt.Sprintf(":%s", args[0].AsString()))
			}
		}
		return *core.NewStringException("Provided value must be a symbol or string.")
	})

	environment.SetCallable("keyword?", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsSymbol() {
			return *core.NewBoolean(strings.HasPrefix(args[0].AsSymbol(), ":"))
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("keys", func(args ...core.Type) core.Type {
		keys := *core.NewList()
		if len(args) >= 1 && args[0].IsHashmap() {
			for key := range args[0].AsHashmap() {
				keys.Append(*core.NewString(key))
			}
		}
		return keys
	})

	environment.SetCallable("vals", func(args ...core.Type) core.Type {
		values := *core.NewList()
		if len(args) >= 1 && args[0].IsHashmap() {
			for _, value := range args[0].AsHashmap() {
				values.Append(value)
			}
		}
		return values
	})

	environment.SetCallable("get", func(args ...core.Type) core.Type {
		if len(args) >= 2 && args[0].IsHashmap() {
			haystack, needle := args[0].AsHashmap(), args[1]
			for key, value := range haystack {
				if needle.IsString() && key == needle.AsString() {
					return value
				}
				if needle.IsSymbol() && key == needle.AsSymbol() {
					return value
				}
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("contains?", func(args ...core.Type) core.Type {
		if len(args) >= 2 && args[0].IsHashmap() && (args[1].IsString() || args[1].IsSymbol()) {
			haystack, needle := args[0].AsHashmap(), args[1]
			for key := range haystack {
				if needle.IsString() && key == needle.AsString() {
					return *core.NewBoolean(true)
				}
				if needle.IsSymbol() && key == needle.AsSymbol() {
					return *core.NewBoolean(true)
				}
			}
		}
		return *core.NewBoolean(false)
	})

	return environment
}
