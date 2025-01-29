package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKey() *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	fmt.Println("GenerateKey", hex.EncodeToString(crypto.FromECDSA(key)))
	return key
}

func LoadOrGenerateKey(key string) *ecdsa.PrivateKey {
	if key != "" {
		load, err := crypto.HexToECDSA(key)
		if err == nil {
			return load
		}
	}
	return GenerateKey()
}
