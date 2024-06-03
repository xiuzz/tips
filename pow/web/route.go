package web

import (
	"pow/pow"
	"github.com/gin-gonic/gin"
)

func RunRouter(blockchain *pow.Blockchain) {
	r := gin.Default()
	r.GET("/addMiner", addMiner(blockchain))
	r.GET("/getBlockChainInfo", getBlockChainInfo(blockchain))
	r.Run()
}

// 增加矿工
func addMiner(blockchain *pow.Blockchain) gin.HandlerFunc {
	return func(c *gin.Context) {
		blockchain.IncreaseMiner()
		c.JSON(200, gin.H{
			"message": "增加成功",
		})
	}
}

// 打印挖矿信息
func getBlockChainInfo(blockchain *pow.Blockchain) gin.HandlerFunc {
	return func(c *gin.Context) {
		blocks, miners := blockchain.GetBlockInfo()
		c.JSON(200, gin.H{
			"blocks": blocks,
			"miners": miners,
		})
	}
}
