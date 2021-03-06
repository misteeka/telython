package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/client"
	"telython/payments/service/pkg/payments"
	"telython/pkg/http"
	"time"
)

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

func print(data interface{}, error *http.Error, err error, start time.Time) {
	duration := math.Round((float64(time.Now().Sub(start).Microseconds())/1000.0)*100) / 100.0
	if err != nil {
		fmt.Println("ERR: " + err.Error())
		return
	}
	fmt.Println(http.ToReadable(error))
	if data != nil {
		switch data := data.(type) {
		case []payments.Payment:
			if len(data) == 0 {
				return
			}
			fmt.Println("Payments: ")
			for i := 0; i < len(data); i++ {
				printable, err := data[i].SerializeReadable()
				if err != nil {
					fmt.Println("serialization error " + err.Error())
				} else {
					fmt.Println(string(printable))
				}
			}
		default:
			fmt.Println(fmt.Sprintf("Data: %v", data))
		}
	}
	fmt.Println(fmt.Sprintf("Completed in %f ms", duration))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Telython Pay Shell")
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ReplaceAll(text, "\r", "")
		parts := strings.Split(text, " ")
		if len(parts) < 1 {
			fmt.Println("Wrong command")
			continue
		}
		cmd := parts[0]
		args := parts[1:]
		if strings.Compare("getBalance", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			currencyCode, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			start := time.Now()
			amount, status, err := client.GetBalance(username, password, currencyCode)
			currency := &currency.Currency{
				Type:   currency.FromCode(currencyCode),
				Amount: amount,
			}
			print(currency.Readable(), status, err, start)
		} else if strings.Compare("createAccount", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			start := time.Now()
			data, status, err := client.CreateAccount(username, password)
			print(data, status, err, start)
		} else if strings.Compare("sendPayment", cmd) == 0 {
			if len(args) < 4 {
				fmt.Println("Wrong args")
				continue
			}
			sender := args[0]
			receiver := args[1]
			amount, ok := new(big.Int).SetString(args[2], 10)
			if !ok {
				fmt.Println("Wrong amount")
				continue
			}
			currencyFrom, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			currencyTo, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			password := args[5]
			start := time.Now()
			status, err := client.SendPayment(sender, receiver, &currency.Currency{
				Type:   currency.FromCode(currencyFrom),
				Amount: amount,
			}, currencyTo, password)
			print(nil, status, err, start)
		} else if strings.Compare("addPayment", cmd) == 0 {
			if len(args) < 4 {
				fmt.Println("Wrong args")
				continue
			}
			sender := args[0]
			receiver := args[1]
			amount, ok := new(big.Int).SetString(args[2], 10)
			if !ok {
				fmt.Println("Wrong amount")
				continue
			}
			currencyCode, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			secretKey := args[4]
			start := time.Now()
			requestError, err := client.AddPayment(sender, receiver, &currency.Currency{
				Type:   currency.FromCode(currencyCode),
				Amount: amount,
			}, secretKey)
			if err != nil {
				return
			}
			print(nil, requestError, err, start)
		} else if strings.Compare("getHistory", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			start := time.Now()
			result, error, err := client.GetHistory(username, password)
			print(result, error, err, start)
		} else if strings.Compare("accountInfo", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			start := time.Now()
			error, err := client.GetAccountInfo(username, password)
			print(nil, error, err, start)
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
