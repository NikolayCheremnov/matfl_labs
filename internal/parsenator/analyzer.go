package parsenator

import (
	"../../internal/lexinator"
	"errors"
	"fmt"
	"io"
)

// analyzer structure
type Analyzer struct {
	scanner    lexinator.Scanner
	isPrepared bool
	// the output stream of error messages
	writer io.Writer
}

// preparing the analyzer
func Preparing(srsFileName string, scannerErrWriter io.Writer, analyzerErrWriter io.Writer) (Analyzer, error) {
	A := Analyzer{writer: analyzerErrWriter}
	scanner, err := lexinator.ScannerInitializing(srsFileName, scannerErrWriter)
	if err != nil {
		return A, err
	}
	A.scanner = scanner
	A.isPrepared = true
	return A, nil
}

// the results of the error message
func (A *Analyzer) printPanicError(msg string) {
	textPos, line, linePos := A.scanner.StorePosValues()
	_, err := fmt.Fprintf(A.writer, "error: %s position: %d line: %d line position: %d\n", msg, textPos, line, linePos)
	if err != nil {
		panic(err)
	}
	// maybe temporarily
	panic(errors.New("completed with an error. see the error logs")) // critical completion
}

func (A *Analyzer) printError(msg string) {
	textPos, line, linePos := A.scanner.StorePosValues()
	_, err := fmt.Fprintf(A.writer, "error: %s position: %d line: %d line position: %d\n", msg, textPos, line, linePos)
	if err != nil {
		panic(err)
	}
}

// handlers for the nonterminals

// done
// axiom: <глобальные описания> -> <описание процедуры>|<описание>|;|<main>|e U <глобальные описания>|e
// input: source module file, stream for write errors
func (A *Analyzer) GlobalDescriptions() error {
	if !A.isPrepared {
		return errors.New("can't start the analysis: the analyzer is not prepared")
	}

	// support function
	getBranch := func() int {
		textPos, line, linePos := A.scanner.StorePosValues()
		lexType, _ := A.scanner.Scan()
		var branch int
		if lexType == lexinator.Void {
			branch = 1
		} else if lexType == lexinator.Long ||
			lexType == lexinator.Short ||
			lexType == lexinator.Int ||
			lexType == lexinator.Bool ||
			lexType == lexinator.Const {
			branch = 2
		} else if lexType == lexinator.Semicolon {
			branch = 3
		} else if lexType == lexinator.End {
			branch = 4
		} else {
			branch = 0 // error branch
		}
		A.scanner.RestorePosValues(textPos, line, linePos)
		return branch
	}

	isEnd := false
	for !isEnd {
		switch getBranch() {
		case 1:
			A.procedureDescription()
			break
		case 2:
			A.description()
			break
		case 3:
			A.scanner.Scan()
			break
		case 4:
			A.scanner.Scan()
			isEnd = true
			break
		case 0:
			A.printPanicError("invalid program infrastructure")
		}
	}
	A.printError("there are no syntax level errors")
	return nil
}

