package crypto

import (
	"fmt"
	"math/rand"
	"sort"
)

const p = 257

type Secret struct {
	Index int    `json:"index"`
	Share []byte `json:"share"`
}

func makeRandParameter(t int) []int {
	params := make([]int, t)

	for i := 0; i < t; i++ {
		for {
			params[i] = rand.Intn(p)
			if params[i] != 0 {
				break
			}
		}

	}

	return params
}

type shuffle []int

func (s shuffle) Len() int {
	return len(s)
}

func (s shuffle) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s shuffle) Less(i, j int) bool {
	return rand.Int()%2 == 0
}

func makeIndexes() []int {
	var s shuffle = make([]int, 256)
	for i := 0; i < 256; i++ {
		s[i] = i + 1
	}
	sort.Sort(s)
	return s
}

// func Recover(secrets []Secret, t int) byte {
// 	if len(secrets) != t {
// 		panic("aaaa")
// 	}

// 	coffees := make([][]int, t)
// 	values := make([]int, t)
// 	for i := 0; i < t; i++ {
// 		coffee := make([]int, t)
// 		tmp := 1
// 		for j := 0; j < t; j++ {
// 			coffee[j] = tmp
// 			tmp *= secrets[i].Index % p
// 		}
// 		coffees[i] = coffee

// 		values[i] = int(secrets[i].Share)
// 	}

// 	m, s := recursion(coffees, values)

// 	return byte(s * inv(m) % p) // TODO
// }

func inv(a int) int {
	return quick_mi(a, p-2)
}

func quick_mi(a int, b int) int {
	cnt := 1
	for b > 0 {
		if b&1 != 0 {
			cnt = cnt * a % p
		}
		a = a * a % p
		b >>= 1
	}
	return cnt
}

func recursion(coffees [][]int, values []int) (int, int) {
	if len(coffees) != len(values) {
		panic("11111")
	}
	if len(coffees) == 1 {
		c, m := coffees[0][0], values[0]
		c = ((c % p) + p) % p
		m = ((m % p) + p) % p
		return c, m
	}

	t := len(coffees)
	ltrt := coffees[t-1][t-1]
	newCoffees := make([][]int, len(coffees)-1)
	newValues := make([]int, len(coffees)-1)
	for i := 0; i < t-1; i++ {
		multi := coffees[i][t-1]

		coffee := make([]int, len(coffees)-1)
		for j := 0; j < t-1; j++ {
			coffee[j] = (coffees[i][j]*ltrt - coffees[t-1][j]*multi) % p
		}
		newCoffees[i] = coffee
		newValues[i] = (values[i]*ltrt - values[t-1]*multi) % p
	}

	return recursion(newCoffees, newValues)
}

// func SplitSecret(secret byte, n, t int) []Secret {
// 	params := makeRandParameter(t)
// 	params[0] = int(secret)

// 	fmt.Println(params)
// 	indexes := makeIndexes() // TODO
// 	secrets := make([]Secret, n)
// 	cnt := 0
// 	for i := 0; i < 256 && cnt < n; i++ {
// 		tmp := calculate(params, indexes[i])
// 		if tmp == 256 {
// 			continue
// 		}
// 		secrets[cnt] = Secret{
// 			Share: byte(tmp),
// 			Index: indexes[i],
// 		}
// 		cnt++
// 	}

// 	return secrets
// }

//ax + c 
func calculate(param []int, index int) int {
	sum := 0
	tmp := 1
	for i := 0; i < len(param); i++ {
		sum = (sum + param[i]*tmp) % p
		tmp *= index
		tmp = tmp % p
	}
	return sum
}

func EnCrypto(message string, n, t int) []Secret {
	indexes := makeIndexes()
	cnt := 0
	i := 0
	secrets := make([]Secret, n)
	for i := 0; i < len(secrets); i++ {
		secrets[i] = Secret{
			Share: make([]byte, len(message)),
		}
	}
	params := makeRandParameter(t)
	for cnt < n {
		i = (i + 1) % 256
		flag := false
		for j := 0; j < len(message); j++ {
			secret := message[j]
			params[0] = int(secret)
			fmt.Println("params:", params)
			tmp := calculate(params, indexes[i])
			if tmp == 256 {
				flag = true
				break
			}
			secrets[cnt].Share[j] = byte(tmp)
		}
		if !flag {
			secrets[cnt].Index = indexes[i]
			cnt++
		}
	}
	return secrets
}

func DeCrypto(secrets []Secret, t int) []byte {
	if len(secrets) != t {
		panic("aaaa")
	}
	res := make([]byte, len(secrets[0].Share))
	for k := 0; k < len(secrets[0].Share); k++ {
		coffees := make([][]int, t)
		values := make([]int, t)
		for i := 0; i < t; i++ {
			coffee := make([]int, t)
			tmp := 1
			for j := 0; j < t; j++ {
				coffee[j] = tmp
				tmp *= secrets[i].Index % p
			}
			coffees[i] = coffee

			values[i] = int(secrets[i].Share[k])
		}
		m, s := recursion(coffees, values)
		res[k] = byte(s * inv(m) % p)
	}
	return res
}
