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

func readForm(r *parser.Reader) (*ManalispType, error) {
	token, err := r.Next()
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, nil
	}

	if *token == "^" {
		if list, err := readList(r); err == nil {
			symbol := "with-meta"
			subelements := *list.List
			seq := []ManalispType{ManalispType{Symbol: &symbol}, subelements[1], subelements[0]}
			return &ManalispType{List: &seq}, nil
		} else {
			fmt.Printf("%s\n", err.Error())
		}
	}

	if *token == "(" {
		return readList(r)
	} else if *token == "[" {
		return readVector(r)
	} else if *token == "{" {
		return readHashmap(r)
	} else if *token == "'" {
		return readPrefixExpansion(r, "quote")
	} else if *token == "~" {
		return readPrefixExpansion(r, "unquote")
	} else if *token == "`" {
		return readPrefixExpansion(r, "quasiquote")
	} else if *token == "@" {
		return readPrefixExpansion(r, "deref")
	} else if *token == "~@" {
		return readPrefixExpansion(r, "splice-unquote")
	} else if *token != ")" && *token != "]" && *token != "}" {
		return readAtom(token)
	}
	return nil, nil
}

func readSequence(r *parser.Reader) (*[]ManalispType, error) {
	sequence := []ManalispType{}
	for form, err := readForm(r); form != nil || err != nil; form, err = readForm(r) {
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

func readList(r *parser.Reader) (*ManalispType, error) {
	sequence, err := readSequence(r)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{List: sequence}, nil
	}
}

func readVector(r *parser.Reader) (*ManalispType, error) {
	sequence, err := readSequence(r)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{Vector: sequence}, nil
	}
}

func readHashmap(r *parser.Reader) (*ManalispType, error) {
	sequence, err := readSequence(r)
	if err != nil {
		return nil, err
	} else {
		return &ManalispType{Hashmap: sequence}, nil
	}
}

func readPrefixExpansion(r *parser.Reader, symbol string) (*ManalispType, error) {
	if form, err := readForm(r); err != nil {
		return nil, err
	} else if form != nil {
		sequence := []ManalispType{ManalispType{Symbol: &symbol}, *form}
		return &ManalispType{List: &sequence}, nil
	} else {
		return nil, nil
	}
}
