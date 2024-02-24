package bip39

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

// The allowed size of ENT is 128-256 bits
func generateEntropy(bitSize int) []byte {
	if bitSize < 128 || bitSize > 256 || bitSize%32 != 0 {
		panic("bad arg!")
	}
	entropy := make([]byte, bitSize/8)
	rand.Read(entropy)
	return entropy
}

func generateMnemonic(entropy []byte) string {
	entLen := len(entropy) * 8
	csLen := entLen / 32
	msLen := (entLen + csLen) / 11

	hasher := sha256.New()
	hasher.Write(entropy)
	hash := hasher.Sum(nil)
	csByteGet := hash[0]
	bigTwo := big.NewInt(2)
	bigOne := big.NewInt(1)
	dataBigInt := new(big.Int).SetBytes(entropy)
	for i := 0; i < csLen; i++ {
		dataBigInt.Mul(dataBigInt, bigTwo)
		if csByteGet&(1<<(7-i)) != 0 {
			dataBigInt.Or(dataBigInt, bigOne)
		}
	}
	// fmt.Println(dataBigInt)
	mnemonic := make([]string, msLen)
	index := big.NewInt(0)
	//2048
	save := big.NewInt(2047)
	mod := big.NewInt(2048)
	for i := msLen - 1; i >= 0; i-- {
		index.And(save, dataBigInt)
		mnemonic[i] = English[index.Int64()]
		dataBigInt.Div(dataBigInt, mod)
	}

	var result string
	for i := 0; i < msLen; i++ {
		if i != 0 {
			result += " "
		}
		result += mnemonic[i]
	}
	return result
}

func generateSeed(mnemonic string, passphrase string) []byte {
	return pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"+passphrase), 2048, 64, sha512.New)
}
