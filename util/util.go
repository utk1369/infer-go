package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// given a map, pull a property from it at some deeply nested depth
// this reimplements (most of) JS `Pluck` in go: https://github.com/gjohnson/pluck
func Pluck(o map[string]interface{}, path string) interface{} {

	// support dots for now because thats all we need
	parts := strings.Split(path, ".")

	if len(parts) == 1 && o[parts[0]] != nil {
		// if there is only one part, just return that property value
		return o[parts[0]]
	} else if len(parts) > 1 && o[parts[0]] != nil {
		var prev map[string]interface{}
		var ok bool
		if prev, ok = o[parts[0]].(map[string]interface{}); !ok {
			// not an object type! ...or a map, yeah, that.
			return nil
		}

		for i := 1; i < len(parts)-1; i += 1 {
			// we need to check the existence of another
			// map[string]interface for every property along the way
			cp := parts[i]

			if prev[cp] == nil {
				// didn't find the property, it's missing
				return nil
			}

			var ok bool
			//check if prev[cp] is a slice and handle accordingly
			if reflect.TypeOf(prev[cp]).Kind() == reflect.Slice {
				idx, err := strconv.Atoi(parts[i+1])
				if err != nil {
					return nil
				}
				if prev, ok = reflect.ValueOf(prev[cp]).Index(idx).Interface().(map[string]interface{}); !ok {
					return nil
				}
				i += 1
			} else if prev, ok = prev[cp].(map[string]interface{}); !ok {
				return nil
			}
		}

		if prev[parts[len(parts)-1]] != nil {
			fmt.Println("Pluck path: ", path, " val: ", prev[parts[len(parts)-1]])
			return prev[parts[len(parts)-1]]
		} else {
			return nil
		}
	}

	return nil
}

func DeepClone(obj interface{}, copy interface{}) {
	x, _ := json.Marshal(obj)
	json.Unmarshal(x, &copy)
}
