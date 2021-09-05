package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
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

func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}
	// Table testing for Splitter()
	tests := []test{
		{input: "11:0:6", sep: ":", index: 2, output: "6"},
		{input: "11:0:6", sep: ":", index: 5, output: ""},
		{input: "11:0:6", sep: ":", index: -2, output: ""},
		{input: "11:0:6", sep: "/", index: 0, output: "11:0:6"},
	}
	for _, tc := range tests {
		result := Splitter(tc.input, tc.sep, tc.index)
		if result != tc.output {
			t.Errorf("Expected %s, got %s", tc.output, result)
		}
	}
}

func TestErrorHandler(t *testing.T) {
	oldPanic := panic
	defer func() {
		panic = oldPanic
	}()
	called := false
	panic = func(v ...interface{}) {
		called = true
	}
	newError := errors.New("test error")
	ErrorHandler(newError)
	if !called {
		t.Error("ErrorHandler did not call the panic function")
	}
}

func TestFromBytes(t *testing.T) {
	type test struct {
		Test string
	}
	input := test{"test"}
	var restored test
	b := ToBytes(input)
	FromBytes(&restored, b)
	if !reflect.DeepEqual(input, restored) {
		t.Error("FromBytes was unable to restore the input.")
	}
}

func TestToJSON(t *testing.T) {
	type test struct {
		Test string
	}
	input := test{"test"}
	b := ToJSON(input)
	t.Run("Output is a slice of bytes", func(t *testing.T) {
		k := reflect.TypeOf(b).Kind()
		if k != reflect.Slice {
			t.Errorf("Expected a slice of bytes, got %s", k)
		}
	})
	t.Run("Output can be unmarshaled to a struct", func(t *testing.T) {
		var restored test
		err := json.Unmarshal(b, &restored)
		if err != nil {
			t.Errorf("Error while unmarshaling JSON: %s", err.Error())
		}
		if !reflect.DeepEqual(input, restored) {
			t.Error("Unmarshaled JSON is not the same as the input.")
		}
	})
}
