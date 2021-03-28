package pnghider

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHidePayload(t *testing.T) {
	pngBytes, err := HidePayload([]byte("IDAB"), []byte("dabbed on"), "test.png")
	assert.Nil(t, err)
	os.WriteFile("test-out.png", pngBytes, 0644)
}
