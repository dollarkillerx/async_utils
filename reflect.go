package async_utils

import (
	"fmt"
	"reflect"
)

type StructItem struct {
	Data interface{}
	Type string
}

func GetStructVal(st interface{}, field ...string) (result map[string]StructItem, err error) {
	result = map[string]StructItem{}

	t := reflect.TypeOf(st)
	exs := sliceToSet(field)
	if t.Kind().String() != "struct" {
		return nil, fmt.Errorf("Not Struct")
	}

	v := reflect.ValueOf(st)
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i)
		value := v.Field(i).Interface()

		_, ex := exs[key.Name]
		if ex {
			result[key.Name] = StructItem{
				Data: value,
				Type: key.Type.String(),
			}
		}
	}

	return result, nil
}

func sliceToSet(field []string) map[string]bool {
	output := map[string]bool{}
	for _, v := range field {
		output[v] = true
	}

	return output
}
