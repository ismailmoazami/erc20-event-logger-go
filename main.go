package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"eventlogger/database"
	"eventlogger/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

func main() {

	// --------------- API Logic -----------------
	router := gin.Default()
	router.GET("/records", getTransferRecords)
	database.InitDB()

	go func() {
		router.Run("localhost:8080")
	}()
	// --------------------------------------------

	infura_wss := os.Getenv("INFURA_WSS")
	token_address := os.Getenv("ADDRESS")
	client, err := ethclient.DialContext(context.Background(), infura_wss)

	if err != nil {
		log.Fatal("Error while creating client!")
	}

	defer client.Close()

	contract_address := common.HexToAddress(token_address)

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
			var transferEvent struct {
				Value *big.Int
			}

			From := vLog.Topics[1].Hex()
			To := vLog.Topics[2].Hex()

			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				fmt.Println(err)
			}

			var transferRecord models.TransferRecord
			transferRecord.From = From
			transferRecord.To = To
			transferRecord.Value = transferEvent.Value.String()

			database.DB.Create(&transferRecord)
		}
	}

}

func getTransferRecords(c *gin.Context) {
	var transferRecords []models.TransferRecord

	if err := database.DB.Find(&transferRecords).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not get data"})
		return
	}
	c.JSON(http.StatusOK, transferRecords)
}
