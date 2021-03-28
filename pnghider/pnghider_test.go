package pnghider

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHidePayload(t *testing.T) {
	pic, err := os.ReadFile("../presentation/images/png_hex.png")
	if err != nil {
		t.Errorf("Error opening file: %v", err)
		t.FailNow()
	}
	// NOTE: capital first character marks this as a critical chunk
	pngBytes, err := HidePayload([]byte("iDAB"), []byte("dabbed on"), pic)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	home, _ := os.UserHomeDir()
}

func TestRecoverPayload(t *testing.T) {
	home, _ := os.UserHomeDir()
	pic, err := os.ReadFile(home + "/test-out.png")
	if err != nil {
		t.Errorf("Error opening file: %v", err)
		t.FailNow()
	}
	payload, err := RecoverPayload([]byte("iDAB"), pic)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	fmt.Printf("recovered payload: %s\n", string(payload))
	assert.Equal(t, "dabbed on", string(payload))
}