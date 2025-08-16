package utils

import (
	"encoding/json"
)

// ToJsonString serializes any Go value into a formatted JSON string.
// If the serialization fails (e.g., due to an unsupported type), it returns
// a literal string "{}".
func ToJSONString(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
