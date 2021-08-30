package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/achung3071/gpcoin/utils"
)

const (
	privateKey string = "307702010104208204ea0102e28a6c245e4e1b511cb6440f8688f239115cc41b0053e993eeec93a00a06082a8648ce3d030107a1440342000491014d888c86de295e9b3eb1c4ef03f18240a0dcc184ea1f596e80bb92448c8e3f686a22328da52afda498a5e89179e4774887033bea830268ea74078cba9924"
	hashedMsg  string = "cb13313fb0d834900302ebee6fe7e4cc0ee48d568969549700a526e44a4f27e9"
	signature  string = "1d9cb6868e124052ae2e95687c2cc6dc88f95fde10b4d1d28a995b69f155e00254f099dc0afc6629329e2dc7bd821d11986b53f1c5b28c2ffd7d694b6d4098ae"
)

// How signature verification works:
// 1. Generate public-private key pair and hash the data to be signed
// 2. Add the hash with the private key to generate the signature
// 3. To verify the signature was created by a specific private key,
//    add the (hash + signature + public key) to get a true/false value

func Start() {
	// Get private key from hex string
	privKeyBytes, err := hex.DecodeString(privateKey)
	utils.ErrorHandler(err)
	privateKey, err := x509.ParseECPrivateKey(privKeyBytes)
	utils.ErrorHandler(err)

	// Get r and s (signature) from hex string
	sigBytes, err := hex.DecodeString(signature)
	utils.ErrorHandler(err)
	rBytes, sBytes := sigBytes[:len(sigBytes)/2], sigBytes[len(sigBytes)/2:]
	var r, s big.Int
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)

	// Get bytes of hash
	hashBytes, err := hex.DecodeString(hashedMsg)
	utils.ErrorHandler(err)

	// If private key used to sign hash to get signature, prints true.
	// Otherwise, prints false.
	ok := ecdsa.Verify(&privateKey.PublicKey, hashBytes, &r, &s)
	fmt.Println(ok)
}
