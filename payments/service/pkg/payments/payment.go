package payments

import (
	"encoding/json"
	"strconv"
	"strings"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/database"
	"telython/pkg/eplidr"
)

type Payment struct {
	Id           uint64 // Id unique value for a payment
	Sender       string // Sender is sender username
	Receiver     string // Receiver is receiver username
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

func New(From string, To string, currencyFrom *currency.Currency, currencyTo *currency.Currency, timestamp uint64) *Payment {
	payment := Payment{
		Id:           fnv64(From + To + strconv.FormatUint(timestamp, 10)),
		Sender:       From,
		Receiver:     To,
		CurrencyFrom: currencyFrom,
		CurrencyTo:   currencyTo,
		Timestamp:    timestamp,
	}
	return &payment
}

func (payment *Payment) Commit() error {
	if payment.Sender == "system" {
		receiverShardNum := database.Accounts.Table.GetShardNum(fnv64(payment.Receiver))
		err := database.Payments.GetShard(receiverShardNum).Put(
			eplidr.PlainToColumns(
				[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
				[]interface{}{payment.Id, payment.Sender, payment.Receiver, payment.CurrencyFrom.Amount.String(), payment.CurrencyTo.Amount.String(), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
			),
		)
		if err != nil {
			return err
		}
		return nil
	} else if payment.Receiver == "system" {
		senderShardNum := database.Accounts.Table.GetShardNum(fnv64(payment.Sender))
		err := database.Payments.GetShard(senderShardNum).Put(
			eplidr.PlainToColumns(
				[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
				[]interface{}{payment.Id, payment.Sender, payment.Receiver, payment.CurrencyFrom.Amount.String(), payment.CurrencyTo.Amount.String(), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
			),
		)
		if err != nil {
			return err
		}
		return nil
	} else {
		senderShardNum := database.Accounts.Table.GetShardNum(fnv64(payment.Receiver))
		receiverShardNum := database.Accounts.Table.GetShardNum(fnv64(payment.Receiver))
		err := database.Payments.GetShard(senderShardNum).Put(
			eplidr.PlainToColumns(
				[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
				[]interface{}{payment.Id, payment.Sender, payment.Receiver, payment.CurrencyFrom.Amount.String(), payment.CurrencyTo.Amount.String(), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
			),
		)
		if err != nil {
			return err
		}
		if senderShardNum != receiverShardNum {
			err = database.Payments.GetShard(receiverShardNum).Put(
				eplidr.PlainToColumns(
					[]string{"id", "sender", "receiver", "amountFrom", "amountTo", "timestamp", "currencyFrom", "currencyTo"},
					[]interface{}{payment.Id, payment.Sender, payment.Receiver, payment.CurrencyFrom.Amount.String(), payment.CurrencyTo.Amount.String(), payment.Timestamp, payment.CurrencyFrom.Type.Id, payment.CurrencyTo.Type.Id},
				),
			)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (payment *Payment) Serialize() (string, error) {
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

func (payment Payment) SerializeReadable() (string, error) {
	/*serializedCurrencyFrom, err := json.Marshal(payment.CurrencyFrom)
	if err != nil {
		return "", err
	}
	serializedCurrencyTo, err := json.Marshal(payment.CurrencyTo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"Id": %d, "Sender": "%s", "Receiver": "%s", "Timestamp": %d, "CurrencyFrom": "%s", "CurrencyTo": "%s"}`, payment.Id, payment.Sender, payment.Receiver, payment.Timestamp,
		base64.StdEncoding.EncodeToString(serializedCurrencyFrom), base64.StdEncoding.EncodeToString(serializedCurrencyTo)), nil*/

	jsonData, err := json.Marshal(payment)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializePayment(serialized string) (Payment, error) {
	payment := Payment{}
	err := json.Unmarshal([]byte(serialized), &payment)
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func SerializePayments(payments []Payment) (string, error) {
	serializedPayments := ""
	for i := 0; i < len(payments); i++ {
		serialized, err := payments[i].Serialize()
		if err != nil {
			return "", err
		}
		if i == len(payments)-1 {
			serializedPayments += serialized
		} else {
			serializedPayments += serialized + "\n"
		}
	}
	return serializedPayments, nil
}

func DeserializePayments(serialized []byte) (*[]Payment, error) {
	if string(serialized) == "" {
		return &[]Payment{}, nil
	}
	var payments []Payment
	jsons := strings.Split(string(serialized), ",")
	for i := 0; i < len(jsons); i++ {
		payment, err := DeserializePayment(jsons[i])
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return &payments, nil
}
