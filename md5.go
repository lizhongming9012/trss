package main

import (
	"crypto/md5"
	"fmt"
)

func mdf5(sour string) (res string) {
	data := []byte(sour)
	has := md5.Sum(data)
	res = fmt.Sprintf("%x", has)
	return
}
