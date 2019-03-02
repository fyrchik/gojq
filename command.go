package gojq

type (
	cmdType uint8

	command struct {
		Type cmdType
		Selector string
	}
)

const (
	_ cmdType = iota
	fieldT
	indexT
	arrayT
	builtinT
)

var builtins = []string{
	"len", "flatten", "keys", "values",
}