package keypair

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
)

// KeyPair defines the pair of public and private key
type KeyPair struct {
	PublicKey  *btcec.PublicKey
	PrivateKey *btcec.PrivateKey
}

// FromPrivateKey takes a hex encoded private key  and returns a KeyPair instance
func FromPrivateKey(hexPriv string) (*KeyPair, error) {
	privateKeyBytes, err := hex.DecodeString(hexPriv)
	if err != nil {
		return nil, err
	}
	privateKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	return &KeyPair{publicKey, privateKey}, nil
}
