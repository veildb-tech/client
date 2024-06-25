/*
Copyright © 2024 Bridge Digital
*/
package util

func MapKeyByValue(m map[string]string, v string) string {
	var key string = ""

	for key, value := range m {
		if value == v {
			return key
		}
	}

	return key
}
