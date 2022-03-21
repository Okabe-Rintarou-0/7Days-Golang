package main

import (
	"fmt"
	"net/url"
	"testing"
)

func Test(t *testing.T) {
	slc := []string{"a", "b", "c", "d"}
	// Insert
	rear := append([]string{"Inserted"}, slc[2:]...)
	slc = append(slc[:1], rear...)
	fmt.Println(slc)

	// Delete conditionally
	for i := 0; i < len(slc); {
		if slc[i] == "d" {
			slc = append(slc[:i], slc[i+1:]...)
		} else {
			i++
		}
	}
	fmt.Println(slc)
	if false {
		t.Error()
	}

	fmt.Println(url.QueryEscape("1"))
}
