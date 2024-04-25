package crypto

import (
	"fmt"
	"testing"
)

func TestCrypto1(t *testing.T) {
	New(111)
	flag := false
	for i := 0; i < 256; i++ {
		x := make([]int, 3)
		y := make([]int, 3)
		for i := 0; i < 3; i++ {
			x[i] = i + 1
			y[i] = int(EnCrypto(i + 1))
			// fmt.Println(arr[i])
		}

		ans := DeCrypto(x, y)
		if ans != c {
			fmt.Println(i, ans)
			flag = true
		}
	}
	if flag {
		t.Fatal("err")
	}
}

func TestCrypto2(t *testing.T) {
	for j := 0; j < P; j++ {
		New(byte(j))
		flag := false
		x := make([]int, 3)
		y := make([]int, 3)

		for i := 0; i < 3; i++ {
			x[i] = i + 1
			y[i] = int(EnCrypto(i + 1))
			// fmt.Println(arr[i])
		}

		ans := DeCrypto(x, y)
		if ans != c {
			fmt.Println(j, ans)
			flag = true
		}
		if flag {
			t.Fatal("err")
		}
	}
}
