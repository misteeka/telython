package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/valyala/fastjson"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"syscall"
	"telython/payments/service/pkg/log"
	"telython/payments/service/pkg/payments"
	"telython/payments/service/pkg/status"
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

func SendPayment(sender uint64, receiver uint64, amount uint64, password string) (status.Status, error) {
	json, err := Post("sendPayment", fmt.Sprintf(`{"sender":%d,"receiver":%d,"amount":%d, "password":"%s"}`, sender, receiver, amount, password))
	return getStatus(json), err
}
func GetBalance(accountId uint64, password string) (uint64, status.Status, error) {
	json, err := Get("getBalance?a=" + strconv.FormatUint(accountId, 10) + "&p=" + password)
	return json.GetUint64("data"), getStatus(json), err
}
func CreateAccount(username string, password string, currency int) (uint64, status.Status, error) {
	json, err := Post("createAccount", fmt.Sprintf(`{"username":"%s","password":"%s","currency":%d}`, username, password, currency))
	return json.GetUint64("data"), getStatus(json), err
}

func gunzipWrite(w io.Writer, data []byte) error {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		gr.Close()
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		gr.Close()
		return err
	}
	gr.Close()
	return nil
}

func GetHistory(accountId uint64) ([]payments.Payment, error) { // TODO timestamp and currency mismatch in SerializeRreadable
	f, err := syscall.Open(fmt.Sprintf("./client/data/txs/%d.gz", accountId), syscall.O_RDONLY, 0644)
	if err != nil {
		log.WarnLogger.Println("No transaction file found. Use loadHistory")
		return nil, nil
	}
	var data []byte
	buf := make([]byte, 4096)
	n, err := syscall.Read(f, buf)
	if err != nil {
		return nil, err
	}
	fmt.Println(n)
	if n == 0 {
		log.WarnLogger.Println("Broken transaction files. Use loadHistory")
		return nil, nil
	}
	for n > 0 {
		data = append(data, buf[:n]...)
		n, err = syscall.Read(f, buf)
		if err != nil {
			return nil, err
		}
	}
	var paymentsBytes bytes.Buffer
	err = gunzipWrite(&paymentsBytes, data)
	if err != nil {
		return nil, err
	}
	return *payments.DeserializePayments(paymentsBytes.Bytes()), nil
}
func LoadHistory(accountId uint64, password string) (status.Status, error) {
	json, err := Get("getHistory?a=" + strconv.FormatUint(accountId, 10) + "&p=" + password)
	if err != nil {
		return 0, err
	}
	paymentsBytes, err := base64.StdEncoding.DecodeString(string(json.GetStringBytes("data")))
	if err != nil {
		return status.Status(json.GetInt("status")), err
	}

	// Write with BestSpeed.
	f, err := os.Create(fmt.Sprintf("./client/data/txs/%d.gz", accountId))
	if err != nil {
		return status.Status(json.GetInt("status")), err
	}
	w, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
	if err != nil {
		return status.Status(json.GetInt("status")), err
	}
	_, err = w.Write(paymentsBytes)
	if err != nil {
		return status.Status(json.GetInt("status")), err
	}
	err = w.Close()
	if err != nil {
		return status.Status(json.GetInt("status")), err
	}
	err = f.Close()
	if err != nil {
		return status.Status(json.GetInt("status")), err
	}

	return status.Status(json.GetInt("status")), nil
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
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			password := args[1]
			start := time.Now()
			balance, status, err := GetBalance(id, password)
			print(balance, status, err, start)
		} else if strings.Compare("createAccount", cmd) == 0 {
			if len(args) < 3 {
				fmt.Println("Wrong args")
				continue
			}
			username := args[0]
			password := args[1]
			currency, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			start := time.Now()
			data, status, err := CreateAccount(username, password, currency)
			print(data, status, err, start)
		} else if strings.Compare("sendPayment", cmd) == 0 {
			if len(args) < 4 {
				fmt.Println("Wrong args")
				continue
			}
			sender, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			receiver, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			password := args[3]
			start := time.Now()
			status, err := SendPayment(sender, receiver, 1, password)
			for i := 1; i < int(amount); i++ {
				SendPayment(sender, receiver, 1, password)
			}
			print(nil, status, err, start)
		} else if strings.Compare("loadHistory", cmd) == 0 {
			if len(args) < 2 {
				fmt.Println("Wrong args")
				continue
			}
			account, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			password := args[1]
			start := time.Now()
			status, err := LoadHistory(account, password)
			print(nil, status, err, start)
		} else if strings.Compare("getHistory", cmd) == 0 {
			if len(args) < 1 {
				fmt.Println("Wrong args")
				continue
			}
			account, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			start := time.Now()
			result, err := GetHistory(account)
			print(result, -1, err, start)
		} else {
			fmt.Println("Unknown command.")
		}
	}

}
