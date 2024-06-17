package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"
)

const (
	Dif         = 1
	INT64_MAX   = math.MaxInt64
	MaxProbably = 255
	MinProbably = 235
)

// 创建一种名为Miner的结构体，包含miner的地址和持币数量，以及记录的币龄
type Miner struct {
	addr    []byte
	num     int64
	coinAge int64
}

// 初始化Miner的函数，默认addr为用sha256方法对字符串miner和现在的时间拼接后的字符串处理后的结果,num为0，coinAge为0
func createMiner() *Miner {
	temp := sha256.Sum256([]byte("miner" + time.Now().String()))
	miner := Miner{
		addr:    temp[:],
		num:     0,
		coinAge: 0,
	}
	return &miner
}

// 初始化Miners数组的函数，调用AddMiner函数，生成一个Miner，然后将其添加到Miners数组中
func InitMiners() []Miner {
	miner := createMiner()
	Miners := []Miner{*miner}
	return Miners
}

// 传入一个Miner和Miners数组，将miner添加到Miners数组中

func AddMiner(miner Miner, Miners *[]Miner) {
	*Miners = append(*Miners, miner)
}

// 添加矿工
func AddMiners(Miners *[]Miner) {
	var MinerNum int
	MinerNum = 4
	// fmt.Print("请输入创建矿工的数量：")
	// fmt.Scanf("%d", &MinerNum)
	for i := 0; i < MinerNum; i++ {
		AddMiner(*createMiner(), Miners)
	}
}

// 更新Miners数组函数，传入Coins数组和Miners数组，遍历Coins数组，将Coins数组中的币的矿工序号与Miners数组中的矿工序号相同的矿工的币龄加上（现在的时间-Coin的时间戳）*Coin的数量
func UpdateMiners(Coins []Coin, Miners []Miner) []Miner {
	for i := 0; i < len(Coins); i++ {
		index := Coins[i].MinerIndex
		Miners[index].coinAge += (time.Now().Unix() - Coins[i].Time) * Coins[i].Num
		Coins[i].Time = time.Now().Unix()
	}
	return Miners
}

type MinerTime struct {
	minerIndex int
	totalTime  int64
}

var start int64
var end int64

func AddMinerData(minerDatas *[]MinerTime, minerData *MinerTime) {
	*minerDatas = append(*minerDatas, *minerData)
}

// 函数名：Pos,传入Miners数组，当前难度值Dif和一个string类型变量tradeData，内设一个int变量timeCounter, 从0递增到Intmax，
// hash值为SHA256(SHA256(tradeData|timeCounter)),循环内遍历Miners数组，目标值target=Dif乘当前Miner的币龄，
// 要求hash小于target，返回满足要求的第一个Miner的序号并清空这个Miner的币龄，一旦满足要求则退出整个循环
func Pos(Miners Miner, Dif int64, tradeData string) bool {
	var timeCounter int
	var realDif int64
	realDif = int64(MinProbably)
	if realDif+Dif*Miners.coinAge > int64(MaxProbably) {
		realDif = MaxProbably
	} else {
		realDif += Dif * Miners.coinAge
	}
	target := big.NewInt(1)
	// 数据长度为8位
	// 需求：需要满足前两位为0，才能解决问题
	// 1 * 2 << (8-2) = 64
	// 0100 0000
	// 00xx xxxx
	// 32 * 8
	target.Lsh(target, uint(realDif))
	for timeCounter = 0; timeCounter < INT64_MAX; timeCounter++ {
		hash := sha256.Sum256([]byte(tradeData + string(timeCounter)))
		hash = sha256.Sum256(hash[:])
		var hashInt big.Int
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(target) == -1 {
			return true
		}
	}
	return false
}

func CorrectMiner(Miners []Miner, Dif int64, tradeData string) int {
	var minTime int64 = INT64_MAX
	var correctMiner int
	var MinerData []MinerTime
	var wait sync.WaitGroup
	var lock sync.Mutex
	wait.Add(len(Miners))
	for i := 0; i < len(Miners); i++ {
		go func(i int) {
			defer wait.Done()
			start = time.Now().UnixNano()
			//最小持币量为2才能挖矿
			time.Sleep(1000)
			if Miners[i].num >= 2 {
				success := Pos(Miners[i], Dif, tradeData)
				if success {
					end = time.Now().UnixNano()
					MinerDataDemo := MinerTime{
						minerIndex: i,
						totalTime:  end - start,
					}
					lock.Lock()
					AddMinerData(&MinerData, &MinerDataDemo)
					lock.Unlock()
				}
			}
		}(i)
	}
	wait.Wait()
	if MinerData != nil {
		fmt.Println(MinerData)
		for j := range MinerData {
			if MinerData[j].totalTime < minTime {
				minTime = MinerData[j].totalTime
				correctMiner = MinerData[j].minerIndex
			}
		}
		Miners[correctMiner].coinAge = 0
		return correctMiner
	}
	return -1
}

//这么看下来，pos并没有说那么不公平，所谓的抽奖依然要进行hash运算，只不过每个人难度值不同

func Mine(Miners []Miner, Dif int64, tradeData string, BlockChain []Block) {
	fmt.Println("开始挖矿")
	Coins := make([]Coin, 0)
	winnerIndex := CorrectMiner(Miners, Dif, tradeData)
	if winnerIndex == -1 {
		panic("挖矿失败")
	}
	fmt.Println("挖矿成功")
	fmt.Println("本轮获胜矿工:", winnerIndex)
	AddCoin(NewCoin(int64(winnerIndex), Miners), &Coins)
	GenerateBlock(winnerIndex, Miners, Coins[len(Coins)-1], tradeData, BlockChain)
	time.Sleep(5 * time.Second)
	UpdateMiners(Coins, Miners)
	PrintMiners(Miners)
}

//传入Miners数组，打印矿工数组每个矿工信息的函数

func PrintMiners(Miners []Miner) {
	for i := 0; i <= len(Miners)-1; i++ {
		fmt.Println("Miner", i, ":", hex.EncodeToString(Miners[i].addr), Miners[i].num, Miners[i].coinAge)
	}
}

func IsContinueMining(Miners []Miner, BlockChain []Block) {
	var isContinue string
	for {
		Mine(Miners, Dif, "New block", BlockChain)
		fmt.Println("是否继续挖矿?y or n")
		fmt.Scanf("%s", &isContinue)
		if isContinue == "y" {
			continue
		} else if isContinue == "n" {
			fmt.Println("挖矿结束")
			break
		} else {
			fmt.Println("输入错误")
			continue
		}
	}
}
