package main

import (
	"../../internal/parsenator"
	"fmt"
)

func main() {
	err := parsenator.GoodTesting()
	if err != nil {
		fmt.Println(err)
	}
}
