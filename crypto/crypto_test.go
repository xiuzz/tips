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
		cnt := 0
		for i := 0; i < P; i++ {
			tmp, flag := EnCrypto(i + 1)
			if flag {
				continue
			}
			x[cnt] = i + 1
			y[cnt] = int(tmp)
			cnt++
			if cnt == 3 {
				break
			}
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

		cnt := 0
		for i := 0; i < P; i++ {
			tmp, flag := EnCrypto(i + 1)
			if flag {
				continue
			}
			x[cnt] = i + 1
			y[cnt] = int(tmp)
			cnt++
			if cnt == 3 {
				break
			}
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
