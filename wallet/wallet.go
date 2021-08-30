package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/achung3071/gpcoin/utils"
)

// How signature verification works:
// 1. Generate public-private key pair and hash the data to be signed
// 2. Add the hash with the private key to generate the signature
// 3. To verify the signature was created by a specific private key,
//    add the (hash + signature + public key) to get a true/false value

func Start() {
	// Generate privateKey-publicKey pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.ErrorHandler(err)

	// Generate the hash for the data
	data := "Keep this secret"
	hash := utils.Hash(data)
	hashInBytes, err := hex.DecodeString(hash)
	utils.ErrorHandler(err)

	// Signature generated in two parts (r,s)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashInBytes)
	utils.ErrorHandler(err)
	fmt.Printf("R: %d\nS: %d\n", r, s)

	// Verify key
	ok := ecdsa.Verify(&privateKey.PublicKey, hashInBytes, r, s)
	fmt.Println(ok)
}
