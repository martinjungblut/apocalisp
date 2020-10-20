package apocalisp

import (
	"apocalisp/core"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/peterh/liner"
)

func DefaultEnvironment(parser core.Parser, eval func(*core.Type, *core.Environment, bool) (*core.Type, error)) *core.Environment {
	environment := core.NewEnvironment(nil, []string{}, []core.Type{})

	environment.SetCallable("+", func(inputs ...core.Type) core.Type {
		var result float64
		for _, input := range inputs {
			result += input.AsNumber()
		}

		val := core.Type{Float: &result}
		return *val.CoerceNumber()
	})

	environment.SetCallable("-", func(inputs ...core.Type) core.Type {
		var result float64
		if len(inputs) >= 1 {
			result = inputs[0].AsNumber()
		}
		if len(inputs) >= 2 {
			for _, input := range inputs[1:] {
				result -= input.AsNumber()
			}
		}

		val := core.Type{Float: &result}
		return *val.CoerceNumber()
	})

	environment.SetCallable("/", func(inputs ...core.Type) core.Type {
		var result float64
		if len(inputs) >= 1 {
			result = inputs[0].AsNumber()
		}
		if len(inputs) >= 2 {
			for _, input := range inputs[1:] {
				result /= input.AsNumber()
			}
		}

		val := core.Type{Float: &result}
		return *val.CoerceNumber()
	})

	environment.SetCallable("*", func(inputs ...core.Type) core.Type {
		var result float64
		if len(inputs) >= 1 {
			result = inputs[0].AsNumber()
		}
		if len(inputs) >= 2 {
			for _, input := range inputs[1:] {
				result *= input.AsNumber()
			}
		}

		val := core.Type{Float: &result}
		return *val.CoerceNumber()
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
					if rval := first.CallFunction(e); rval.IsException() {
						return *rval.Exception
					} else {
						result.Append(rval)
					}
				}
			} else if first.IsCallable() {
				for _, e := range iterable {
					if rval := first.CallCallable(e); rval.IsException() {
						return *rval.Exception
					} else {
						result.Append(rval)
					}
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
			return *core.NewBoolean(args[0].IsSymbol() && !args[0].IsKeyword())
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
		if len(args) >= 1 {
			if converted, node := args[0].ToKeyword(); converted {
				return *node
			}
		}
		return *core.NewStringException("Provided value must be a symbol or string.")
	})

	environment.SetCallable("keyword?", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsSymbol() {
			return *core.NewBoolean(args[0].IsKeyword())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("keys", func(args ...core.Type) core.Type {
		keys := *core.NewList()
		if len(args) >= 1 && args[0].IsHashmap() {
			for key := range args[0].AsHashmap() {
				if key.IsSymbol {
					keys.Append(*core.NewSymbol(key.Identifier))
				} else {
					keys.Append(*core.NewString(key.Identifier))
				}
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
			if haystack, needle := args[0].AsHashmap(), args[1].AsHashmapKey(); needle != nil {
				if value, ok := haystack[*needle]; ok {
					return value
				}
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("contains?", func(args ...core.Type) core.Type {
		if len(args) >= 2 && args[0].IsHashmap() && (args[1].IsString() || args[1].IsSymbol()) {
			if haystack, needle := args[0].AsHashmap(), args[1].AsHashmapKey(); needle != nil {
				if _, ok := haystack[*needle]; ok {
					return *core.NewBoolean(ok)
				}
			}
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("assoc", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsHashmap() {
			newHashmap := core.NewHashmap().AsHashmap()
			for key, value := range args[0].AsHashmap() {
				newHashmap[key] = value
			}
			for key, value := range core.NewHashmapFromSequence(args[1:]).AsHashmap() {
				newHashmap[key] = value
			}
			return core.Type{Hashmap: &newHashmap}
		}
		return *core.NewHashmap()
	})

	environment.SetCallable("dissoc", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsHashmap() {
			blacklist := core.NewHashmap().AsHashmap()
			for _, k := range args[1:] {
				if key := k.AsHashmapKey(); key != nil {
					blacklist[*key] = *core.NewNil()
				}
			}

			newHashmap := core.NewHashmap().AsHashmap()
			for key, value := range args[0].AsHashmap() {
				if _, blacklisted := blacklist[key]; !blacklisted {
					newHashmap[key] = value
				}
			}
			return core.Type{Hashmap: &newHashmap}
		}
		return *core.NewHashmap()
	})

	environment.SetCallable("readline", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsString() {
			var input *core.Type
			withLiner(func(state *liner.State) {
				if line, err := state.Prompt(args[0].AsString()); err == nil {
					input = core.NewString(line)
				}
			})
			if input != nil {
				return *input
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("number?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsNumber())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("string?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsString())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("fn?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			if args[0].IsSymbol() {
				node := environment.Get(args[0].AsSymbol())
				return *core.NewBoolean((node.IsFunction() || node.IsCallable()) && !node.IsMacroFunction())
			} else {
				return *core.NewBoolean((args[0].IsFunction() || args[0].IsCallable()) && !args[0].IsMacroFunction())
			}
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("macro?", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			return *core.NewBoolean(args[0].IsMacroFunction())
		}
		return *core.NewBoolean(false)
	})

	environment.SetCallable("seq", func(args ...core.Type) core.Type {
		if len(args) >= 1 {
			arg := args[0]
			if arg.IsList() && !arg.IsEmptyIterable() {
				return arg
			} else if arg.IsVector() && !arg.IsEmptyIterable() {
				return *core.NewList(arg.AsIterable()...)
			} else if arg.IsString() && len(arg.AsString()) > 0 {
				result := *core.NewList()
				for _, s := range strings.Split(arg.AsString(), "") {
					result.Append(*core.NewString(s))
				}
				return result
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("conj", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].IsIterable() {
			iterable := args[0].AsIterable()
			oargs := args[1:]

			if args[0].IsList() {
				nl := *core.NewList(iterable...)
				for _, oarg := range oargs {
					nl.Prepend(oarg)
				}
				return nl
			} else if args[0].IsVector() {
				nl := *core.NewVector(iterable...)
				for _, oarg := range oargs {
					nl.Append(oarg)
				}
				return nl
			}
		}
		return *core.NewNil()
	})

	environment.SetCallable("time-ms", func(args ...core.Type) core.Type {
		time.Sleep(time.Millisecond)
		timestamp := time.Now().UnixNano() / int64(time.Millisecond)
		return core.Type{Integer: &timestamp}
	})

	environment.SetCallable("meta", func(args ...core.Type) core.Type {
		if len(args) >= 1 && args[0].Metadata != nil {
			return *args[0].Metadata
		}
		return *core.NewNil()
	})

	environment.SetCallable("with-meta", func(args ...core.Type) core.Type {
		if len(args) >= 2 {
			args[0].Metadata = &args[1]
			return args[0]
		}
		return *core.NewNil()
	})

	return environment
}
