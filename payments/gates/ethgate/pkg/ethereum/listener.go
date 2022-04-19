package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"telython/payments/gates/ethgate/pkg/database"
	ethapi "telython/payments/gates/ethgate/pkg/ethereum/api"
	tpay "telython/payments/service/pkg/client"
	"telython/pkg/log"
)

var headers chan *types.Header
var sub ethereum.Subscription

func initEthereumClient() error {
	var err error
	Client, err = ethclient.Dial("ws://127.0.0.1:3334")
	return err
}

func initEthereumListener() error {
	var err error
	headers = make(chan *types.Header)
	sub, err = Client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return err
	}
	go newBlockHandler()
	return nil
}

func newBlockHandler() {
	for {
		select {
		case err := <-sub.Err():
			log.ErrorLogger.Println(err.Error())
			break
		case header := <-headers:
			newBlock, err := Client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.ErrorLogger.Println(err.Error())
				continue
			}
			log.InfoLogger.Println("New block #", newBlock.Number().Uint64())
			if newBlock.NumberU64() > 12 {
				block, err := Client.BlockByNumber(context.Background(), new(big.Int).Sub(newBlock.Number(), big.NewInt(12)))
				if err != nil {
					log.ErrorLogger.Println(err.Error())
					continue
				}
				for _, tx := range block.Transactions() {
					accountId, found, err := database.WalletToAccount.GetUint(ethapi.AddressToBase64(tx.To()), "id")
					if err != nil {
						log.ErrorLogger.Println(err.Error())
						continue
					}
					if !found {
						continue
					}
					paymentError, err := tpay.AddPayment(accountId, tx.Value(), "secret")
					if err != nil {
						log.ErrorLogger.Println(err.Error())
						continue
					}
					if paymentError != nil {
						log.WarnLogger.Println("Error ("+paymentError.Message+") during adding ", tx.Value().String(), "to", accountId, "("+tx.To().Hex()+")", "confirmations:", 12)
					} else {
						log.InfoLogger.Println("Added", tx.Value().String(), "to", accountId, "("+tx.To().Hex()+")", "confirmations:", 12)
					}
				}
			}
		}
	}
}