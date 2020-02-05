package strutils

import "encoding/json"

// StructToString converts struct to string by json.Marshal
func StructToString(v interface{}) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// StringToStruct does the reverse by json.Unmarshal
func StringToStruct(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), v)
}
