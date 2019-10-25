package apocalisp

import (
	"apocalisp/parser"
	"apocalisp/typing"
	"fmt"
	"strconv"
	"strings"
)

func Parse(sexpr string) (*typing.Type, error) {
	tokens := parser.Tokenize(sexpr)
	reader := parser.NewReader(tokens)
	return readForm(reader)
}

func readForm(reader *parser.Reader) (*typing.Type, error) {
	token, err := reader.Next()
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, nil
	}

	if *token == "^" {
		if list, err := readList(reader); err == nil {
			symbol := "with-meta"
			subelements := *list.List
			seq := []typing.Type{typing.Type{Symbol: &symbol}, subelements[1], subelements[0]}
			return &typing.Type{List: &seq}, nil
		} else {
			fmt.Printf("%s\n", err.Error())
		}
	}

	if *token == "(" {
		return readList(reader)
	} else if *token == "[" {
		return readVector(reader)
	} else if *token == "{" {
		return readHashmap(reader)
	} else if *token == "'" {
		return readPrefixExpansion(reader, "quote")
	} else if *token == "~" {
		return readPrefixExpansion(reader, "unquote")
	} else if *token == "`" {
		return readPrefixExpansion(reader, "quasiquote")
	} else if *token == "@" {
		return readPrefixExpansion(reader, "deref")
	} else if *token == "~@" {
		return readPrefixExpansion(reader, "splice-unquote")
	} else if *token != ")" && *token != "]" && *token != "}" {
		return readAtom(token)
	}
	return nil, nil
}

func readSequence(reader *parser.Reader) (*[]typing.Type, error) {
	sequence := []typing.Type{}
	for form, err := readForm(reader); form != nil || err != nil; form, err = readForm(reader) {
		if err != nil {
			return nil, err
		} else if form != nil {
			sequence = append(sequence, *form)
		}
	}
	return &sequence, nil
}

func readAtom(token *string) (*typing.Type, error) {
	i, err := strconv.ParseInt(*token, 10, 64)
	if err == nil {
		return &typing.Type{Integer: &i}, nil
	}

	if *token == "nil" {
		return &typing.Type{Nil: true}, nil
	}

	if *token == "true" {
		v := true
		return &typing.Type{Boolean: &v}, nil
	}

	if *token == "false" {
		v := false
		return &typing.Type{Boolean: &v}, nil
	}

	if strings.HasPrefix(*token, "\"") && strings.HasSuffix(*token, "\"") {
		return &typing.Type{String: token}, nil
	}

	return &typing.Type{Symbol: token}, nil
}

func readList(reader *parser.Reader) (*typing.Type, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &typing.Type{List: sequence}, nil
	}
}

func readVector(reader *parser.Reader) (*typing.Type, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &typing.Type{Vector: sequence}, nil
	}
}

func readHashmap(reader *parser.Reader) (*typing.Type, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &typing.Type{Hashmap: sequence}, nil
	}
}

func readPrefixExpansion(reader *parser.Reader, symbol string) (*typing.Type, error) {
	if form, err := readForm(reader); err != nil {
		return nil, err
	} else if form != nil {
		sequence := []typing.Type{typing.Type{Symbol: &symbol}, *form}
		return &typing.Type{List: &sequence}, nil
	} else {
		return nil, nil
	}
}
