package utils

import (
	"log"
	"reflect"
)

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

// ListEncryptField 提取所有的需要加解密字段
func ListEncryptField(ptr interface{}) (fieldValues []reflect.Value) {
	return ListFieldWithTag(ptr, "encrypted")
}

func ListTagValue(ptr interface{}, tagName string) []string {
	result := make([]string, 0)
	v := reflect.ValueOf(ptr)
	t := reflect.TypeOf(ptr)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {
			newTagValues := ListTagValue(fieldValue.Interface(), tagName)
			result = append(result, newTagValues...)
			continue
		}
		tag := fieldType.Tag

		tagValue, ok := tag.Lookup(tagName)
		if !ok {
			continue
		}

		result = append(result, tagValue)
	}
	return result
}

func ListTagValueByCustom(ptr interface{}, tagName string, customFieldNames []string) []string {
	result := make([]string, 0)
	v := reflect.ValueOf(ptr)
	t := reflect.TypeOf(ptr)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	customFieldNameMap := make(map[string]bool)
	for _, customFieldName := range customFieldNames {
		customFieldNameMap[customFieldName] = true
	}

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {
			newTagValues := ListTagValue(fieldValue.Interface(), tagName)
			result = append(result, newTagValues...)
			continue
		}
		tag := fieldType.Tag

		tagValue, ok := tag.Lookup(tagName)
		if !ok {
			continue
		}
		customValue, ok := tag.Lookup("custom")
		if !ok {
			continue
		}

		if !customFieldNameMap[customValue] {
			continue
		}

		result = append(result, tagValue)
	}
	return result
}

func ListFieldWithTag(ptr interface{}, tagName string) (fieldValues []reflect.Value) {
	v := reflect.ValueOf(ptr)
	t := reflect.TypeOf(ptr)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {
			newFieldValues := ListFieldWithTag(fieldValue.Interface(), tagName)
			fieldValues = append(fieldValues, newFieldValues...)
			continue
		}
		tag := fieldType.Tag

		_, ok := tag.Lookup(tagName)
		if !ok {
			continue
		}

		fieldValues = append(fieldValues, fieldValue)
	}

	return
}

func ListFieldWithTagByCustom(ptr interface{}, tagName string, customFieldNames []string) (fieldValues []reflect.Value) {
	v := reflect.ValueOf(ptr)
	t := reflect.TypeOf(ptr)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	customFieldNameMap := make(map[string]bool)
	for _, customFieldName := range customFieldNames {
		customFieldNameMap[customFieldName] = true
	}

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {
			newFieldValues := ListFieldWithTag(fieldValue.Interface(), tagName)
			fieldValues = append(fieldValues, newFieldValues...)
			continue
		}
		tag := fieldType.Tag

		customValue, ok := tag.Lookup(tagName)
		if !ok {
			continue
		}
		if !customFieldNameMap[customValue] {
			continue
		}

		fieldValues = append(fieldValues, fieldValue)
	}

	return
}
