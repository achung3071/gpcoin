package utils

import (
	"encoding/hex"
	"fmt"
	"reflect"
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

func TestToBytes(t *testing.T) {
	t.Run("Output is a slice of bytes", func(t *testing.T) {
		input := "test"
		bytes := ToBytes(input)
		k := reflect.TypeOf(bytes).Kind()
		if k != reflect.Slice {
			t.Errorf("Expected a slice of bytes, got %s", k)
		}
	})
}
