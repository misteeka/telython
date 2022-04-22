package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"telython/payments/exchange/pkg/client"
	"telython/payments/pkg/currency"
	"telython/pkg/http"
	"telython/pkg/utils"
	"time"
)

func print(data interface{}, error *http.Error, err error, start time.Time) {
	duration := math.Round((float64(time.Now().Sub(start).Microseconds())/1000.0)*100) / 100.0
	if err != nil {
		fmt.Println("ERR: " + err.Error())
		return
	}
	if error != nil {
		fmt.Println(http.ToReadable(error))
	}
	if data != nil {
		fmt.Println(fmt.Sprintf("Data: %v", data))
	}
	fmt.Println(fmt.Sprintf("Completed in %f ms", duration))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Telython Exchange")
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
		if strings.Compare("getPrice", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			symbol := args[0]
			key := args[1]
			start := time.Now()
			price, requestError, err := client.GetPrice(symbol, key)
			if price == nil {
				print(nil, requestError, err, start)
			} else {
				print(price.Readable(), requestError, err, start)
			}
		} else if strings.Compare("convert", cmd) == 0 {
			if len(args) < 4 {
				fmt.Println("Wrong args")
				continue
			}
			fromCode, err := utils.ParseUint(args[0])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			to, err := utils.ParseUint(args[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			amount, ok := new(big.Int).SetString(args[2], 10)
			if !ok {
				fmt.Println("Wrong amount")
				continue
			}
			from := &currency.Currency{
				Type:   currency.FromCode(fromCode),
				Amount: amount,
			}
			key := args[3]
			start := time.Now()
			requestError, result, err := client.Convert(from, to, key)
			print(result.Readable(), requestError, err, start)
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
