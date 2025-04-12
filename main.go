package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	type LogNumber struct {
		Number *big.Int
	}

	infura_wss := os.Getenv("INFURA_WSS")

	client, err := ethclient.DialContext(context.Background(), infura_wss)

	if err != nil {
		log.Fatal("Error while creating client!")
	}

	defer client.Close()

	contract_address := common.HexToAddress("0x93d10B0DBDFe0b541F79Eb4Af8Fe8F32ba78d7cc")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract_address},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	abiBytes, err := os.ReadFile("abi.json")
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(contractAbi)
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)

		case vLog := <-logs:
			var event LogNumber
			err := contractAbi.UnpackIntoInterface(&event, "LogNumber", vLog.Data)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(event.Number)
		}
	}

}
