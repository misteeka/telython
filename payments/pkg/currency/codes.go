package currency

import (
	"encoding/base64"
	"math/big"
)

var Types map[string]*Type
var typesByCode map[uint64]*Type

func init() {
	Types = make(map[string]*Type)
	Types["USD"] = &Type{
		Symbol:   "USD",
		Id:       0,
		Decimals: new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil),
	}
	Types["ETH"] = &Type{
		Symbol:   "ETH",
		Id:       1,
		Decimals: new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
	}
	typesByCode = make(map[uint64]*Type)
	for _, Type := range Types {
		typesByCode[Type.Id] = Type
	}
}

func FromCode(code uint64) *Type {
	return typesByCode[code]
}

type Type struct {
	Symbol   string
	Id       uint64
	Decimals *big.Int
}

type Currency struct {
	Type   *Type
	Amount *big.Int
}

func (currency *Currency) Readable() string {
	return new(big.Float).Quo(new(big.Float).SetInt(currency.Amount), new(big.Float).SetInt(currency.Type.Decimals)).String()
}

func (currency *Currency) Serialize() (uint64, string) {
	return currency.Type.Id, base64.StdEncoding.EncodeToString(currency.Amount.Bytes())
}

func Deserialize(currencyCode uint64, serialized string) (*Currency, error) {
	bytes, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}
	return &Currency{
		Type:   typesByCode[currencyCode],
		Amount: new(big.Int).SetBytes(bytes),
	}, nil
}
