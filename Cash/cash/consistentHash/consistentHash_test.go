package consistentHash

import (
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	nodes := []string{"2", "4", "6"}
	ch := New(3, nodes, func(bytes []byte) uint32 {
		n, _ := strconv.Atoi(string(bytes))
		return uint32(n)
	})
	testMap := map[string]string{
		"3":  "4",
		"4":  "4",
		"13": "4",
		"15": "6",
		"21": "2",
		"27": "2",
		"25": "6",
	}

	for k, v := range testMap {
		if ch.Get(k) != v {
			t.Errorf("Error occurs here k = %s! expected %s, but got %s\n", k, v, ch.Get(k))
		}
	}

	ch.AddNode("1")

	testMap2 := map[string]string{
		"3":  "4",
		"23": "4",
		"27": "1",
		"2":  "2",
	}

	for k, v := range testMap2 {
		if ch.Get(k) != v {
			t.Errorf("Error occurs here k = %s! expected %s, but got %s\n", k, v, ch.Get(k))
		}
	}

	ch.DeleteNode("6")

	testMap3 := map[string]string{
		"2":  "2",
		"25": "1",
		"27": "1",
		"1":  "1",
		"15": "1",
		"5":  "1",
	}

	for k, v := range testMap3 {
		if ch.Get(k) != v {
			t.Errorf("Error occurs here k = %s! expected %s, but got %s\n", k, v, ch.Get(k))
		}
	}
}
