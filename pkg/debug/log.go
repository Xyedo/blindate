package debug

import (
	"encoding/json"
	"log"
	"reflect"
)

func Log(value any) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		Log(v.Elem())
		return
	}

	out := make(map[string]interface{})

	if v.Kind() != reflect.Struct {
		b, _ := json.MarshalIndent(value, "", "  ")
		log.Println(string(b))
		return
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get("json"); tagv != "" {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).Interface()
		}
	}

	b, _ := json.MarshalIndent(out, "", "  ")
	log.Println(string(b))
}
