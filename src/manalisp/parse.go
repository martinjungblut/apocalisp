package manalisp

import (
	"fmt"
	"manalisp/parser"
	"strconv"
)

func Parse(sexpr string) (*ManalispType, error) {
	tokens := parser.Tokenize(sexpr)
	reader := parser.NewReader(tokens)
	return readForm(reader)
}

func readForm(reader *parser.Reader) (*ManalispType, error) {
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
			seq := []ManalispType{ManalispType{Symbol: &symbol}, subelements[1], subelements[0]}
			return &ManalispType{List: &seq}, nil
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

func readSequence(reader *parser.Reader) (*[]ManalispType, error) {
	sequence := []ManalispType{}
	for form, err := readForm(reader); form != nil || err != nil; form, err = readForm(reader) {
		if err != nil {
			return nil, err
		} else if form != nil {
			sequence = append(sequence, *form)
		}
	}
	return &sequence, nil
}

func readAtom(token *string) (*ManalispType, error) {
	i, err := strconv.ParseInt(*token, 10, 64)
	if err == nil {
		return &ManalispType{Integer: &i}, nil
	}
	return &ManalispType{Symbol: token}, nil
}

func readList(reader *parser.Reader) (*ManalispType, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{List: sequence}, nil
	}
}

func readVector(reader *parser.Reader) (*ManalispType, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{Vector: sequence}, nil
	}
}

func readHashmap(reader *parser.Reader) (*ManalispType, error) {
	sequence, err := readSequence(reader)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{Hashmap: sequence}, nil
	}
}

func readPrefixExpansion(reader *parser.Reader, symbol string) (*ManalispType, error) {
	if form, err := readForm(reader); err != nil {
		return nil, err
	} else if form != nil {
		sequence := []ManalispType{ManalispType{Symbol: &symbol}, *form}
		return &ManalispType{List: &sequence}, nil
	} else {
		return nil, nil
	}
}
