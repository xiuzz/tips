package bip39

import (
	"fmt"
	"testing"
)

func TestBip39(t *testing.T) {
	entropy := generateEntropy(128)
	mnemonic := generateMnemonic(entropy)
	fmt.Println(mnemonic)
	seed := generateSeed(mnemonic, "")
	fmt.Println(seed)
}
