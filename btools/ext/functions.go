package ext

import "reflect"

func Contains[T comparable](slice []T, item T) bool {
	if slice == nil {
		return false
	}
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	} else {
		return b
	}
}

func Union[T comparable](a, b []T) []T {
	uniqueMap := make(map[T]struct{})
	var result []T

	for _, val := range a {
		uniqueMap[val] = struct{}{}
	}

	for _, val := range b {
		uniqueMap[val] = struct{}{}
	}

	for key := range uniqueMap {
		result = append(result, key)
	}

	return result
}

func ExtractField[T any, V any](items []T, fieldName string) []V {
	var result []V
	for _, item := range items {
		val := reflect.ValueOf(item)

		// Разыменовываем указатель, если это нужно
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// Проверяем, есть ли поле
		field := val.FieldByName(fieldName)
		if field.IsValid() && field.CanInterface() {
			if v, ok := field.Interface().(V); ok {
				result = append(result, v)
			}
		}
	}
	return result
}

func Diff[T comparable](slice1, slice2 []T) []T {
	var result []T

	set := make(map[T]struct{})
	for _, v := range slice1 {
		set[v] = struct{}{}
	}

	for _, v := range slice2 {
		if _, exists := set[v]; !exists {
			result = append(result, v)
		}
	}
	return result
}
