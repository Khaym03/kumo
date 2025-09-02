package controller

import (
	"encoding/json"
	"os"
)

type CombineFingerprint struct {
	Fingerprint map[string]any    `json:"fingerprint"`
	Headers     map[string]string `json:"headers"`
}

func LoadFingerprints(filepath string) ([]CombineFingerprint, error) {
	var cf []CombineFingerprint

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &cf)
	if err != nil {
		return nil, err
	}

	return cf, err
}
