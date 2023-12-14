package debug

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func Log(value any) {
	log.Println(LogToString(value))
}

func LogToString(value any) string {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return fmt.Sprintln(nil)
		}

		return LogToString(v.Elem().Interface())

	case reflect.Struct:
		b, _ := json.MarshalIndent(parseStruct(v), "", "  ")
		return fmt.Sprintln(string(b))

	default:
		b, _ := json.MarshalIndent(value, "", "  ")
		return fmt.Sprintln(string(b))
	}

}

func parseStruct(v reflect.Value) map[string]interface{} {
	out := make(map[string]interface{})
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldType := typ.Field(i)
		if !fieldType.IsExported() {
			continue
		}
		fieldValue := v.Field(i)
		parseByKind(fieldValue, fieldType, out)

	}
	return out
}

func parseByKind(fieldValue reflect.Value, fieldType reflect.StructField, out map[string]interface{}) {
	switch fieldValue.Kind() {
	case reflect.Pointer:
		if fieldValue.IsNil() {
			setField(fieldType, out, nil)
			break
		}

		parseByKind(fieldValue.Elem(), fieldType, out)
	case reflect.Struct:
		if marshaler, ok := fieldValue.Interface().(json.Marshaler); ok {
			b, _ := marshaler.MarshalJSON()
			v, _ := strconv.Unquote(string(b))

			setField(fieldType, out, v)
			break
		}
		uVal := parseStruct(fieldValue)
		setField(fieldType, out, uVal)
	default:
		uVal := fieldValue.Interface()
		setField(fieldType, out, uVal)
	}
}

func setField(fieldType reflect.StructField, out map[string]interface{}, uVal any) {
	if tagv := fieldType.Tag.Get("json"); tagv != "" && tagv != "-" {

		out[strings.Split(tagv, ",omitempty")[0]] = uVal
	} else {
		out[fieldType.Name] = uVal
	}
}
