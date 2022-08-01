package pkg

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// SetNestedField update or add a key-value field
func SetNestedField(obj map[string]interface{}, target interface{}, fields ...string) (err error) {
	m := obj

	for i, field := range fields[:len(fields)-1] {
		rawField := field
		if isIndexedField(field) {
			field = removeIndexFromField(field)
		}

		if val, ok := m[field]; ok {
			switch typedVal := val.(type) {
			case map[string]interface{}:
				m = typedVal
			case []interface{}:
				index := getIndex(rawField)
				if index >= 0 && index < len(typedVal) {
					m, ok = typedVal[index].(map[string]interface{})
				}
			default:
				ok = false
			}

			if !ok {
				err = fmt.Errorf("unexpect type in filed %s", strings.Join(fields[:i+1], ""))
				break
			}
		} else {
			emptyMap := make(map[string]interface{})
			m[field] = emptyMap
			m = emptyMap
		}
	}
	m[fields[len(fields)-1]] = target
	return
}

// NestedFieldAsString returns the unstructured map value as string format
func NestedFieldAsString(obj map[string]interface{}, fields ...string) (strVal string, ok bool, err error) {
	var result interface{}
	if result, ok, err = NestedField(obj, fields...); ok && err == nil {
		if strVal, ok = result.(string); !ok {
			err = fmt.Errorf("expect string type, got %s", reflect.TypeOf(result))
		}
	}
	return
}

// NestedFieldAsInt returns the unstructured map value as int format
func NestedFieldAsInt(obj map[string]interface{}, fields ...string) (intVal int, ok bool, err error) {
	var result interface{}
	if result, ok, err = NestedField(obj, fields...); ok && err == nil {
		if intVal, ok = result.(int); !ok {
			err = fmt.Errorf("expect int type, got %s", reflect.TypeOf(result))
		}
	}
	return
}

// NestedField returns the unstructured map value
func NestedField(obj map[string]interface{}, fields ...string) (interface{}, bool, error) {
	var val interface{} = obj

	for i, field := range fields {
		if val == nil {
			return nil, false, nil
		}

		rawField := field
		if isIndexedField(field) {
			field = removeIndexFromField(field)
		}

		if m, ok := val.(map[string]interface{}); ok {
			if val, ok = m[field]; !ok {
				return nil, false, nil
			}

			if isIndexedField(rawField) {
				if m, ok := val.([]map[string]interface{}); ok {
					index := getIndex(rawField)
					if index >= 0 && index < len(m) {
						val = m[index]
					} else {
						return nil, false, fmt.Errorf("invalid index: %d, field: %v", index, fields[0:i])
					}
				} else if m, ok := val.([]interface{}); ok {
					index := getIndex(rawField)
					if index >= 0 && index < len(m) {
						val = m[index]
					} else {
						return nil, false, fmt.Errorf("invalid index: %d, field: %v", index, fields[0:i])
					}
				} else if m, ok := val.([]string); ok {
					index := getIndex(rawField)
					if index >= 0 && index < len(m) {
						val = m[index]
					} else {
						return nil, false, fmt.Errorf("invalid index: %d, field: %v", index, fields[0:i])
					}
				}
			}
		} else {
			return nil, false, fmt.Errorf("expect %s type is map[string]interface{}", fields[0:i])
		}
	}
	return val, true, nil
}

func removeIndexFromField(field string) string {
	items := strings.Split(field, "[")
	if len(items) == 2 {
		return items[0]
	}
	return field
}

func isIndexedField(field string) (matched bool) {
	matched, _ = regexp.MatchString(".*\\[\\d\\]", field)
	return
}

// getIndex returns the index number from format xxx[10]
func getIndex(word string) (index int) {
	index = -1 // this is an invalid value

	items := strings.Split(word, "[")
	if len(items) == 2 {
		index, _ = strconv.Atoi(strings.TrimSuffix(items[1], "]"))
	}
	return
}
