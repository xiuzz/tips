package crypto

const (
	P = 251
	A = 101
	B = 39
)

// f(x) = Ax^2 + Bx + C

/*
*

	a1 b1 c1 z1
	a2 b2 c2 z2
	a3 b3 c3 z3
*/
var c byte

func quick_mi(a int, b int) int {
	ans := 1
	for b != 0 {
		if b&1 == 1 {
			ans = ans * a % P
		}
		a = a * a % P
		b >>= 1
	}
	return ans
}

func inv(a int) int {
	return quick_mi(a, P-2)
}
func EnCrypto(index int) byte {
	return byte((A*index*index + B*index + int(c)) % P)
}

func DeCrypto(x []int, y []int) byte {
	b := (((x[1]*x[1]-x[2]*x[2])*(y[0]-y[1]) - (x[0]*x[0]-x[1]*x[1])*(y[1]-y[2])) * inv((x[1]*x[1]-x[2]*x[2])*(x[0]-x[1])-(x[0]*x[0]-x[1]*x[1])*(x[1]-x[2])) % P + P) % P
	a := ((y[0] - y[1] - b*(x[0]-x[1])) * inv(x[0]*x[0]-x[1]*x[1]) % P + P) % P
	c := ((y[0] - a*x[0]*x[0] - b*x[0]) % P + P)% P
	return byte(c)
}

func New(val byte) {
	c = val
}
