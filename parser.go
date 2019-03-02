package gojq

import (
	"strings"

	"github.com/pkg/errors"
)

// TODO make it more robust
// right now it wont work with e.g. string map keys with spaces
func Parse(cmd string) ([]command, error) {
	cs := strings.Split(cmd, " ")
	commands := make([]command, 0, len(cs))
	for _, c := range cs {
		switch {
		case strings.HasPrefix(c, "."):
			commands = append(commands, command{
				Type:     fieldT,
				Selector: c[1:],
			})
		case strings.HasPrefix(c, "[") && strings.HasSuffix(c, "]"):
			if len(c) == 2 {
				commands = append(commands, command{Type: arrayT})
			} else {
				commands = append(commands, command{
					Type:     indexT,
					Selector: c[1 : len(c)-1],
				})
			}
		case strings.HasPrefix(c, "!"):
			var found bool
			c = c[1:]
			for _, b := range builtins {
				if c == b {
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("unknown builtin")
			}
			commands = append(commands, command{
				Type:     builtinT,
				Selector: c,
			})
		default:
			return nil, errors.Errorf("unknown operation: %s", c)
		}
	}
	return commands, nil
}

func MustCompile(cmds []command, err error) []command {
	if err != nil {
		panic(err)
	}
	return cmds
}

func field(s string) command {
	return command{
		Type: fieldT,
		Selector: s,
	}
}

func builtin(s string) command {
	return command{
		Type: builtinT,
		Selector: s,
	}
}

func array() command {
	return command{Type: arrayT}
}

func index(s string) command {
	return command{
		Type: indexT,
		Selector: s,
	}
}