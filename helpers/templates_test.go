package helpers

import (
	"fmt"
	"reflect"
	"testing"
)

type StudentNested struct {
	Common string
}

type Student struct {
	StudentNested
	First  string
	Last   string
	City   string
	Mobile int64
}

func TestThing(t *testing.T) {

	s := Student{
		First:  "Chetan",
		Last:   "Kumar",
		City:   "Bangalore",
		Mobile: 7777777777,
	}
	v := reflect.ValueOf(s)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
}
