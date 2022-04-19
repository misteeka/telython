package api

import (
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func PublicKeyBytesToAddress(publicKey []byte) *common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]
	result := common.BytesToAddress(address)
	return &result
}

func PublicBase64ToAddress(publicBase64 string) (*common.Address, error) {
	publicBytes, err := base64.StdEncoding.DecodeString(publicBase64)
	if err != nil {
		return nil, err
	}
	return PublicKeyBytesToAddress(publicBytes), nil
}

func Base64ToPrivate(privateBase64 string) (*ecdsa.PrivateKey, error) {
	privateBytes, err := base64.StdEncoding.DecodeString(privateBase64)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func PrivateToAddress(privateKey *ecdsa.PrivateKey) (*common.Address, bool) {
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, false
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address, true
}

func AddressToBase64(address *common.Address) string {
	return base64.StdEncoding.EncodeToString(address.Bytes())
}

func Base64ToAddress(base64Address string) (*common.Address, error) {
	addressBytes, err := base64.StdEncoding.DecodeString(base64Address)
	if err != nil {
		return nil, err
	}
	address := common.BytesToAddress(addressBytes)
	return &address, nil
}
