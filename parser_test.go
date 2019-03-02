package gojq

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testCases = []struct {
	s    string
	cmds []command
}{
	{
		s:    ".Data !flatten",
		cmds: []command{field("Data"), builtin("flatten")},
	},
	{
		s:    ".Data [1] [2] []",
		cmds: []command{field("Data"), index("1"), index("2"), array()},
	},
	{
		s:    `[mapkey]`,
		cmds: []command{index("mapkey")},
	},
}

var errorCases = []string{
	"!kek",
	"!",
	".f !",
	"[2] [",
	"]",
	"1",
}

func TestParse(t *testing.T) {
	var cmds []command

	Convey("parse", t, func() {
		for _, tc := range testCases {
			cmds = MustCompile(Parse(tc.s))
			So(cmds, ShouldResemble, tc.cmds)
		}

		for _, tc := range errorCases {
			So(func() { MustCompile(Parse(tc)) }, ShouldPanic)
		}
	})
}
