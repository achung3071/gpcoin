package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/achung3071/gpcoin/utils"
)

// How signature verification works:
// 1. Generate public-private key pair and hash the data to be signed
// 2. Add the hash with the private key to generate the signature
// 3. To verify the signature was created by a specific private key,
//    add the (hash + signature + public key) to get a true/false value

const (
	walletFileName string = "gpcoin.wallet"
)

// Interface for isolating filesystem side effects (allows for unit testing)
type fileLayer interface {
	walletFileExists() bool
	writeFile(name string, data []byte) error
	readFile(name string) ([]byte, error)
}

// Struct for implementing fileLayer interface
type layer struct{}

func (layer) walletFileExists() bool {
	_, err := os.Stat(walletFileName)
	exists := !os.IsNotExist(err)
	return exists
}

func (layer) writeFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644) // read-write perms
}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Note that the wallet address is actualy the public key associated with
// the private key (which people can use to verify that you signed transactions)
type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet
var files fileLayer = layer{}

// NON-MUTATING FUNCTIONS
// Access singleton instance of wallet
func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if files.walletFileExists() {
			// yes -> load existing wallet
			w.privateKey = restoreKey()
		} else {
			// no -> create new wallet file
			w.privateKey = createPrivateKey()
			commitWallet(w)
		}
		w.Address = keyToAddress(w.privateKey)
	}
	return w
}

// Save wallet with private key & read-write permissions
func commitWallet(w *wallet) {
	privKeyBytes, err := x509.MarshalECPrivateKey(w.privateKey)
	utils.ErrorHandler(err)
	err = files.writeFile(walletFileName, privKeyBytes)
	utils.ErrorHandler(err)
}

// Creates a new private key
func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.ErrorHandler(err)
	return privateKey
}

// Get address (public key) from private key
func keyToAddress(k *ecdsa.PrivateKey) string {
	// Note that since PublicKey is an embedded struct in PrivateKey,
	// all its fields (X and Y) are "promoted" to PrivateKey, making
	// them directly accessible.
	return encodeBigInts(k.X.Bytes(), k.Y.Bytes())
}

// Encodes big ints (r/s for signature or x/y for public key) into hex string
func encodeBigInts(a, b []byte) string {
	return fmt.Sprintf("%x", append(a, b...))
}

// Restore two big ints (either r/s for signature or x/y for public key)
// from the hexadecimal string encoding
func restoreBigInts(encoding string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(encoding)
	if err != nil {
		return nil, nil, err
	}
	aBytes, bBytes := bytes[:len(bytes)/2], bytes[len(bytes)/2:]
	var a, b big.Int
	a.SetBytes(aBytes)
	b.SetBytes(bBytes)
	return &a, &b, nil
}

// Restore a private key from a wallet file
func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := files.readFile(walletFileName)
	utils.ErrorHandler(err)
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.ErrorHandler(err)
	return key
}

// Sign a hash (i.e., a new transaction id) using wallet's private key
func Sign(hash string, w *wallet) string {
	hashBytes, err := hex.DecodeString(hash)
	utils.ErrorHandler(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, hashBytes)
	utils.ErrorHandler(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

// Verify a hash (transaction) has been signed by the private key (wallet) associated w/ address
func Verify(hash, signature, address string) bool {
	hashBytes, err := hex.DecodeString(hash)
	utils.ErrorHandler(err)
	r, s, err := restoreBigInts(signature)
	utils.ErrorHandler(err)
	x, y, err := restoreBigInts(address)
	utils.ErrorHandler(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return ecdsa.Verify(&publicKey, hashBytes, r, s)
}
