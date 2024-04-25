package crypto

import (
	"fmt"
	"testing"
)

func TestCrypto(t *testing.T) {
	New(10)
	arr := make([]byte, 3);
	for i := 0; i < 3; i++ {
		arr[i] = EnCrypto(i+1)
	}
	dec := make([][]int, 3)
	for i := 0; i < len(dec); i++ {
		dec[i] = []int{((i+1)*(i+1)), int(i+1), 1, int(arr[i])}
	}
	ans := DeCrypto(dec)
	if ans != c {
		fmt.Println(ans)
		t.Fatal("err")
	}
}
