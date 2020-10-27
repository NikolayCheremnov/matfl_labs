package lexinator

import "errors"

// scanner function
// input parameter: source module text, position in text, line index, position in line,
// output: lexType, lexImage, error and changed parameters
func Scanner(sourceModule string, textPos int, line int, linePos int) (lexType int, lex string, err error, tP int, l int, lP int) {
	incPos := func() { textPos++; linePos++ }
	incLine := func() { line++; linePos = 0 }
Start:
	// preparing
	lex = "" // lexeme image

	// ignored symbols passing
	for sourceModule[textPos] == ' ' || sourceModule[textPos] == '\n' || sourceModule[textPos] == '\t' {
		if sourceModule[textPos] == '\n' {
			incLine()
		}
		incPos()
	}

	// next conditions
	// 1) is a latin letter
	if (sourceModule[textPos] >= 'a' && sourceModule[textPos] <= 'z') ||
		(sourceModule[textPos] >= 'A' && sourceModule[textPos] <= 'Z') {
		for (sourceModule[textPos] >= 'a' && sourceModule[textPos] <= 'z') ||
			(sourceModule[textPos] >= 'A' && sourceModule[textPos] <= 'Z') ||
			(sourceModule[textPos] >= '0' && sourceModule[textPos] <= '9') {
			lex += sourceModule[textPos : textPos+1] // add symbol to lexeme
			incPos()
			if len(lex) > MaxLexLen {
				return Err, lex, errors.New("too long lexeme"), textPos, line, linePos
			}
		}
		// id was retrieved => check keyword
		switch lex {
		case "main":
			return Main, lex, nil, textPos, line, linePos
		case "int":
			return Int, lex, nil, textPos, line, linePos
		case "short":
			return Short, lex, nil, textPos, line, linePos
		case "long":
			return Long, lex, nil, textPos, line, linePos
		case "bool":
			return Bool, lex, nil, textPos, line, linePos
		case "for":
			return For, lex, nil, textPos, line, linePos
		case "const":
			return Const, lex, nil, textPos, line, linePos
		case "void":
			return Void, lex, nil, textPos, line, linePos
		default:
			return Id, lex, nil, textPos, line, linePos // identifier
		}
	}

	// 2) is a digital
	if sourceModule[textPos] >= '0' && sourceModule[textPos] <= '9' {
		for sourceModule[textPos] >= '0' && sourceModule[textPos] <= '9' {
			lex += sourceModule[textPos : textPos+1]
			incPos()
			if len(lex) > MaxLexLen {
				return Err, lex, errors.New("too long lexeme"), textPos, line, linePos
			}
		}
		return IntConst, lex, nil, textPos, line, linePos
	}

	// 3) single special symbols
	if sourceModule[textPos] == ';' {
		lex = sourceModule[textPos : textPos+1]
		incPos()
		return Semicolon, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == ',' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return Comma, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == '(' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return OpeningBracket, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == ')' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return ClosingBracket, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == '{' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return OpeningBrace, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == '}' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return ClosingBrace, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == '+' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return Plus, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == '-' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return Minus, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == '*' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return Mul, lex, nil, textPos, line, linePos
	} else if sourceModule[textPos] == '%' {
		lex += sourceModule[textPos : textPos+1]
		incPos()
		return Mod, lex, nil, textPos, line, linePos
	}

	// 4) /
	if sourceModule[textPos] == '/' {
		lex = sourceModule[textPos : textPos+1]
		incPos()
		if sourceModule[textPos] == '/' { // line comments
			for sourceModule[textPos] != '\n' {
				incPos()
			}
			incPos()
			incLine()
			goto Start
		} else if sourceModule[textPos] == '*' { // multiline comment
			incPos()
			for sourceModule[textPos] != '*' || (sourceModule[textPos] == '*' && sourceModule[textPos+1] != '/') {
				if sourceModule[textPos] == '\n' {
					incPos()
					incLine()
				} else if sourceModule[textPos] == '\000' {
					return Err, "", errors.New("unclosed comment in end"), textPos, line, linePos
				} else {
					incPos()
				}
			}
			incPos() // eat *
			incPos() // eat /
			goto Start
		} else {
			return Div, lex, nil, textPos, line, linePos
		}
	}

	// <
	if sourceModule[textPos] == '<' {
		lex = sourceModule[textPos : textPos+1]
		incPos()
		if sourceModule[textPos] == '=' {
			lex += sourceModule[textPos : textPos+1]
			incPos()
			return LessEqu, lex, nil, textPos, line, linePos
		} else {
			return less, lex, nil, textPos, line, linePos
		}
	}

	// >
	if sourceModule[textPos] == '>' {
		lex = sourceModule[textPos : textPos+1]
		incPos()
		if sourceModule[textPos] == '=' {
			lex += sourceModule[textPos : textPos+1]
			incPos()
			return MoreEqu, lex, nil, textPos, line, linePos
		} else {
			return More, lex, nil, textPos, line, linePos
		}
	}

	// =
	if sourceModule[textPos] == '=' {
		lex = sourceModule[textPos : textPos+1]
		incPos()
		if sourceModule[textPos] == '=' {
			lex += sourceModule[textPos : textPos+1]
			incPos()
			return Equ, lex, nil, textPos, line, linePos
		} else {
			return Assignment, lex, nil, textPos, line, linePos
		}
	}

	// module end
	if sourceModule[textPos] == '\000' {
		return End, "null_char", nil, textPos + 1, line, linePos
	}

	// error on input
	incPos()
	return Err, "", errors.New("invalid input symbol"), textPos, line, linePos
}
