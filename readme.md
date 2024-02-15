# 随手笔记和代码
写一些小知识点
## markle tree
默克尔树又叫hash树，本质就是一个完全二叉树,因此实现起来还是比较容易
```go
package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// 因为完全二叉树的性质，且在区块中一旦确定交易数量，后续是无法改变树上节点数的，因此用固长的数组反而更加容易？

// 这个 balance 可以是任何类型的账单数据，或者二进制，json
type node struct {
	balance string // 这里简单用string表示
	curHash string
	seedId  int //if seed id == 0 not seed
}

type MerkleTree struct {
	RootHash string
	treeCore []node
}

// 网上随便找的sha256方法
func GetSHA256HashCode(message []byte) string {
	//方法一：
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode

	//方法二：
	//bytes2:=sha256.Sum256(message)//计算哈希值，返回一个长度为32的数组
	//hashCode2:=hex.EncodeToString(bytes2[:])//将数组转换成切片，转换成16进制，返回字符串
	//return hashCode2
}

var cur int

func createTree(balances []string, index int, treeCore []node) string {
	if index >= len(treeCore) {
		return ""
	}
	treeCore[index].seedId = -1
	l := createTree(balances, index<<1, treeCore)
	r := createTree(balances, index<<1|1, treeCore)

	if index*2 >= len(treeCore) && cur < len(balances) {
		// fmt.Println(index, len(treeCore), cur)
		treeCore[index].seedId = cur
		treeCore[index].balance = balances[cur]
		treeCore[index].curHash = GetSHA256HashCode([]byte(balances[cur]))
		cur++
	} else {
		if r == "" {
			r = l
		}
		metaHash := l + r
		treeCore[index].curHash = GetSHA256HashCode([]byte(metaHash))
	}
	return treeCore[index].curHash
}

func bo(n int) int {
	begin := 1
	cnt := n
	for n > begin {
		cnt += begin
		begin *= 2
	}
	return cnt
}

func (mt *MerkleTree) VerifyForRootHash() bool {
	return mt.Verify(-1, "")
}

func (mt *MerkleTree) Verify(seedId int, balance string) bool {
	if seedId >= cur {
		panic("illegal index")
	}
	return mt.query(seedId, 1, balance) == mt.RootHash
}
func (mt *MerkleTree) query(seedId int, index int, balance string) string {
	if index >= len(mt.treeCore) {
		return ""
	}

	l := mt.query(seedId, index<<1, balance)
	r := mt.query(seedId, index<<1|1, balance)
	if mt.treeCore[index].seedId >= 0 {
		if mt.treeCore[index].seedId == seedId {
			return GetSHA256HashCode([]byte(balance))
		} else {
			return mt.treeCore[index].curHash
		}
	} else {
		if r == "" {
			r = l
		}
		metaHash := l + r
		hash := GetSHA256HashCode([]byte(metaHash))
		if hash != mt.treeCore[index].curHash {
			fmt.Println(hash, mt.treeCore[index].curHash)
		}
		return hash
	}
}

func Start(balances []string) *MerkleTree { //初始化位置 给定账单长度
	if len(balances) == 0 {
		panic("nil!")
	}
	if len(balances)&1 != 0 {
		balances = append(balances, balances[len(balances)-1])
	}
	len := bo(len(balances))
	treeCore := make([]node, len+1)
	createTree(balances, 1, treeCore)

	merkle := MerkleTree{
		RootHash: treeCore[1].curHash,
		treeCore: treeCore,
	}
	return &merkle
}

```
简单做了下test
![alt text](picture/image-5.png)
![alt text](picture/image-6.png)
![alt text](picture/image.png)
![alt text](picture/image-3.png)
![alt text](picture/image-1.png)