package main

import "../../internal/lexinator"

func main() {
	err := lexinator.ScannerTesting("D:/CherepNick/ASTU/4_course/7_semester/MATFL/lr3/internal/lexinator/testData/good_test.c",
		"")
	if err != nil {
		panic(err)
	}
	/*/err = lexinator.ScannerTesting("D:/CherepNick/ASTU/4_course/7_semester/MATFL/lr3/internal/lexinator/testData/bad_test.c")
	if err != nil {
		panic(err)
	}
	/*/
}
