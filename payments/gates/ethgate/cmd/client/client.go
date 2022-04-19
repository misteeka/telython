package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"telython/payments/gates/ethgate/pkg/client"
	"telython/pkg/http"
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
			wallet, status, err := client.CreateWallet(id)
			if wallet == nil {
				print(nil, status, err, start)
			} else {
				print(wallet.GetAddressHEX(), status, err, start)
			}
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
			address, status, err := client.GetAddress(id)
			if address == nil {
				print(nil, status, err, start)
			} else {
				print(address.Hex(), status, err, start)
			}
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
			private, status, err := client.GetPrivate(id)
			print(private, status, err, start)
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
