package to

import (
	"encoding/json"
	"fmt"
)

// String converts the object to a string.
func String(o interface{}) string {
	res := fmt.Sprintf("%v", o)
	return res
}

// JSON converts the object to a valid JSON string.
func JSON(o interface{}) (string, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LogString converts the object to a valid JSON string or string
func LogString(obj interface{}) (logStr string) {
	strBytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}
	return string(strBytes)
}
