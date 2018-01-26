package arraytool

import "reflect"

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func RemoveItem(items []string, removeTarget string) (results []string) {

	for _, item := range items {
		if item != removeTarget {
			results = append(results, item)
		}
	}

	return results
}
