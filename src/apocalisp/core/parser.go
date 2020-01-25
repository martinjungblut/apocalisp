package core

type Parser interface {
	Parse(sexpr string) (*Type, error)
}
