package payments

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/database"
	"telython/pkg/eplidr"
)

type Payment struct {
	Id           uint64 // Id unique value for a payment
	Sender       uint64 // Sender is fvn64 of sender username
	Receiver     uint64 // Receiver is fvn64 of receiver username
	Timestamp    uint64 // Timestamp UNIX timestamp in microseconds
	CurrencyFrom *currency.Currency
	CurrencyTo   *currency.Currency
}

func fnv64(key string) uint64 {
	hash := uint64(4332272522)
	const prime64 = uint64(33555238)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime64
		hash ^= uint64(key[i])
	}
	return hash
}

func New(From uint64, To uint64, currencyFrom *currency.Currency, currencyTo *currency.Currency, timestamp uint64) *Payment {
	payment := Payment{
		Id:           fnv64(strconv.FormatUint(From, 10) + strconv.FormatUint(To, 10) + strconv.FormatUint(timestamp, 10)),
		Sender:       From,
		Receiver:     To,
		CurrencyFrom: currencyFrom,
		CurrencyTo:   currencyTo,
		Timestamp:    timestamp,
	}
	return &payment
}

func (payment *Payment) Commit() error {
	if payment.Sender == 0 {
		receiverShardNum := database.Accounts.Table.GetShardNum(payment.Receiver)
		err := database.Payments.GetShard(receiverShardNum).Put(
			eplidr.PlainToColumns(
				[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
				[]interface{}{payment.Id, payment.Sender, payment.Receiver, base64.StdEncoding.EncodeToString(payment.CurrencyFrom.Amount.Bytes()), base64.StdEncoding.EncodeToString(payment.CurrencyTo.Amount.Bytes()), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
			),
		)
		if err != nil {
			return err
		}
		return nil
	} else if payment.Receiver == 0 {
		senderShardNum := database.Accounts.Table.GetShardNum(payment.Sender)
		err := database.Payments.GetShard(senderShardNum).Put(
			eplidr.PlainToColumns(
				[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
				[]interface{}{payment.Id, payment.Sender, payment.Receiver, base64.StdEncoding.EncodeToString(payment.CurrencyFrom.Amount.Bytes()), base64.StdEncoding.EncodeToString(payment.CurrencyTo.Amount.Bytes()), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
			),
		)
		if err != nil {
			return err
		}
		return nil
	} else {
		senderShardNum := database.Accounts.Table.GetShardNum(payment.Sender)
		receiverShardNum := database.Accounts.Table.GetShardNum(payment.Receiver)
		err := database.Payments.GetShard(senderShardNum).Put(
			eplidr.PlainToColumns(
				[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
				[]interface{}{payment.Id, payment.Sender, payment.Receiver, base64.StdEncoding.EncodeToString(payment.CurrencyFrom.Amount.Bytes()), base64.StdEncoding.EncodeToString(payment.CurrencyTo.Amount.Bytes()), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
			),
		)
		if err != nil {
			return err
		}
		if senderShardNum != receiverShardNum {
			err = database.Payments.GetShard(receiverShardNum).Put(
				eplidr.PlainToColumns(
					[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
					[]interface{}{payment.Id, payment.Sender, payment.Receiver, base64.StdEncoding.EncodeToString(payment.CurrencyFrom.Amount.Bytes()), base64.StdEncoding.EncodeToString(payment.CurrencyTo.Amount.Bytes()), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
				),
			)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (payment *Payment) Serialize() ([]byte, error) {
	return payment.SerializeReadable()
	/*
		buff := new(bytes.Buffer)
		err := binary.Write(buff, binary.BigEndian, payment.Id)
		if err != nil {
			return nil, err
		}
		err = binary.Write(buff, binary.BigEndian, payment.Sender)
		if err != nil {
			return nil, err
		}
		err = binary.Write(buff, binary.BigEndian, payment.Receiver)
		if err != nil {
			return nil, err
		}
		err = binary.Write(buff, binary.BigEndian, payment.Amount)
		if err != nil {
			return nil, err
		}
		err = binary.Write(buff, binary.BigEndian, payment.Timestamp)
		if err != nil {
			return nil, err
		}
		err = binary.Write(buff, binary.BigEndian, payment.Currency)
		if err != nil {
			return nil, err
		}
		return buff.Bytes(), nil*/
}

func (payment Payment) SerializeReadable() ([]byte, error) {
	jsonData, err := json.Marshal(payment)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func DeserializePayment(serialized []byte) (Payment, error) {
	payment := Payment{}
	err := json.Unmarshal(serialized, &payment)
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func SerializePayments(payments []Payment) ([]byte, error) {
	buff := new(bytes.Buffer)
	for i := 0; i < len(payments); i++ {
		serialized, err := payments[i].Serialize()
		if err != nil {
			return nil, err
		}
		buff.Write(serialized)
		buff.Write([]byte("\n"))
	}
	return buff.Bytes(), nil
}

func DeserializePayments(serialized []byte) (*[]Payment, error) {
	var payments []Payment
	for i := 0; i < (len(serialized) / 48); i++ {
		payment, err := DeserializePayment(serialized[i*48 : (i+1)*48])
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return &payments, nil
}
