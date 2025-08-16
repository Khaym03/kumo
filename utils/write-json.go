package utils

import (
	"encoding/json"
	"os"
)

func WriteJSONFile(path string, data any, pretty bool) error {
	var jsonData []byte
	var err error

	if pretty {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}
