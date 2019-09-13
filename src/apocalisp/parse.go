package apocalisp

import (
	"fmt"
	"apocalisp/parser"
	"strconv"
)

func Parse(sexpr string) (*ApocalispType, error) {
	tokens := parser.Tokenize(sexpr)
	reader := parser.NewReader(tokens)
	return readForm(reader)
}

func readForm(reader *parser.Reader) (*ApocalispType, error) {
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
			seq := []ApocalispType{ApocalispType{Symbol: &symbol}, subelements[1], subelements[0]}
			return &ApocalispType{List: &seq}, nil
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

func readSequence(reader *parser.Reader) (*[]ApocalispType, error) {
	sequence := []ApocalispType{}
	for form, err := readForm(reader); form != nil || err != nil; form, err = readForm(reader) {
		if err != nil {
			return nil, err
		} else if form != nil {
			sequence = append(sequence, *form)
		}
	}
	return &sequence, nil
}

func readAtom(token *string) (*ApocalispType, error) {
	i, err := strconv.ParseInt(*token, 10, 64)
	if err == nil {
		return &ApocalispType{Integer: &i}, nil
	}
	return &ApocalispType{Symbol: token}, nil
}

func readList(reader *parser.Reader) (*ApocalispType, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &ApocalispType{List: sequence}, nil
	}
}

func readVector(reader *parser.Reader) (*ApocalispType, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &ApocalispType{Vector: sequence}, nil
	}
}

func readHashmap(reader *parser.Reader) (*ApocalispType, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &ApocalispType{Hashmap: sequence}, nil
	}
}

func readPrefixExpansion(reader *parser.Reader, symbol string) (*ApocalispType, error) {
	if form, err := readForm(reader); err != nil {
		return nil, err
	} else if form != nil {
		sequence := []ApocalispType{ApocalispType{Symbol: &symbol}, *form}
		return &ApocalispType{List: &sequence}, nil
	} else {
		return nil, nil
	}
}
