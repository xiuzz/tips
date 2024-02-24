package main

import (
	"fmt"
	"math/big"
)

func main() {
	x := []byte{byte(0), byte(32)}
	bigInt := new(big.Int).SetBytes(x)
	fmt.Println(bigInt)
}
