package main

import (
	"fmt"
	"pow/pow"
	"pow/web"
	"time"
)

func main() {
	var count int
	fmt.Printf("请输入初始矿工数量：")
	fmt.Scanf("%d", &count)
	time.Sleep(10000)
	fmt.Printf("开始挖矿")
	//新建区块链网络
	work := pow.NewBlockChainNetWork(pow.BlockchainConfig{
		//矿工数量
		MinerCount: count,
		//平均出块时间
		OutBlockTime: 10,
		//初始难道值
		InitialDifficulty: 20,
		//每多少个区块修改一次难度值
		ModifyDifficultyBlockNumber: 10,
		//每次记账奖励
		BookkeepingIncentives: 20,
	})
	//运行区块链网络
	work.RunBlockChainNetWork()
	//启动web服务
	web.RunRouter(work)
}
