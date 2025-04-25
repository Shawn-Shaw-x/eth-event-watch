package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func TestGetLogs(t *testing.T) {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 获取环境变量
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	apiKey := os.Getenv("ETHEREUM_API_KEY")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	eventSignature := os.Getenv("EVENT_SIGNATURE")

	if rpcURL == "" || apiKey == "" || contractAddress == "" {
		log.Fatal(" ❌ Missing required environment variables")
	}

	// 连接以太坊节点（使用你的 RPC）
	client, err := ethclient.Dial(rpcURL + apiKey)
	if err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()

	// 构造筛选参数
	var addresses []common.Address
	addresses = append(addresses, common.HexToAddress(contractAddress))
	var eventSignatures [][]common.Hash
	eventSignatures = append(eventSignatures, []common.Hash{common.HexToHash(eventSignature)})
	query := ethereum.FilterQuery{
		BlockHash: nil,
		FromBlock: big.NewInt(22341322),
		ToBlock:   big.NewInt(22341324),
		Addresses: addresses,
		Topics:    eventSignatures,
	}
	// 发起调用
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("❌ 获取 log失败: %v", err)
	}
	// 遍历日志
	for _, vLog := range logs {
		fmt.Println("📦 日志地址:", vLog.Address.Hex())
		fmt.Println("📝 topics:")
		for i, topic := range vLog.Topics {
			fmt.Printf("  - topic[%d]: %s\n", i, topic.Hex())
		}
		fmt.Println("📨 data:", hex.EncodeToString(vLog.Data))
	}

}
