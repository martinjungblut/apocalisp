package apocalisp

import (
	"apocalisp/core"
	"apocalisp/parser"
)

// Expose DefaultEnvironment() through the 'apocalisp' namespace.
func DefaultEnvironment() *core.Environment {
	return core.DefaultEnvironment(parser.Parser{})
}
