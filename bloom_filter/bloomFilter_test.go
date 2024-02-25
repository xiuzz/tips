package bloom_filter

import "testing"

func TestBloomFilter(t *testing.T) {
	metaData := []string{"test", "scientist", "computer", "block", "chain"}
	bf := New(len(metaData), 0.1)
	for i := 0; i < len(metaData); i++ {
		bf.Insert([]byte(metaData[i]))
	}
	if !bf.Query([]byte("test")) {
		t.Fail()
	} 
	if bf.Query([]byte("ttttt")) {
		t.Fail()
	}
}