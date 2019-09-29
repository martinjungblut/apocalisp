package parser

import (
	"fmt"
	"testing"
)

func Test_Next_Should_Return_Next_Token(t *testing.T) {
	tokens := []string{"(", ")"}
	reader := NewReader(tokens)

	token, err := reader.Next()
	if err != nil {
		t.Error(err)
	}
	if *token != "(" {
		t.Error("Token should have been `(`.")
	}

	token, err = reader.Next()
	if err != nil {
		t.Error(err)
	}
	if *token != ")" {
		t.Error("Token should have been `)`.")
	}
}

func Test_Next_Should_Return_Nil_If_There_Are_No_More_Tokens(t *testing.T) {
	tokens := []string{}
	reader := NewReader(tokens)

	token, err := reader.Next()
	if err != nil {
		t.Error(err)
	}
	if token != nil {
		t.Error("Token should have been nil.")
	}
}

func Test_Next_Should_Return_Error_If_Syntax_Is_Invalid(t *testing.T) {
	parensMessage := "Error: unexpected ')'."
	bracesMessage := "Error: unexpected '}'."
	bracketsMessage := "Error: unexpected ']'."
	eofMessage := "Error: unexpected EOF."

	mapping := map[string]string{
		")":      parensMessage,
		"())":    parensMessage,
		"}":      bracesMessage,
		"{}}":    bracesMessage,
		"]":      bracketsMessage,
		"[]]":    bracketsMessage,
		"\"":     eofMessage,
		"\"test": eofMessage,
		"\"\"\"": eofMessage,
	}

	for input, output := range mapping {
		tokens := Tokenize(input)
		reader := NewReader(tokens)

		var err error
		for i := 0; i < len(tokens); i++ {
			_, err = reader.Next()
		}
		if err == nil {
			t.Error(fmt.Sprintf("Next() should have failed, but didn't. Input: `%s`.", input))
		} else if err.Error() != output {
			t.Error(fmt.Sprintf("Input `%s` should have yielded error `%s`.", input, output))
		}
	}
}
