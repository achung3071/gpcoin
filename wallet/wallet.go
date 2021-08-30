package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
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

// Note that the wallet address is actualy the public key associated with
// the private key (which people can use to verify that you signed transactions)
type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

// NON-MUTATING FUNCTIONS
// Access singleton instance of wallet
func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if walletFileExists() {
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

// Get address (public key) from private key
func keyToAddress(k *ecdsa.PrivateKey) string {
	// Note that since PublicKey is an embedded struct in PrivateKey,
	// all its fields (X and Y) are "promoted" to PrivateKey, making
	// them directly accessible.
	publicKeyBytes := append(k.X.Bytes(), k.Y.Bytes()...)
	return fmt.Sprintf("%x", publicKeyBytes)
}

// Save wallet with private key & read-write permissions
func commitWallet(w *wallet) {
	privKeyBytes, err := x509.MarshalECPrivateKey(w.privateKey)
	utils.ErrorHandler(err)
	err = os.WriteFile(walletFileName, privKeyBytes, 0644)
	utils.ErrorHandler(err)
}

// Creates a new private key
func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.ErrorHandler(err)
	return privateKey
}

// Restore a private key from a wallet file
func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := os.ReadFile(walletFileName)
	utils.ErrorHandler(err)
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.ErrorHandler(err)
	return key
}

// Check if a wallet file exists
func walletFileExists() bool {
	_, err := os.Stat(walletFileName)
	exists := !os.IsNotExist(err) // check if error caused by nonexistent file
	return exists
}
