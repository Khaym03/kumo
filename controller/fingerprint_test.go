package controller

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFingerprint(t *testing.T) {
	filename := filepath.Join("..", "data", "fingerprints.json")
	fingerprints, err := LoadFingerprints(filename)

	assert.NotNil(t, fingerprints, "nil fingerprints")
	assert.Nil(t, err)
}
