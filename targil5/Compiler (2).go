//avital parasha 318376696
//tali kushelevsky 207822339

package Compiler

import (
	"bufio"
	etree "targil5_207822339_318376696/etree-master"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type hashT struct {
	name    string
	varType string
	kind    string
	index   int
}

var classScopeIdx int
var subroutineScopeIdx int

type symbolTableStruct struct {
	classScope, subroutineScope []hashT
}

//var sTable symbolTableStruct
var staticIndex int
var fieldIndex int
var varIndex int
var argIndex int
var kindOfSub string
var currentSub string

func (smblTable *symbolTableStruct) symbolTableConstructor() {
	staticIndex = 0
	fieldIndex = 0
	varIndex = 0
	argIndex = 0
	subroutineScopeIdx = 0
	classScopeIdx = 0
	(*smblTable).classScope = make([]hashT, 100)
	(*smblTable).subroutineScope = make([]hashT, 100)

}
func (smblTable *symbolTableStruct) getVarNum() int {
	return varIndex
}

func (smblTable *symbolTableStruct) startSubroutine() {
	subroutineScopeIdx = 0
	varIndex = 0
	argIndex = 0
}
func (smblTable *symbolTableStruct) define(name, varType, kind string) {
	switch kind {
	case "static":
		{
			(*smblTable).classScope[classScopeIdx].name = name
			(*smblTable).classScope[classScopeIdx].varType = varType
			(*smblTable).classScope[classScopeIdx].kind = "static"
			(*smblTable).classScope[classScopeIdx].index = staticIndex
			classScopeIdx++
			staticIndex++
			break
		}
	case "field":
		{
			(*smblTable).classScope[classScopeIdx].name = name
			(*smblTable).classScope[classScopeIdx].varType = varType
			(*smblTable).classScope[classScopeIdx].kind = "this"
			(*smblTable).classScope[classScopeIdx].index = fieldIndex
			classScopeIdx++
			fieldIndex++
			break
		}
	case "var":
		{
			(*smblTable).subroutineScope[subroutineScopeIdx].name = name
			(*smblTable).subroutineScope[subroutineScopeIdx].varType = varType
			(*smblTable).subroutineScope[subroutineScopeIdx].kind = "local"
			(*smblTable).subroutineScope[subroutineScopeIdx].index = varIndex
			subroutineScopeIdx++
			varIndex++
			break
		}
	case "arg":
		{
			(*smblTable).subroutineScope[subroutineScopeIdx].name = name
			(*smblTable).subroutineScope[subroutineScopeIdx].varType = varType
			(*smblTable).subroutineScope[subroutineScopeIdx].kind = "argument"
			(*smblTable).subroutineScope[subroutineScopeIdx].index = argIndex
			subroutineScopeIdx++
			argIndex++
			break
		}
	}
}

func varCount(kind string) int {
	switch kind {
	case "static":
		return staticIndex
	case "field":
		return fieldIndex
	case "var":
		return varIndex
	case "arg":
		return argIndex
	default:
		return -1
	}
}

func (smblTable *symbolTableStruct) get(name, tableName string) int {
	idx := 0
	switch tableName {
	case "subroutineScope":
		{
			for _, elem := range (*smblTable).subroutineScope {
				if name == elem.name {
					return idx
				}
				idx++
				if idx > subroutineScopeIdx {
					break
				}
			}
		}
	case "classScope":
		{
			for _, elem := range (*smblTable).classScope {
				if name == elem.name {
					return idx
				}
				idx++
				if idx > classScopeIdx {
					break
				}
			}
		}
	}
	return -1
}

//need to check the implementation
func (smblTable *symbolTableStruct) kindOf(name string) string {
	//checks in current subroutine first
	idx := (*smblTable).get(name, "subroutineScope")
	if idx != -1 {
		return (*smblTable).subroutineScope[idx].kind
	} else {
		idx = (*smblTable).get(name, "classScope")
		if idx != -1 {
			return (*smblTable).classScope[idx].kind
		} else {
			return "error"
		}
	}
}

func (smblTable *symbolTableStruct) typeOf(name string) string {
	//checks in current subroutine first
	idx := (*smblTable).get(name, "subroutineScope")
	if idx != -1 {
		return (*smblTable).subroutineScope[idx].varType
	} else {
		idx = (*smblTable).get(name, "classScope")
		if idx != -1 {
			return (*smblTable).classScope[idx].varType
		} else {
			return "error"
		}
	}
}

func (smblTable *symbolTableStruct) indexOf(name string) int {
	//checks in current subroutine first
	idx := (*smblTable).get(name, "subroutineScope")
	if idx != -1 {
		return (*smblTable).subroutineScope[idx].index
	} else {
		idx = (*smblTable).get(name, "classScope")
		if idx != -1 {
			return (*smblTable).classScope[idx].index
		} else {
			return -1
		}
	}
}

//var destElament *etree.Element
var srcElements []*etree.Element
var symbolT symbolTableStruct

var w *bufio.Writer
var className string
var elementsIdx int
var fieldNum int
var ifCounter int
var whileCounter int
var listNum int

func cmpStr(val string, cnst string) bool {

	return strings.Contains(val, cnst)
}
func isExpression(val string) bool {
	return cmpStr(val, "+") ||
		cmpStr(val, "-") ||
		cmpStr(val, "*") ||
		cmpStr(val, "/") ||
		cmpStr(val, "&") ||
		cmpStr(val, "|") ||
		cmpStr(val, "<") ||
		cmpStr(val, ">") ||
		cmpStr(val, "=")
}

func compileExpression() {
	var operator string
	compileTerm()
	for isExpression(srcElements[elementsIdx].Text()) {
		operator = srcElements[elementsIdx].Text()
		elementsIdx++
		compileTerm()

		if strings.Contains(operator, "+") {
			fmt.Fprintln(w, "add")
		} else if strings.Contains(operator, "-") {
			fmt.Fprintln(w, "sub")
		} else if strings.Contains(operator, "*") {
			fmt.Fprintln(w, "call Math.multiply 2")
		} else if strings.Contains(operator, "/") {
			fmt.Fprintln(w, "call Math.divide 2")
		} else if strings.Contains(operator, "&") {
			fmt.Fprintln(w, "and")
		} else if strings.Contains(operator, "|") {
			fmt.Fprintln(w, "or")
		} else if strings.Contains(operator, "<") {
			fmt.Fprintln(w, "lt")
		} else if strings.Contains(operator, ">") {
			fmt.Fprintln(w, "gt")
		} else if strings.Contains(operator, "=") {
			fmt.Fprintln(w, "eq")
		}
	}
}
func compileTerm() {
	var unaryOp string
	if cmpStr(srcElements[elementsIdx].Text(), "(") {
		// (
		elementsIdx++
		compileExpression()
		// )
		elementsIdx++
	} else if cmpStr(srcElements[elementsIdx].Text(), "-") || cmpStr(srcElements[elementsIdx].Text(), "~") {
		unaryOp = srcElements[elementsIdx].Text()
		// unaryOp
		elementsIdx++
		compileTerm()

		if cmpStr(unaryOp, "-") {
			fmt.Fprintln(w, "neg")
		} else if cmpStr(unaryOp, "~") {
			fmt.Fprintln(w, "not")
		}
	} else if cmpStr(srcElements[elementsIdx].Tag, "integerConstant") { //???????
		fmt.Fprintln(w, "push constant "+srcElements[elementsIdx].Text())
		elementsIdx++
	} else if cmpStr(srcElements[elementsIdx].Tag, "stringConstant") {
		str := srcElements[elementsIdx].Text()
		lenStr := len(str)
		if lenStr > 0 {
			fmt.Fprintln(w, "push constant "+strconv.Itoa(lenStr)) //length
			fmt.Fprintln(w, "call String.new 1")
		}
		byteArray := []byte(str)
		for i := range byteArray {
			fmt.Fprintln(w, "push constant "+strconv.Itoa(int(byteArray[i])))
			fmt.Fprintln(w, "call String.appendChar 2")
		}

		elementsIdx++
	} else if cmpStr(srcElements[elementsIdx].Text(), "true") {
		fmt.Fprintln(w, "push constant 0")
		fmt.Fprintln(w, "not")
		elementsIdx++
	} else if cmpStr(srcElements[elementsIdx].Text(), "false") || cmpStr(srcElements[elementsIdx].Text(), "null") {
		fmt.Fprintln(w, "push constant 0")
		elementsIdx++
	} else if cmpStr(srcElements[elementsIdx].Text(), "this") {
		fmt.Fprintln(w, "push pointer 0")
		elementsIdx++
	} else if cmpStr(srcElements[elementsIdx+1].Text(), "[") {
		var varName string
		varName = srcElements[elementsIdx].Text()
		// varName
		elementsIdx++
		// [
		elementsIdx++
		compileExpression()
		// ]
		elementsIdx++
		fmt.Fprintln(w, "push "+symbolT.kindOf(varName)+" "+strconv.Itoa(symbolT.indexOf(varName)))
		fmt.Fprintln(w, "add")
		fmt.Fprintln(w, "pop pointer 1")
		fmt.Fprintln(w, "push that 0")
	} else if cmpStr(srcElements[elementsIdx+1].Text(), "(") {
		var subName string
		fmt.Fprintln(w, "push argument 0")
		subName = srcElements[elementsIdx].Text()
		// subroutineName
		elementsIdx++
		// (
		elementsIdx++
		compileExpressionList()
		// )
		fmt.Fprintln(w, "call "+subName+strconv.Itoa((symbolT.getVarNum()+1)))
		elementsIdx++
	} else if cmpStr(srcElements[elementsIdx+1].Text(), ".") {
		var temp1 string
		var temp2 string
		temp1 = srcElements[elementsIdx].Text()
		// className|varName
		elementsIdx++
		// .
		elementsIdx++
		temp2 = srcElements[elementsIdx].Text()
		// subroutineName
		elementsIdx++
		// (
		elementsIdx++

		if !cmpStr(symbolT.kindOf(temp1), "error") {
			fmt.Fprintln(w, "push "+symbolT.kindOf(temp1)+" "+strconv.Itoa(symbolT.indexOf(temp1)))
			temp1 = symbolT.typeOf(temp1)
			listNum++
		}

		compileExpressionList()

		fmt.Fprintln(w, "call "+temp1+"."+temp2+" "+strconv.Itoa(listNum))
		listNum = 0
		// )
		elementsIdx++
	} else {
		fmt.Fprintln(w, "push "+symbolT.kindOf(srcElements[elementsIdx].Text())+" "+strconv.Itoa(symbolT.indexOf(srcElements[elementsIdx].Text())))
		// varName
		elementsIdx++
	}


}

func compileClassVarDec() {
	var name string
	var symbolType string
	var kind string

	kind = srcElements[elementsIdx].Text()
	elementsIdx++
	symbolType = srcElements[elementsIdx].Text()
	elementsIdx++
	name = srcElements[elementsIdx].Text()
	elementsIdx++

	symbolT.define(name, symbolType, kind)

	for cmpStr(srcElements[elementsIdx].Text(), ",") {
		if cmpStr(kind, "field") {
			fieldNum++
		}
		//destElament.AddChild(srcElements[elementsIdx]) //,
		elementsIdx++
		name = srcElements[elementsIdx].Text() //identifier
		symbolT.define(name, symbolType, kind)
		elementsIdx++
	}

	// ;
	elementsIdx++
}
func isSubroutine(val string) bool {
	var routines = "constructor function method"
	return strings.Contains(routines, val)
}
func compileSubroutineDec() {
	kindOfSub = srcElements[elementsIdx].Text()

	elementsIdx++
	// Return type
	elementsIdx++
	// Subroutine name
	currentSub = srcElements[elementsIdx].Text()

	elementsIdx++
	// (
	elementsIdx++

	compileParameterList()

	// )
	elementsIdx++

	compileSubroutineBody()

}
func compileParameterList() {

	kind := "arg"

	if !cmpStr(srcElements[elementsIdx].Text(), ")") {

		varType := srcElements[elementsIdx].Text() //type
		elementsIdx++
		name := srcElements[elementsIdx].Text() //varName
		elementsIdx++
		symbolT.define(name, varType, kind)
		for cmpStr(srcElements[elementsIdx].Text(), ",") {
			//,
			elementsIdx++
			varType = srcElements[elementsIdx].Text() //type
			elementsIdx++
			name = srcElements[elementsIdx].Text() //varName
			elementsIdx++
			symbolT.define(name, varType, kind)
		}
	}

}

func compileSubroutineBody() {
	// {
	elementsIdx++

	for cmpStr(srcElements[elementsIdx].Text(), "var") {
		compileVarDec()
	}
	fmt.Fprintln(w, "function "+className+"."+currentSub+" "+strconv.Itoa(symbolT.getVarNum()))
	if kindOfSub == "method" {
		fmt.Fprintln(w, "push argument 0")
		fmt.Fprintln(w, "pop pointer 0")
	}
	if kindOfSub == "constructor" {
		fmt.Fprintln(w, "push constant "+strconv.Itoa(fieldNum))
		fmt.Fprintln(w, "call Memory.alloc 1")
		fmt.Fprintln(w, "pop pointer 0")
	}

	compileStatements()

	// }
	elementsIdx++
}
func compileVarDec() {
	var name string
	var symboltype string
	kind := "var"
	// var
	elementsIdx++
	symboltype = srcElements[elementsIdx].Text()
	elementsIdx++
	name = srcElements[elementsIdx].Text()
	elementsIdx++
	symbolT.define(name, symboltype, kind)

	for cmpStr(srcElements[elementsIdx].Text(), ",") {
		// ,
		elementsIdx++
		name = srcElements[elementsIdx].Text()
		elementsIdx++
		symbolT.define(name, symboltype, kind)
	}

	// ;
	elementsIdx++
}

func compileStatements() {
	for !cmpStr(srcElements[elementsIdx].Text(), "}") {
		compileStatement()
	}
}
func isStatement(val string) bool {
	const statements = "let if while do return"
	return strings.Contains(statements, val)
}

func compileStatement() {
	for isStatement(srcElements[elementsIdx].Text()) {
		switch srcElements[elementsIdx].Text() {
		case "let":
			compileLetStatement()
			break
		case "if":
			compileIfStatement()
			break
		case "while":
			compileWhileStatement()
			break
		case "do":
			compileDoStatement()
			fmt.Fprintln(w, "pop temp 0")
			break
		default:
			compileReturnStatement()
			break
		}
		if cmpStr(srcElements[elementsIdx].Text(), ";") {
			elementsIdx++
		}

	}
}
func compileLetStatement() {
	arrayFlag := false
	var varName string
	elementsIdx++
	varName = srcElements[elementsIdx].Text()
	elementsIdx++
	if cmpStr(srcElements[elementsIdx].Text(), "[") {
		arrayFlag = true
		elementsIdx++
		compileExpression()
		elementsIdx++
		fmt.Fprintln(w, "push "+symbolT.kindOf(varName)+" "+strconv.Itoa(symbolT.indexOf(varName)))
		fmt.Fprintln(w, "add")
	}
	// =
	elementsIdx++
	compileExpression()
	if arrayFlag {
		fmt.Fprintln(w, "pop temp 0")
		fmt.Fprintln(w, "pop pointer 1")
		fmt.Fprintln(w, "push temp 0")
		fmt.Fprintln(w, "pop that 0")
	} else {
		fmt.Fprintln(w, "pop "+symbolT.kindOf(varName)+" "+strconv.Itoa(symbolT.indexOf(varName)))
	}
	// ;
	elementsIdx++
}

func compileIfStatement() {
	var ifCounterLocal int
	ifCounter++
	ifCounterLocal = ifCounter
	// if
	elementsIdx++
	// (
	elementsIdx++
	compileExpression()
	// )
	fmt.Fprintln(w, "if-goto IF_TRUE"+strconv.Itoa(ifCounterLocal))
	fmt.Fprintln(w, "goto IF_FALSE"+strconv.Itoa(ifCounterLocal))
	fmt.Fprintln(w, "label IF_TRUE"+strconv.Itoa(ifCounterLocal))
	elementsIdx++
	// {
	elementsIdx++
	compileStatements()
	// }
	elementsIdx++
	chkElse := "0"
	if cmpStr(srcElements[elementsIdx].Text(), "else") {
		fmt.Fprintln(w, "goto IF_END"+strconv.Itoa(ifCounterLocal))
		chkElse = "1"
	}

	fmt.Fprintln(w, "label IF_FALSE"+strconv.Itoa(ifCounterLocal))

	if cmpStr(srcElements[elementsIdx].Text(), "else") {
		// else
		elementsIdx++
		// {
		elementsIdx++
		compileStatements()
		// }
		elementsIdx++
	}
	if chkElse == "1" {
		fmt.Fprintln(w, "label IF_END"+strconv.Itoa(ifCounterLocal))
	}
}

func compileWhileStatement() {
	var whileCounterLocal int
	whileCounter++
	whileCounterLocal = whileCounter
	fmt.Fprintln(w, "label WHILE_EXP"+strconv.Itoa(whileCounterLocal))
	// while
	elementsIdx++
	// (
	elementsIdx++
	compileExpression()
	// )
	elementsIdx++
	fmt.Fprintln(w, "not")
	fmt.Fprintln(w, "if-goto WHILE_END"+strconv.Itoa(whileCounterLocal))
	// {
	elementsIdx++
	compileStatements()
	// }
	elementsIdx++
	fmt.Fprintln(w, "goto WHILE_EXP"+strconv.Itoa(whileCounterLocal))
	fmt.Fprintln(w, "label WHILE_END"+strconv.Itoa(whileCounterLocal))
}

func compileDoStatement() {
	var temp1 string
	var temp2 string
	// do
	elementsIdx++
	switch srcElements[elementsIdx+1].Text() {
	case "(":
		{
			temp1 = srcElements[elementsIdx].Text()
			// subroutineName
			elementsIdx++
			// (
			elementsIdx++

			fmt.Fprintln(w, "push pointer 0")

			compileExpressionList()

			fmt.Fprintln(w, "call "+className+"."+temp1+" "+strconv.Itoa(listNum+1))
			listNum = 0
			// )
			elementsIdx++
			break
		}
	case ".":
		{
			// className|varName
			temp1 = srcElements[elementsIdx].Text()
			elementsIdx++
			// .
			elementsIdx++
			// subroutineName
			temp2 = srcElements[elementsIdx].Text()
			elementsIdx++
			// (
			elementsIdx++

			if !cmpStr(symbolT.kindOf(temp1), "error") {
				fmt.Fprintln(w, "push "+symbolT.kindOf(temp1)+" "+strconv.Itoa(symbolT.indexOf(temp1)))
				temp1 = symbolT.typeOf(temp1)
				listNum++
			}

			compileExpressionList()

			fmt.Fprintln(w, "call "+temp1+"."+temp2+" "+strconv.Itoa(listNum))
			listNum = 0
			// )
			elementsIdx++
			break
		}
	}
	// ;
	elementsIdx++
}
func compileReturnStatement() {
	elementsIdx++
	if !cmpStr(srcElements[elementsIdx].Text(), ";") {
		compileExpression()
	} else {
		fmt.Fprintln(w, "push constant 0")
	}

	fmt.Fprintln(w, "return")
	// ;
	elementsIdx++
}

func compileSubroutineCall() {
	//destElament.AddChild(srcElements[elementsIdx]) //identifier
	elementsIdx++
	if cmpStr(srcElements[elementsIdx].Text(), "(") {
		//destElament.AddChild(srcElements[elementsIdx]) //(
		elementsIdx++
		compileExpressionList()
		//destElament.AddChild(srcElements[elementsIdx]) //)
		elementsIdx++
	} else {
		//destElament.AddChild(srcElements[elementsIdx]) //.
		elementsIdx++
		//destElament.AddChild(srcElements[elementsIdx]) //identifier
		elementsIdx++
		//destElament.AddChild(srcElements[elementsIdx]) //(
		elementsIdx++
		compileExpressionList()
		//destElament.AddChild(srcElements[elementsIdx]) //)
		elementsIdx++
	}
}

func compileExpressionList() {
	if !cmpStr(srcElements[elementsIdx].Text(), ")") {
		listNum++
		compileExpression()
		for cmpStr(srcElements[elementsIdx].Text(), ",") {
			// ,
			listNum++
			elementsIdx++
			compileExpression()
		}
	}
}

func compileClass() {
	// class
	elementsIdx++
	// className
	elementsIdx++
	// {
	elementsIdx++
	for cmpStr(srcElements[elementsIdx].Text(), "static") || cmpStr(srcElements[elementsIdx].Text(), "field") {
		if cmpStr(srcElements[elementsIdx].Text(), "field") {
			fieldNum++
		}
		compileClassVarDec()
	}

	for elementsIdx < len(srcElements) && (cmpStr(srcElements[elementsIdx].Text(), "constructor") || cmpStr(srcElements[elementsIdx].Text(), "function") || cmpStr(srcElements[elementsIdx].Text(), "method")) {
		symbolT.startSubroutine()
		ifCounter = -1
		whileCounter = -1
		if cmpStr(srcElements[elementsIdx].Text(), "method") {
			symbolT.define("this", className, "arg")
		}
		compileSubroutineDec()
	}
	// }
	elementsIdx++
}
func Compile(path string) {
	fieldNum = 0
	srcDoc := etree.NewDocument()
	if err := srcDoc.ReadFromFile(path); err != nil {
		panic(err)
	}
	symbolT = symbolTableStruct{classScope: nil, subroutineScope: nil}
	symbolT.symbolTableConstructor()
	newFileName := strings.Replace(path, "T.xml", ".vm", -1)
	//fmt.Println(newFileName)
	wfile, err := os.Create(newFileName)
	if err != nil {
		panic(err)
	}
	defer wfile.Close()

	w = bufio.NewWriter(wfile)

	elementsIdx = 0
	className = srcDoc.ChildElements()[0].ChildElements()[1].Text()
	srcElements = srcDoc.Element.ChildElements()[0].ChildElements()
	compileClass()
	elementsIdx = 2

	w.Flush()
}
