package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

var Client *ethclient.Client

func Init() error {
	err := initEthereumClient()
	if err != nil {
		return err
	}
	err = initEthereumListener()
	if err != nil {
		return err
	}
	return nil
}
