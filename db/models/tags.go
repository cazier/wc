package models

import (
	"reflect"
	"strings"
)

func TrimFields(endpoint string, items []interface{}) []map[string]any {
	result := make([]map[string]any, len(items))

	for index, element := range items {
		result[index] = trim(endpoint, element)
	}

	return result
}

func trim(endpoint string, i interface{}) map[string]any {
	contains := func(slice []string, search string) bool {
		for _, element := range slice {
			if element == search {
				return true
			}
		}
		return false
	}

	results := make(map[string]any)

	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		if result, err := t.Field(i).Tag.Lookup("endpoints"); err {
			if contains(strings.Split(result, ","), endpoint) {
				if key, err := t.Field(i).Tag.Lookup("json"); err && key != "-" {
					results[key] = v.Field(i).Interface()
				} else if key == "" {
					results[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
				}
			}
		}
	}
	return results
}
