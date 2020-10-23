package parser

import (
	"apocalisp/core"
	"apocalisp/escaping"
	"strconv"
	"strings"
)

type Parser struct{}

func (parser Parser) Parse(sexpr string) (*core.Type, error) {
	return readForm(newReader(tokenize(sexpr)))
}

func readForm(reader *reader) (*core.Type, error) {
	token, err := reader.next()
	if err != nil {
		return core.NewStringException(err.Error()), nil
	} else if token == nil {
		return nil, nil
	}

	if *token == "^" {
		firstForm, err := readForm(reader)
		if err != nil {
			return nil, err
		}

		secondForm, err := readForm(reader)
		if err != nil {
			return nil, err
		}

		return core.NewList(*core.NewSymbol("with-meta"), *secondForm, *firstForm), nil
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
	if i, err := strconv.ParseInt(*token, 10, 64); err == nil {
		return &core.Type{Integer: &i}, nil
	}

	if f, err := strconv.ParseFloat(*token, 64); err == nil {
		return &core.Type{Float: &f}, nil
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
			return core.NewStringException(err.Error()), nil
		} else {
			return &core.Type{String: &t}, nil
		}
	}

	return &core.Type{Symbol: token}, nil
}

func readList(reader *reader) (*core.Type, error) {
	if sequence, err := readSequence(reader); err != nil {
		return nil, err
	} else {
		return &core.Type{List: sequence}, nil
	}
}

func readVector(reader *reader) (*core.Type, error) {
	if sequence, err := readSequence(reader); err != nil {
		return nil, err
	} else {
		return &core.Type{Vector: sequence}, nil
	}
}

func readHashmap(reader *reader) (*core.Type, error) {
	if sequence, err := readSequence(reader); err != nil {
		return nil, err
	} else {
		return core.NewHashmapFromSequence(*sequence), nil
	}
}

func readPrefixExpansion(reader *reader, symbol string) (*core.Type, error) {
	if form, err := readForm(reader); err != nil {
		return nil, err
	} else if form != nil {
		sequence := []core.Type{{Symbol: &symbol}, *form}
		return &core.Type{List: &sequence}, nil
	} else {
		return nil, nil
	}
}
