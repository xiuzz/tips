package bloom_filter

import (
	"crypto/sha256"
	"math"
	"math/big"
)

type BloomFilter struct {
	bitSet  []byte
	hashLen int
}

const Itor = 114

var bitLen *big.Int

func New(n int, expect float64) *BloomFilter {
	m := int(-float64(n) * math.Log(expect) / math.Pow((math.Ln2), 2))
	k := int(float64(m) / float64(n) * math.Ln2)
	bloomFilter := &BloomFilter{
		bitSet:  make([]byte, (m+7)/8),
		hashLen: k,
	}
	bitLen = big.NewInt(int64(m))
	return bloomFilter
}

func (bloomFilter *BloomFilter) Insert(element []byte) {
	hasher := sha256.New()
	temp := make([]byte, len(element))
	copy(temp, element)
	for i := 0; i < bloomFilter.hashLen; i++ {
		temp = append(temp, byte(Itor))
		hasher.Write(temp)
		hash := hasher.Sum(nil)
		dataBigInt := new(big.Int).SetBytes(hash)
		dataBigInt.Mod(dataBigInt, bitLen)
		bloomFilter.bitSet[dataBigInt.Int64()/8] |= (1 << byte(dataBigInt.Int64()%8))
	}
}

func (bloomFilter *BloomFilter) Query(element []byte) bool{
	hasher := sha256.New()
	temp := make([]byte, len(element))
	copy(temp, element)
	cnt := 0
	for i := 0; i < bloomFilter.hashLen; i++ {
		temp = append(temp, byte(Itor))
		hasher.Write(temp)
		hash := hasher.Sum(nil)
		dataBigInt := new(big.Int).SetBytes(hash)
		dataBigInt.Mod(dataBigInt, bitLen)
		if bloomFilter.bitSet[dataBigInt.Int64()/8]&(1<<byte(dataBigInt.Int64()%8)) != 0 {
			cnt++
		}
	}
	return cnt == bloomFilter.hashLen
}
