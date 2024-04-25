package crypto

const (
	P = 257
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

func EnCrypto(index int) byte {
	return byte((A*index*index + B*index + int(c)) % P)
}

func DeCrypto(arr [][]int) byte {
	xa := (arr[0][2]*arr[1][0]%P - arr[1][2]*arr[0][0]%P + P) % P
	xb := (arr[1][1]*arr[2][0]%P - arr[2][1]*arr[1][0]%P + P) % P
	xc := (arr[1][2]*arr[2][0]%P - arr[2][2]*arr[1][0]%P + P) % P
	xd := (arr[0][1]*arr[1][0]%P - arr[1][1]*arr[0][0]%P + P) % P
	xz1 := (arr[0][3]*arr[1][0]%P - arr[1][3]*arr[0][0]%P + P) % P
	xz2 := (arr[1][3]*arr[2][0]%P - arr[2][3]*arr[1][0]%P + P) % P
	y1 := ((xa*xb)%P - (xc*xd)%P + P) % P
	y2 := (xz1*xb%P - xz2*xd%P + P) % P
	ans := y2 / y1
	return byte(ans)
}

func New(val byte) {
	c = val
}
