package utils

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

var jsonTool = jsoniter.ConfigCompatibleWithStandardLibrary

//JSONMarshal .
func JSONMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

//JSONMarshalToString .
func JSONMarshalToString(v interface{}) (string, error) {
	return jsonTool.MarshalToString(v)
}

//JSONUnMarshal .
func JSONUnMarshal(str string, v interface{}) error {
	return jsonTool.UnmarshalFromString(str, v)
}
