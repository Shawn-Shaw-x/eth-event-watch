package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func TestGetTransactionReceipt(t *testing.T) {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 获取环境变量
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	apiKey := os.Getenv("ETHEREUM_API_KEY")
	eventSignature := os.Getenv("EVENT_SIGNATURE")

	if rpcURL == "" || apiKey == "" {
		log.Fatal(" ❌ Missing required environment variables")
	}

	// 连接以太坊节点（使用你的 RPC）
	client, err := ethclient.Dial(rpcURL + apiKey)
	if err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()

	// 获取最新区块号
	latestBlockNum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取最新区块号失败: %v", err)
	}
	// 根据区块号获取区块
	latestBlock, err := client.BlockByNumber(context.Background(), big.NewInt(int64(latestBlockNum)))
	if err != nil {
		log.Fatalf("❌ 获取最新区块内容失败: %v", err)
	}

	// todo 筛选交易
	// 替换为你想要监听的交易哈希
	txHash := latestBlock.Transactions()[10].Hash()

	// 获取交易回执（包含事件 logs）
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatalf("❌ 获取回执失败: %v", err)
	}

	// 遍历日志
	for _, vLog := range receipt.Logs {
		if (vLog.Topics[0].Hex()) != eventSignature {
			// 跳过非 Transfer 事件
			continue
		}
		fmt.Println("📦 日志地址:", vLog.Address.Hex())
		fmt.Println("📝 topics:")
		for i, topic := range vLog.Topics {
			fmt.Printf("  - topic[%d]: %s\n", i, topic.Hex())
		}
		fmt.Println("📨 data:", hex.EncodeToString(vLog.Data))
	}
}
