package gojq

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

var nilV = reflect.ValueOf(nil)

func Exec(cmds []command, elem interface{}) (interface{}, error) {
	r, err := exec(cmds, reflect.ValueOf(elem))
	if err != nil {
		return nil, err
	}
	return r.Interface(), nil
}

func exec(cmds []command, elem reflect.Value) (reflect.Value, error) {
	if len(cmds) == 0 {
		return elem, nil
	}

	switch cmds[0].Type {
	case fieldT:
		if elem.Kind() != reflect.Struct {
			return nilV, errors.Errorf("expected struct, got: %s", elem.Kind())
		}
		return exec(cmds[1:], elem.FieldByName(cmds[0].Selector))
	case indexT:
		switch elem.Kind() {
		case reflect.Map:
			return exec(cmds[1:], elem.MapIndex(reflect.ValueOf(cmds[0].Selector)))
		case reflect.Slice, reflect.Array:
			i, err := strconv.Atoi(cmds[0].Selector)
			if err != nil {
				return nilV, errors.Errorf("cant convert index to int")
			}
			return exec(cmds[1:], elem.Index(i))
		}
		return nilV, errors.Errorf("index operation must be applied to %s, %s or %s",
			reflect.Slice, reflect.Array, reflect.Map)
	case arrayT:
		switch elem.Kind() {
		case reflect.Slice, reflect.Array:
			var r reflect.Value
			for i := 0; i < elem.Len(); i++ {
				newElem, err := exec(cmds[1:], elem.Index(i))
				if err != nil {
					return nilV, err
				} else if i == 0 {
					r = reflect.MakeSlice(reflect.SliceOf(newElem.Type()), 0, elem.Len())
				}
				r = reflect.Append(r, newElem)
			}
			return r, nil
		default:
			return nilV, errors.Errorf("array operation must be applied to %s or %s",
				reflect.Slice, reflect.Array)
		}
	case builtinT:
		switch cmds[0].Selector {
		case "len":
			return exec(cmds[1:], reflect.ValueOf(elem.Len()))
		case "keys", "values":
			if elem.Kind() != reflect.Map {
				return nilV, errors.Errorf("keys operation must be applied to %s", reflect.Map)
			}
			var r reflect.Value
			for i, k := range elem.MapKeys() {
				e := k
				if cmds[0].Selector == "values" {
					e = elem.MapIndex(k)
				}
				if i == 0 {
					r = reflect.MakeSlice(reflect.SliceOf(e.Type()), 0, elem.Len())
				}
				r = reflect.Append(r, e)
			}
			r, err := exec(cmds[1:], r)
			if err != nil {
				return nilV, err
			}
			return r, nil
		case "flatten":
			if elem.Kind() != reflect.Slice && elem.Kind() != reflect.Array {
				return nilV, errors.New("must be slice")
			}
			var r reflect.Value
			for i := 0; i < elem.Len(); i++ {
				newElem, err := exec(cmds[1:], elem.Index(i))
				if err != nil {
					return nilV, err
				} else if i == 0 {
					r = reflect.MakeSlice(newElem.Type(), 0, elem.Len())
				}
				r = reflect.AppendSlice(r, newElem)
			}
			return r, nil
		default:
			panic("unknown builtin")
		}
	default:
		panic("unknown operation type")
	}
	return nilV, nil
}
