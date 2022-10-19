package common

import "encoding/json"

func StructToString(obj interface{}) (newString string, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	newString = string(data)
	return
}
