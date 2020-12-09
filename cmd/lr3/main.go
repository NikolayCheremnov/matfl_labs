package main

import (
	"../../internal/lexinator"
	"fmt"
)

func main() {
	// test good and bad files
	fmt.Println("good test:")
	err := lexinator.ScannerTesting("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/lexinator/testData/good_test.c",
		"")
	if err != nil {
		panic(err)
	}
	fmt.Println("bad test:")
	err = lexinator.ScannerTesting("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/lexinator/testData/bad_test.c", "errors.err")
	if err != nil {
		panic(err)
	}
}
