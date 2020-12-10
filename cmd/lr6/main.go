package main

import (
	"../../internal/parsenator"
	"fmt"
	"log"
)

func main() {
	buf := ""
	// 1. good test
	err := parsenator.Testing("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/parsenator/testData/good_test.c")
	if err != nil {
		log.Println(err)
	}
	fmt.Scanln(&buf)
	// 2. expression error
	err = parsenator.Testing("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/parsenator/testData/expression_err.c")
	if err != nil {
		log.Println(err)
	}
	fmt.Scanln(&buf)
	// 3. id err
	err = parsenator.Testing("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/parsenator/testData/id_err.c")
	if err != nil {
		log.Println(err)
	}
	fmt.Scanln(&buf)
	// 4. invalid lexeme
	err = parsenator.Testing("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/parsenator/testData/invalid_lex_err.c")
	if err != nil {
		log.Println(err)
	}
	fmt.Scanln(&buf)
	// 5. type err
	err = parsenator.Testing("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/parsenator/testData/type_err.c")
	if err != nil {
		log.Println(err)
	}
	fmt.Scanln(&buf)
	// 6. infrastructure err
	err = parsenator.Testing("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/parsenator/testData/infrastructure_err.c")
	if err != nil {
		log.Println(err)
	}
}
