package apocalisp

import (
	"apocalisp/core"
)

// Expose DefaultEnvironment() through the 'apocalisp' namespace.
func DefaultEnvironment() *core.Environment {
	return core.DefaultEnvironment()
}
