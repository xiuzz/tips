package bip39

import (
	"fmt"
	"testing"
)

func TestBip39(t *testing.T) {
	entropy := generateEntropy(128)
	mnemonic := generateMnemonic(entropy)
	var format string
	for i := 0; i < len(mnemonic); i++ {
		if i != 0 {
			format += " "
		}
		format += mnemonic[i]
	}
	fmt.Println(format)
	t.Fail()
}
