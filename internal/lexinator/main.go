package lexinator

func main() {
	err := scannerTesting("testData/bad_test.c")
	if err != nil {
		panic(err)
	}
	err = scannerTesting("testData/bad_test.c")
	if err != nil {
		panic(err)
	}
}
