package utils

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	hash := Hash(struct{ test string }{test: "Test"})
	t.Run("Hash is always the same", func(t *testing.T) {
		expectedHash := "a4fbea6e8fec7ba1429275042bedbe3f43c708e448031efe4292be04d1f58325"
		if hash != expectedHash {
			t.Errorf("Expected %s, got %s", expectedHash, hash)
		}
	})
	t.Run("Hash is in hexadecimal format", func(t *testing.T) {
		_, err := hex.DecodeString(hash)
		if err != nil {
			t.Error("Hash is not in hexadecimal format")
		}
	})
}

func ExampleHash() {
	input := struct{ test string }{test: "Test"}
	hash := Hash(input)
	fmt.Println(hash)
	// Output: a4fbea6e8fec7ba1429275042bedbe3f43c708e448031efe4292be04d1f58325
}
