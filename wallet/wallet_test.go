package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"testing"
)

const (
	testKey       string = "30770201010420f3bcf606539ffbca0186ebc47731b85677860a1f33ce857af4829bc69f1acd39a00a06082a8648ce3d030107a144034200043bccfbe8b721bb0d1f5c7d8917585da9d85bb840be05fd1c1d00606979e57528f2b63c4178ca5459a6666f936b85298ae562f3bb15e22cf207c44b87585fd58b"
	testHash      string = "395727ca97a9d1e0ac2d21bac0d8f928859f15d775d03c29b1a928714a8fde0c"
	testSignature string = "bea7c82061e4e14d08bb6ea12a5764afb52efac9c77986e4b653ce22a15cea45c964b565261cf02ee6cc33ae6ccdb5b3cd097576203e85dd213fc4dd20dbc5be"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	keyBytes, _ := hex.DecodeString(testKey)
	w.privateKey, _ = x509.ParseECPrivateKey(keyBytes)
	w.Address = keyToAddress(w.privateKey)
	return w
}

func TestSign(t *testing.T) {
	signature := Sign(testHash, makeTestWallet())
	t.Run("Signature is hex encoded", func(t *testing.T) {
		_, err := hex.DecodeString(signature)
		if err != nil {
			t.Errorf("Could not decode hex string: %s", err.Error())
		}
	})
}

func TestVerify(t *testing.T) {
	w := makeTestWallet()
	type test struct {
		payload string
		ok      bool
	}
	incorrectHash := "1" + testHash[1:]
	tests := []test{
		{testHash, true},
		{incorrectHash, false},
	}
	for _, tc := range tests {
		ok := Verify(tc.payload, testSignature, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and test case payload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx") // not a hex encoding
	if err == nil {
		t.Error("restoreBigInts() should return error when given a non-hexadecimal string")
	}
}
