package parsenator

import (
	"io"
	"os"
)

func GoodTesting() (err error) {
	// writers preparing
	var sw, aw io.Writer
	sw, err = os.Create("lexinatorErrors.err")
	if err != nil {
		return err
	}
	aw, err = os.Create("parsenatorErrors.err")
	if err != nil {
		return err
	}

	// analyzer preparation
	A, err := Preparing("D:/CherepNick/ASTU/4_course/7_semester/MATFL/labs/internal/lexinator/testData/good_test.c", sw, aw)
	if err != nil {
		return nil
	}

	// testing
	err = A.GlobalDescriptions()
	if err != nil {
		return err
	}
	return nil
}
