package httphelper

import (
	"reflect"
)

func UnpackToResp(f any, resp map[string]any) {
	userReflect := reflect.ValueOf(f)
	userReflectType := reflect.TypeOf(f)

	for i := 0; i < userReflect.NumField(); i++ {
		value := userReflect.Field(i)
		key := userReflectType.Field(i).Tag.Get("json")

		resp[key] = value.Interface()
	}
}
