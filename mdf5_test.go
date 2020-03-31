package main

import (
	"fmt"
	"testing"
)

func TestMdf5(t *testing.T) {
	var s = "88888888"
	res := mdf5(s)
	fmt.Println(res)
}
