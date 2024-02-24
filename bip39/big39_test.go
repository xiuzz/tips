package bip39

import (
	"fmt"
	"testing"
)

func TestBip39(t *testing.T) {
	entropy := generateEntropy(128)
	mnemonic := generateMnemonic(entropy)
	for i := 0; i < len(mnemonic); i++ {
		fmt.Println(mnemonic[i])
	}
		
}
