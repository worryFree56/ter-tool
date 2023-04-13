package tool

import (
	"bytes"
	"encoding/json"
)

func PrettyPrint(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		return ""
	}
	return out.String()
}
