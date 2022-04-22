package utils

import (
	"encoding/base64"
	"math/big"
	"strconv"
)

func Fnv64(key string) uint64 {
	hash := uint64(4332272522)
	const prime64 = uint64(33555238)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime64
		hash ^= uint64(key[i])
	}
	return hash
}

func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBigInt(base64Int string) (*big.Int, error) {
	bytes, err := base64.StdEncoding.DecodeString(base64Int)
	if err != nil {
		return nil, err
	}
	bigint := new(big.Int).SetBytes(bytes)
	return bigint, nil
}
func EncodeBigInt(int *big.Int) string {
	return base64.StdEncoding.EncodeToString(int.Bytes())
}

func ParseUint(v string) (uint64, error) {
	return strconv.ParseUint(v, 10, 64)
}
