package parser

import (
	"apocalisp/core"
	"apocalisp/escaping"
	"fmt"
	"strconv"
	"strings"
)

func Parse(sexpr string) (*core.Type, error) {
	tokens := tokenize(sexpr)
	reader := newReader(tokens)
	return readForm(reader)
}

func readForm(reader *reader) (*core.Type, error) {
	token, err := reader.next()
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, nil
	}

	if *token == "^" {
		if list, err := readList(reader); err == nil {
			symbol := "with-meta"
			subelements := *list.List
			seq := []core.Type{core.Type{Symbol: &symbol}, subelements[1], subelements[0]}
			return &core.Type{List: &seq}, nil
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

func readSequence(reader *reader) (*[]core.Type, error) {
	sequence := []core.Type{}
	for form, err := readForm(reader); form != nil || err != nil; form, err = readForm(reader) {
		if err != nil {
			return nil, err
		} else if form != nil {
			sequence = append(sequence, *form)
		}
	}
	return &sequence, nil
}

func readAtom(token *string) (*core.Type, error) {
	i, err := strconv.ParseInt(*token, 10, 64)
	if err == nil {
		return &core.Type{Integer: &i}, nil
	}

	if *token == "nil" {
		return &core.Type{Nil: true}, nil
	}

	if *token == "true" || *token == "false" {
		if t, err := strconv.ParseBool(*token); err == nil {
			return &core.Type{Boolean: &t}, nil
		}
	}

	if strings.HasPrefix(*token, "\"") && strings.HasSuffix(*token, "\"") {
		if t, err := escaping.UnescapeString(strings.TrimPrefix(strings.TrimSuffix(*token, "\""), "\"")); err != nil {
			return nil, err
		} else {
			return &core.Type{String: &t}, nil
		}
	}

	return &core.Type{Symbol: token}, nil
}

func readList(reader *reader) (*core.Type, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &core.Type{List: sequence}, nil
	}
}

func readVector(reader *reader) (*core.Type, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &core.Type{Vector: sequence}, nil
	}
}

func readHashmap(reader *reader) (*core.Type, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &core.Type{Hashmap: sequence}, nil
	}
}

func readPrefixExpansion(reader *reader, symbol string) (*core.Type, error) {
	if form, err := readForm(reader); err != nil {
		return nil, err
	} else if form != nil {
		sequence := []core.Type{core.Type{Symbol: &symbol}, *form}
		return &core.Type{List: &sequence}, nil
	} else {
		return nil, nil
	}
}
