package pow

import "fmt"

type Miner struct {
	//矿工ID
	Id int64 `json:"id"`
	//矿工账户余额
	Balance uint `json:"balance"`
	//当前矿工正在挖的区块
	blockchain *Blockchain
	// 用于通知 当接收到新区块的时候 不应该从原有的链继续往后挖
	waitForSignal chan interface{} `json:"-"`
}

// 挖矿逻辑
func (m Miner) run() {
	count := 0
	//死循环
	for ; ; count++ {
		//根据全局信息组装去了
		blockWithoutProof := m.blockchain.assembleNewBlock(m.Id, []byte(fmt.Sprintf("模拟区块数据:%d:%d", m.Id, count)))
		block, finish := blockWithoutProof.Mine(m.waitForSignal)
		if !finish {
			//如果不满足条件则计数器增加继续计算hash并判断
			continue
		} else {
			//如果条件满足则增加区块
			m.blockchain.AddBlock(block, m.waitForSignal)
		}
	}
}
