package util

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldDesc struct {
	Name      string
	Type      string
	Optional  bool
	Kind      reflect.Kind
	Value     interface{}
	SliceElem *FieldDesc
	Struct    []*FieldDesc
}

func isUnexported(sf reflect.StructField) bool {
	return strings.ToLower(sf.Name[0:1]) == sf.Name[0:1]
}

func dumpType(t reflect.Type) (res []*FieldDesc) {
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		if isUnexported(sf) {
			continue
		}

		if sf.Type.Kind() == reflect.Map {
			continue
		}

		desc := &FieldDesc{
			Name:     sf.Name,
			Type:     sf.Type.String(),
			Kind:     sf.Type.Kind(),
			Optional: sf.Type.Kind() == reflect.Ptr,
		}

		switch sf.Type.Kind() {
		case reflect.Array, reflect.Ptr, reflect.Slice:
			desc.SliceElem = dumpType(sf.Type.Elem())[0]
		case reflect.Struct:
			desc.Struct = dumpType(sf.Type)
		}

		fmt.Printf("%+v\n", desc)
		res = append(res, desc)
	}

	return res
}

func DumpStruct(x interface{}) (res []*FieldDesc) {
	t := reflect.TypeOf(x)
	res = dumpType(t)

	// v := reflect.ValueOf(x)

	// for i := 0; i < v.NumField(); i++ {
	// 	ft := t.Field(i)
	// 	fv := v.Field(i)

	// 	if isUnexported(ft.Name) {
	// 		continue
	// 	}

	// 	if fv.Kind() == reflect.Map {
	// 		continue
	// 	}

	// 	desc := &FieldDesc{
	// 		Name:     ft.Name,
	// 		Type:     ft.Type.String(),
	// 		Kind:     fv.Kind(),
	// 		Optional: ft.Type.Kind() == reflect.Ptr,
	// 	}

	// 	// if desc.Optional {
	// 	// 	desc.Type = ft.Type.Elem().String()
	// 	// }

	// 	if fv.Kind() != reflect.Struct {
	// 		desc.Value = fv.Interface()
	// 	}

	// 	// fmt.Printf("%s: %s (%v)\n", ft.Name, fv.Kind(), fv.Interface())
	// 	// fmt.Printf("  tags: %s\n", ft.Tag.Get("mapstructure"))

	// 	switch fv.Kind() {
	// 	case reflect.Slice:
	// 		for i := 0; i < fv.Len(); i++ {
	// 			desc.Elem = append(desc.Elem, DumpStruct(fv.Index(i))...)
	// 		}
	// 	case reflect.Struct:
	// 		desc.Struct = DumpStruct(fv.Interface())
	// 	}

	// 	fmt.Printf("%+v\n", desc)
	// 	res = append(res, desc)
	// }

	return res
}
