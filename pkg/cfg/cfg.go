package cfg

import (
	"fmt"
	"github.com/valyala/fastjson"
	"syscall"
	"telython/pkg/log"
)

const configDir = "./config.json"

var (
	Value *fastjson.Value
)

func GetString(key string) string {
	return string(Value.GetStringBytes(key))
}
func GetInt(key string) int {
	return Value.GetInt(key)
}

func LoadConfig() error {
	var json string
	f, err := syscall.Open(configDir, syscall.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(1)
		log.ErrorLogger.Print(err.Error())
		return err
	}
	buf := make([]byte, 4096)
	n, err := syscall.Read(f, buf)
	if err != nil {
		fmt.Println(2)
		log.ErrorLogger.Print(err.Error())
		return err
	}
	for n > 0 {
		json += string(buf[:n])
		n, err = syscall.Read(f, buf)
		if err != nil {
			fmt.Println(3)
			log.ErrorLogger.Print(err.Error())
			return err
		}
	}
	var p fastjson.Parser
	Value, err = p.Parse(json)
	if err != nil {
		fmt.Println(1)
		log.ErrorLogger.Print(err.Error())
		return err
	}
	return nil
}
