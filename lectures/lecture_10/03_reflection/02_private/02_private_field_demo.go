package main

import (
	"demo/internal"
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	p := new(internal.PublicStruct)
	fmt.Println(p)

	setPrivateField(p, "privateField", "newData")
	fmt.Println(p)
}

func setPrivateField(
	s *internal.PublicStruct,
	fieldName string,
	value any,
) {
	objElem := reflect.ValueOf(s).Elem()
	field := objElem.FieldByName(fieldName)

	p := unsafe.Pointer(field.UnsafeAddr())
	internalField := reflect.NewAt(field.Type(), p)

	internalField.Elem().Set(reflect.ValueOf(value))
}
