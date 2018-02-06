package arraytool

import "reflect"

//InArray reports whether value exists and inside array
func InArray(value interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

//RemoveItem remove value from array
func RemoveItem(value string, array []string) (results []string) {

	for _, item := range array {
		if item != value {
			results = append(results, item)
		}
	}

	return results
}
