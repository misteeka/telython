package currency

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
)

var Types map[string]*Type
var typesByCode map[uint64]*Type

func init() {
	Types = make(map[string]*Type)
	Types["USD"] = &Type{
		Symbol:   "USD",
		Id:       0,
		Decimals: new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil),
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

func GetCurrency(symbol string, amount uint64, precision uint64) *Currency {
	Type := Types[symbol]
	return &Currency{
		Type:   Type,
		Amount: new(big.Int).Div(new(big.Int).Mul(new(big.Int).SetUint64(amount), Type.Decimals), new(big.Int).Exp(big.NewInt(10), new(big.Int).SetUint64(precision), nil)),
	}
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
	if currency.Amount == nil {
		return ""
	}
	if currency.Type == nil {
		return ""
	}
	return new(big.Float).Quo(new(big.Float).SetInt(currency.Amount), new(big.Float).SetInt(currency.Type.Decimals)).String() + " " + currency.Type.Symbol
}

func (currency *Currency) Serialize() (uint64, string) {
	return currency.Type.Id, base64.StdEncoding.EncodeToString(currency.Amount.Bytes())
}

func (currency *Currency) Json() string {
	return fmt.Sprintf(`{"Type": %d, "Amount": "%s"}`, currency.Type.Id, currency.Amount)
}

func Deserialize(currencyCode uint64, amountString string) (*Currency, error) {
	amount, ok := new(big.Int).SetString(amountString, 10)
	if !ok {
		return nil, errors.New("Wrong amount " + amountString)
	}
	return &Currency{
		Type:   typesByCode[currencyCode],
		Amount: amount,
	}, nil
}
