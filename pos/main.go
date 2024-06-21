package main

import (
	. "pos/core"
	"time"
)

// 创建币池数组Coins
var Coins []Coin

// 调用InitBlockChain函数，生成一个区块数组
var BlockChain []Block

// 创建矿工数组Miners
var Miners []Miner

func main() {
	//默认难度值dif为1
	//var Dif int64 = 1
	//创建矿工数组Miners
	//var Miners []Miner
	Miners = InitMiners()
	//添加矿工
	AddMiners(&Miners)

	//创建币池数组Coins
	//var Coins []Coin
	//给矿工数组中的矿工添加币
	Coins = InitCoins(Miners)
	for i := 0; i < len(Miners); i++ {
		AddCoin(NewCoin(int64(i), Miners), &Coins)
	}
	//调用InitBlockChain函数，生成一个区块数组
	BlockChain = InitBlockChain(Miners, Coins)
	//时间延迟，给出币龄
	time.Sleep(5 * time.Second)
	UpdateMiners(Coins, Miners)
	PrintMiners(Miners)

	//挖矿
	IsContinueMining(Miners, BlockChain)
}
