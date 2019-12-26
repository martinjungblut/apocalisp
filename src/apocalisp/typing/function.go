package typing

import (
	"apocalisp"
)

type Function struct {
	Params      []string
	Body        Type
	Callable    (func(...Type) Type)
	Environment apocalisp.Environment
}