// <описание параметров> -> long int | short int | int | bool U идентификатор U , <описание параметров> | e
func (A *Analyzer) parameterDescription() {
	var textPos, line, linePos int
	var lexType int
	var lex string
	isFirst := true
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Int && lexType != lexinator.Long && lexType != lexinator.Short && lexType != lexinator.Bool { // int long short bool
			A.printPanicError("'" + lex + "'" + "does not name a type")
		}
		if lexType == lexinator.Long || lexType == lexinator.Short {
			lexType, lex = A.scanner.Scan()
			if lexType != lexinator.Int {
				A.printPanicError("'" + lex + "'" + "does not name a type")
			}
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Id {
			A.printPanicError("'" + lex + "' is not an identifier")
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <параметры> -> идентификатор | константа U , <параметры> | e
func (A *Analyzer) parameters() {
	isFirst := true
	var textPos, line, linePos int
	var lexType int
	var lex string
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Id && lexType != lexinator.IntConst {
			A.printPanicError("'" + lex + "' is not an identifier or constant")
		}
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <эл.выр.> -> (<выражение>) | идентификатор | константа
func (A *Analyzer) simplestExpr() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id && lexType != lexinator.IntConst && lexType != lexinator.OpeningBracket {
		A.printPanicError("'" + lex + "' not allowed in the expression")
	}
	if lexType == lexinator.OpeningBracket {
		A.expression()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {

		}
	}
}

// <множитель> -> <эл.выр.> U e | * U / U % <эл.выр.>
func (A *Analyzer) multiplier() {
	var textPos, line, linePos int
	var lexType int
	isFirst := true
	for isFirst || lexType == lexinator.Mul || lexType == lexinator.Div || lexType == lexinator.Mod {
		isFirst = false
		A.simplestExpr()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <процедура> -> идентификатор ( ) | идентификатор ( <параметры> )
func (A *Analyzer) procedure() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.OpeningBracket {
		A.printPanicError("invalid lexeme '" + lex + "', expected '('")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBracket {
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.parameters()
		lexType, lex = A.scanner.Scan()
	}
	if lexType != lexinator.ClosingBracket {
		A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
	}
}

// <слагаемое> -> <множитель> U +- | e
func (A *Analyzer) summand() {
	var textPos, line, linePos int
	var lexType int
	isFirst := true
	for isFirst || lexType == lexinator.Plus || lexType == lexinator.Minus {
		isFirst = false
		A.multiplier()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <выражение> -> + | - | e U <слагаемое> + +- == <= >= < > <слагаемое> | e
func (A *Analyzer) expression() {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, _ := A.scanner.Scan()
	if lexType != lexinator.Plus && lexType != lexinator.Minus {
		A.scanner.RestorePosValues(textPos, line, linePos)
	}
	isFirst := true
	for isFirst || lexType == lexinator.Plus || lexType == lexinator.Minus ||
		lexType == lexinator.Equ || lexType == lexinator.LessEqu || lexType == lexinator.MoreEqu ||
		lexType == lexinator.Less || lexType == lexinator.More {
		isFirst = false
		A.summand()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, _ = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <оператор for> -> e |
// long int | short int | int | bool | e U идентификатор = выражение
func (A *Analyzer) forOperator() {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.Semicolon {
		A.scanner.RestorePosValues(textPos, line, linePos)
	} else {
		if lexType != lexinator.Long && lexType != lexinator.Short &&
			lexType != lexinator.Int && lexType != lexinator.Bool && lexType != lexinator.Id {
			A.printPanicError("'" + lex + "' is not an identifier or type")
		}
		if lexType == lexinator.Long || lexType == lexinator.Short {
			lexType, lex = A.scanner.Scan()
			if lexType != lexinator.Int {
				A.printPanicError("'" + lex + "'" + "does not name a type")
			}
		}
		if lexType != lexinator.Id {
			lexType, lex = A.scanner.Scan()
			if lexType != lexinator.Id {
				A.printPanicError("'" + lex + "' is not an identifier")
			}
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Assignment {
			A.printPanicError("invalid lexeme '" + lex + "', expected '='")
		}
		A.expression()
	}
}

// <присваивание> -> идентификатор = <выражение>
func (A *Analyzer) assigment() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Assignment {
		A.printPanicError("invalid lexeme '" + lex + "', expected '='")
	}
	A.expression()
}

// <переменная> -> идентификатор U e | = <выражение>
func (A *Analyzer) variable() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Id {
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType == lexinator.Assignment {
		A.expression()
	} else {
		A.scanner.RestorePosValues(textPos, line, linePos)
	}
}

// <for> -> for ( <оператор for> ; U <выражение> | e U ; U <присваивание> | e U ) <оператор>
func (A *Analyzer) _for() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.For {
		A.printPanicError("invalid lexeme '" + lex + "', expected 'for'")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.OpeningBracket {
		A.printPanicError("invalid lexeme '" + lex + "', expected '('")
	}
	A.forOperator()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Semicolon {
		A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Semicolon {
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.expression()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Semicolon {
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
	}
	textPos, line, linePos = A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBracket {
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.assigment()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
	}
	A.operator()
}

// <константы> -> const U long int | short int | int | bool U e | <присваивание> U e | , <присваивание>
func (A *Analyzer) constants() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Const {
		A.printPanicError("invalid lexeme '" + lex + "', expected 'const'")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Long && lexType != lexinator.Short &&
		lexType != lexinator.Int && lexType != lexinator.Bool {
		A.printPanicError("'" + lex + "'" + "does not name a type")
	}
	if lexType == lexinator.Long || lexType == lexinator.Short {
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Int {
			A.printPanicError("'" + lex + "'" + "does not name a type")
		}
	}
	var textPos, line, linePos int
	isFirst := true
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		A.assigment()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <переменные> -> long int | short int | int | bool U e | <присваивание> U e | , <присваивание>
func (A *Analyzer) variables() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Long && lexType != lexinator.Short &&
		lexType != lexinator.Int && lexType != lexinator.Bool {
		A.printPanicError("'" + lex + "'" + "does not name a type")
	}
	if lexType == lexinator.Long || lexType == lexinator.Short {
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Int {
			A.printPanicError("'" + lex + "'" + "does not name a type")
		}
	}
	var textPos, line, linePos int
	isFirst := true
	for isFirst || lexType == lexinator.Comma {
		isFirst = false
		A.variable()
		textPos, line, linePos = A.scanner.StorePosValues()
		lexType, lex = A.scanner.Scan()
	}
	A.scanner.RestorePosValues(textPos, line, linePos)
}

// <оператор> -> <составной оператор> | <for> | <процедура> ; | <присваивание>; | ;
func (A *Analyzer) operator() {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.OpeningBrace {
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.compositeOperator()
	} else if lexType == lexinator.For {
		A.scanner.RestorePosValues(textPos, line, linePos)
		A._for()
	} else if lexType == lexinator.Id {
		lexType, lex = A.scanner.Scan()
		if lexType == lexinator.OpeningBracket {
			A.scanner.RestorePosValues(textPos, line, linePos)
			A.procedure()
		} else if lexType == lexinator.Assignment {
			A.scanner.RestorePosValues(textPos, line, linePos)
			A.assigment()
		} else {
			A.printPanicError("'" + lex + "' is not an procedure or assigment")
		}
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.Semicolon {
			A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
		}
	} else if lexType != lexinator.Semicolon {
		A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
	}
}

// <описание процедуры> -> void идентификатор ( + <описание параметров> | e + ) <составной оператор>
func (A *Analyzer) procedureDescription() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.Void { // void
		A.printPanicError("invalid lexeme '" + lex + "', expected 'void'")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Id { // идентификатор
		A.printPanicError("'" + lex + "' is not an identifier")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.OpeningBracket { // (
		A.printPanicError("invalid lexeme '" + lex + "', expected '('")
	}
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex = A.scanner.Scan()
	if lexType == lexinator.Long || lexType == lexinator.Int ||
		lexType == lexinator.Short || lexType == lexinator.Bool { // <описание параметров>
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.parameterDescription()
		lexType, lex = A.scanner.Scan()
		if lexType != lexinator.ClosingBracket {
			A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
		}
	} else if lexType != lexinator.ClosingBracket { // )
		A.printPanicError("invalid lexeme '" + lex + "', expected ')'")
	}
	A.compositeOperator()
}

// done
// <составной оператор> -> { <операторы и описания> }
func (A *Analyzer) compositeOperator() {
	lexType, lex := A.scanner.Scan()
	if lexType != lexinator.OpeningBrace {
		A.printPanicError("invalid lexeme '" + lex + "', expected '{'")
	}
	A.operatorsAndDescriptions()
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.ClosingBrace {
		A.printPanicError("invalid lexeme '" + lex + "', expected '}'")
	}
}

// <операторы и описания> -> e | <операторы> U e | <операторы и описания>  | <описания> + e | <операторы и описания>
func (A *Analyzer) operatorsAndDescriptions() {
	// support functions
	isOperatorNext := func() bool {
		textPos, line, linePos := A.scanner.StorePosValues()
		lexType, _ := A.scanner.Scan()
		A.scanner.RestorePosValues(textPos, line, linePos)
		return lexType == lexinator.OpeningBrace || lexType == lexinator.For ||
			lexType == lexinator.Id || lexType == lexinator.Semicolon
	}
	isDescriptionNext := func() bool {
		textPos, line, linePos := A.scanner.StorePosValues()
		lexType, _ := A.scanner.Scan()
		A.scanner.RestorePosValues(textPos, line, linePos)
		return lexType == lexinator.Long || lexType == lexinator.Short ||
			lexType == lexinator.Int || lexType == lexinator.Bool || lexType == lexinator.Const
	}

	for isOperatorNext() || isDescriptionNext() {
		if isOperatorNext() {
			A.operator()
		} else {
			A.description()
		}
	}
}

// <описание> -> <переменные>; | <константы>;
func (A *Analyzer) description() {
	textPos, line, linePos := A.scanner.StorePosValues()
	lexType, lex := A.scanner.Scan()
	if lexType == lexinator.Long ||
		lexType == lexinator.Short ||
		lexType == lexinator.Int ||
		lexType == lexinator.Bool { // <переменные>
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.variables()
	} else if lexType == lexinator.Const { // <константы>
		A.scanner.RestorePosValues(textPos, line, linePos)
		A.constants()
	} else {
		A.printPanicError("'" + lex + "'" + "does not name a type")
	}
	lexType, lex = A.scanner.Scan()
	if lexType != lexinator.Semicolon {
		A.printPanicError("invalid lexeme '" + lex + "', expected ';'")
	}
}
