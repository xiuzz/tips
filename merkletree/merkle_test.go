package merkletree

import (
	"fmt"
	"testing"
)

func TestMerkle(t *testing.T) {
	var balance []string = []string{"a", "b", "c", "d", "e"}
	mt := Start(balance)
	fmt.Println(mt)
	if !mt.Verify(-1, "") {
		t.Fatal("root test error")
	}
	if mt.Verify(0, "b") {
		t.Fatal("check test error")
	}
	if !mt.Verify(0, "a") {
		t.Fatal("same check error")
	}
}
