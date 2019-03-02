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
		s:    ".Data [1] [2]",
		cmds: []command{field("Data"), index("1"), index("2")},
	},
}

var errorCases = []string{
	"!kek",
	"!",
	"[",
	"]",
}

func TestParse(t *testing.T) {
	var (
		cmds []command
		err  error
	)

	Convey("parse", t, func() {
		for _, tc := range testCases {
			cmds, err = Parse(tc.s)
			So(err, ShouldBeNil)
			So(cmds, ShouldResemble, tc.cmds)
		}

		for _, tc := range errorCases {
			cmds, err = Parse(tc)
			So(cmds, ShouldBeNil)
			So(err, ShouldNotBeNil)
		}
	})
}
