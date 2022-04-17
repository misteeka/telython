package main

import (
	"bufio"
	"fmt"
	"github.com/valyala/fastjson"
	"math"
	"os"
	"strconv"
	"strings"
	"telython/payments/gates/ethgate/pkg/status"
	"time"
)

func statusToString(s status.Status) string {
	if s == status.SUCCESS {
		return "SUCCESS"
	}
	if s == status.INVALID_REQUEST {
		return "INVALID REQUEST"
	}
	if s == status.INTERNAL_SERVER_ERROR {
		return "INTERNAL SERVER ERROR"
	}
	if s == status.AUTHORIZATION_FAILED {
		return "AUTHORIZATION FAILED"
	}
	if s == status.INVALID_CURRENCY_CODE {
		return "INVALID CURRENCY CODE"
	}
	if s == status.CURRENCY_CODE_MISMATCH {
		return "CURRENCY CODE MISMATCH"
	}
	if s == status.WRONG_AMOUNT {
		return "WRONG AMOUNT"
	}
	if s == status.NOT_FOUND {
		return "NOT FOUND"
	}
	if s == status.INSUFFICIENT_FUNDS {
		return "INSUFFICIENT FUNDS"
	}
	if s == status.TOO_MANY_REQUESTS {
		return "TOO MANY REQUESTS"
	}
	return fmt.Sprintf("%v", s)
}

func getStatus(value *fastjson.Value) status.Status {
	return status.Status(value.GetInt("status"))
}

func print(data interface{}, status status.Status, err error, start time.Time) {
	duration := math.Round((float64(time.Now().Sub(start).Microseconds())/1000.0)*100) / 100.0
	if err != nil {
		fmt.Println("ERR: " + err.Error())
		return
	}
	if status != -1 {
		fmt.Println(fmt.Sprintf("Status: %s", statusToString(status)))
	}
	if data != nil {
		fmt.Println(fmt.Sprintf("Data: %v", data))
	}
	fmt.Println(fmt.Sprintf("Completed in %f ms", duration))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Telython Ethgate")
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
		if strings.Compare("createWallet", cmd) == 0 {
			if len(args) < 1 {
				fmt.Println("Wrong args")
				continue
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			//password := args[1]
			start := time.Now()
			wallet, status, err := CreateWallet(id)
			print(wallet, status, err, start)
		} else if strings.Compare("getAddress", cmd) == 0 {
			if len(args) < 1 {
				fmt.Println("Wrong args")
				continue
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			//password := args[1]
			start := time.Now()
			address, status, err := GetAddress(id)
			print(address, status, err, start)
		} else if strings.Compare("getPrivate", cmd) == 0 {
			if len(args) < 1 {
				fmt.Println("Wrong args")
				continue
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			//password := args[1]
			start := time.Now()
			private, status, err := GetPrivate(id)
			print(private, status, err, start)
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
