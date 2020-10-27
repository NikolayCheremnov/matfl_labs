package lexinator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func scannerTesting(srsFileName string) error {
	byteModule, err := ioutil.ReadFile(srsFileName)
	if err != nil {
		return err
	}
	if len(byteModule) > MaxModuleLen {
		return errors.New("too long module")
	}

	module := string(byteModule) + "\000"
	module = strings.ReplaceAll(module, "\r", "")
	byteModule = nil // clear memory

	textPos := 0
	line := 1
	linePos := 0
	var lexType int
	var lex string
	for textPos != len(module) {
		lexType, lex, err, textPos, line, linePos = Scanner(module, textPos, line, linePos)
		if err != nil {
			fmt.Println("error: ", err, "line:", line, "line pos:", linePos, "text pos:", textPos)
		} else {
			fmt.Println("lexeme:", lex, "type:", lexType)
		}
	}

	return nil
}
