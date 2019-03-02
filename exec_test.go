package gojq

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExec(t *testing.T) {
	type (
		inner struct {
			Field2 uint64
			Field1 map[string]int
		}
		tt struct {
			Data []inner
		}
	)

	var (
		err error
		input, r interface{}
		cmds []command
	)

	Convey("exec", t, func() {
		cmds = MustCompile(Parse(".Field2"))
		input = inner{Field2: 10, Field1: map[string]int{"1": 1, "2": 2, "3": 3}}
		r, err = Exec(cmds, input)
		So(err, ShouldBeNil)
		So(r, ShouldResemble, uint64(10))

		cmds = MustCompile(Parse(".Field1 !len"))
		r, err = Exec(cmds, input)
		So(err, ShouldBeNil)
		So(r, ShouldEqual, 3)

		cmds = MustCompile(Parse(".Field1 !keys"))
		r, err = Exec(cmds, input)
		So(err, ShouldBeNil)

		cmds = MustCompile(Parse(".Field1 !values"))
		r, err = Exec(cmds, input)
		So(err, ShouldBeNil)
		So(r, ShouldHaveSameTypeAs, []int{})
		So(r, ShouldHaveLength, 3)
		for _, v := range r.([]int) {
			So(v, ShouldBeIn, []int{1, 2, 3})
		}

		cmds = MustCompile(Parse("[] .Field2"))
		input = []inner{{Field2: 10}, {Field2: 20}}
		r, err = Exec(cmds, input)
		So(err, ShouldBeNil)
		So(r, ShouldResemble, []uint64{10, 20})

		cmds = MustCompile(Parse("!values !flatten"))
		input = map[int][]int{1: {1, 2, 3}, 2: {4, 5, 6}}
		r, err = Exec(cmds, input)
		So(err, ShouldBeNil)
		// FIXME r should be {1,2,3,4,5,6} or {4,5,6,1,2,3}, not any sequence
		So(r, ShouldHaveSameTypeAs, []int{})
		So(r, ShouldHaveLength, 6)
		for _, v := range r.([]int) {
			So(v, ShouldBeIn, []int{1, 2, 3, 4, 5, 6})
		}
	})
}
