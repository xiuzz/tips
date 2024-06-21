package crypto

import (
	"fmt"
	"testing"
)

func TestSecretshare(t *testing.T) {
	secrets := EnCrypto("qfoqofekjxcnqowhfvqhfoqwnfcoaxcahdowqfyhpopwqhdwqpfnnwfcuncqpn", 10, 2)
	fmt.Println(secrets)
	value := DeCrypto(secrets[:2], 2)
	fmt.Println(string(value))
}
