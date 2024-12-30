package libs

import (
	"reflect"

	goaway "github.com/TwiN/go-away"
)

// ReverseArray reverses any slice in place
func ReverseArray[T any](arr []T) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// censors all string fields in a struct or slice of structs
func Censor(data interface{}) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			censorStructFields(val.Index(i))
		}
	case reflect.Struct:
		censorStructFields(val)
	}
}

// censorStructFields handles the censoring of individual struct fields
func censorStructFields(val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String && field.CanSet() {
			str := field.String()
			censored := goaway.Censor(str)
			field.SetString(censored)
		}
	}
}
