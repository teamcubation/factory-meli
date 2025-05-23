package utils

import "reflect"

func MapSharedFields[S any, D any](src S, dest *D) {
	mapFields(reflect.ValueOf(src), reflect.ValueOf(dest).Elem())
}

func mapFields(srcVal, destVal reflect.Value) {
	srcType := srcVal.Type()
	destType := destVal.Type()

	for i := 0; i < destType.NumField(); i++ {
		destField := destType.Field(i)
		destFieldVal := destVal.Field(i)

		if !destFieldVal.CanSet() {
			continue
		}

		srcField, ok := srcType.FieldByName(destField.Name)
		if !ok {
			continue
		}
		srcFieldVal := srcVal.FieldByName(destField.Name)

		if srcField.Type == destField.Type {
			destFieldVal.Set(srcFieldVal)
		} else if destField.Type.Kind() == reflect.Struct && srcField.Type.Kind() == reflect.Struct {
			mapFields(srcFieldVal, destFieldVal)
		}
	}
}
