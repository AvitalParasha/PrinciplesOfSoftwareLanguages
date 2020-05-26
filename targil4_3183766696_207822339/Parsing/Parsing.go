// Avital Parasha 3183766696
// Tali Kushelevsky 207822339
package Parsing

import (
	etree "etree-master"
	"os"
	"strings"
)

var destElament *etree.Element
var newDoc *etree.Document
var srcElements []*etree.Element
var symbolText string
var symbolText2 string
var elementsIdx int
var text string
var isFrstStat bool

func cmpStr(val string,cnst string) bool{
	return strings.Contains(val,cnst)
}
func copyFromTokFile(num int) {
	for i:=0 ; i < num; i++ {
		destElament.AddChild(srcElements[elementsIdx])
		elementsIdx++
	}
}
func parseClass(){
	destElament = destElament.CreateElement("class")
	newDoc.AddChild(destElament)
	copyFromTokFile(3)
	parseClassVarDec()
	parseSubroutineDec()
	copyFromTokFile(1)
	destElament = destElament.Parent()

}

func parseClassVarDec(){
	for cmpStr(srcElements[elementsIdx].Text(),"static") || cmpStr(srcElements[elementsIdx].Text(),"field"){
			destElament = destElament.CreateElement("classVarDec")
			copyFromTokFile(3) //static or field, //int or char or boolean or identifier, //identifier
			for cmpStr(srcElements[elementsIdx].Text(),","){
				copyFromTokFile(2) //,//identifier
			}
			copyFromTokFile(1) //;
			destElament = destElament.Parent()
		}
}

func isSubroutine(val string)bool{
	var routines = "constructor function method"
	return strings.Contains(routines, val)
}
func parseSubroutineDec(){
	for isSubroutine(srcElements[elementsIdx].Text()) {
		destElament = destElament.CreateElement("subroutineDec")
		copyFromTokFile(4)//function or const or method, //void or type,//identifier, //(
		parseParameterList()
		copyFromTokFile(1) //)
		parseSubroutineBody()
		destElament = destElament.Parent()
	}
	
}

func parseParameterList(){
	destElament = destElament.CreateElement("parameterList")
	if !cmpStr(srcElements[elementsIdx].Text(),")"){
		copyFromTokFile(2) //type, //varName
		for cmpStr(srcElements[elementsIdx].Text(),",") {
			copyFromTokFile(3) //, , //type , //varName
		}
	}
	destElament = destElament.Parent()
}

func parseSubroutineBody(){
	destElament = destElament.CreateElement("subroutineBody")
	copyFromTokFile(1) //{
	parseVarDec()
	parseStatements()
	copyFromTokFile(1) //}
	destElament = destElament.Parent()
}

func parseVarDec(){
	for cmpStr(srcElements[elementsIdx].Text(),"var"){
		destElament = destElament.CreateElement("varDec")
		copyFromTokFile(3) //var, //type, //varName
		for cmpStr(srcElements[elementsIdx].Text(),","){
			copyFromTokFile(2) //,//varName

		}	
		copyFromTokFile(1) //;
		destElament = destElament.Parent()
	}
}

func parseStatements(){
	destElament = destElament.CreateElement("statements")
	parseStatement()
	destElament = destElament.Parent()
}
func isStatement(val string) bool{
	const statements = "let if while do return";
	return strings.Contains(statements,val)
}
func parseStatement(){
	for isStatement(strings.TrimSpace(srcElements[elementsIdx].Text())) {
		switch srcElements[elementsIdx].Text() {
		case "let":
			parseLetStatement()
			break
		case "if":
			parseIfStatement()
			break
		case "while":
			parseWhileStatement()
			break
		case "do":
			parseDoStatement()
			break
		default:
			parseReturnStatement()
			break
		}
	}
}

func parseLetStatement(){
	destElament = destElament.CreateElement("letStatement")
	copyFromTokFile(2) //let, //varName
	if cmpStr(srcElements[elementsIdx].Text(),"["){
		copyFromTokFile(1) //[
		parseExpression()
		copyFromTokFile(1) //]
	}
	copyFromTokFile(1) //=
	parseExpression()
	copyFromTokFile(1) //;
	destElament = destElament.Parent()
}

func parseIfStatement(){
	destElament = destElament.CreateElement("ifStatement")
	copyFromTokFile(2) //if  //(
	parseExpression()
	copyFromTokFile(2) //) {
	parseStatements()
	copyFromTokFile(1) //}
	if cmpStr(srcElements[elementsIdx].Text(),"else"){
		copyFromTokFile(2) //else {
		parseStatements()
		copyFromTokFile(1) //}
	}
	destElament = destElament.Parent()
}

func parseWhileStatement(){
	destElament = destElament.CreateElement("whileStatement")
	copyFromTokFile(2)//while(
	parseExpression()
	copyFromTokFile(2) //){
	parseStatements()
	copyFromTokFile(1) //}
	destElament = destElament.Parent()
}

func parseDoStatement(){
	destElament = destElament.CreateElement("doStatement")
	copyFromTokFile(1)//do
	parseSubroutineCall()
	copyFromTokFile(1) //;
	destElament = destElament.Parent()
}


func parseReturnStatement(){
	destElament = destElament.CreateElement("returnStatement")
	copyFromTokFile(1)//return
	if !cmpStr(srcElements[elementsIdx].Text(),";"){
		parseExpression()
	}
	copyFromTokFile(1) //;
	destElament = destElament.Parent()
}

func isExpression(val string) bool{
	return cmpStr(val,"+") ||
		cmpStr(val,"-") ||
		cmpStr(val,"*") ||
		cmpStr(val,"/") ||
    	cmpStr(val,"|") ||
		cmpStr(val,"&") ||
		cmpStr(val,">") ||
		cmpStr(val,"<") ||
		cmpStr(val,"=")
}

func parseExpression(){
	destElament = destElament.CreateElement("expression")
	parseTerm()
	for isExpression(srcElements[elementsIdx].Text()){
		copyFromTokFile(1) //op
		parseTerm()
	}
	destElament = destElament.Parent()
}

func parseSubroutineCall(){
	copyFromTokFile(1) //identifier
	if cmpStr(srcElements[elementsIdx].Text(),"("){
		copyFromTokFile(1) //(
		parseExpressionList()
		copyFromTokFile(1) //)
	}else{
		copyFromTokFile(3)//. identifier (
		parseExpressionList()
		copyFromTokFile(1) //)
	}

}

func parseExpressionList(){
	destElament = destElament.CreateElement("expressionList")
	if !cmpStr(srcElements[elementsIdx].Text(),")"){
		parseExpression()
	}
	for cmpStr(srcElements[elementsIdx].Text(),","){
		copyFromTokFile(1) //,
		parseExpression()
	}
	destElament = destElament.Parent()
}

func parseTerm(){
	destElament = destElament.CreateElement("term")
	if cmpStr(srcElements[elementsIdx].Text(),"("){
		copyFromTokFile(1) //(
		parseExpression()
		copyFromTokFile(1) //)
	}else if cmpStr(srcElements[elementsIdx+1].Text(),"["){
		copyFromTokFile(2) //varName[
		parseExpression()
		copyFromTokFile(1) //]
	}else if cmpStr(srcElements[elementsIdx].Text(),"-") ||
		cmpStr(srcElements[elementsIdx].Text(),"~"){
		copyFromTokFile(1) //unaryOp
		parseTerm()
	}else if cmpStr(srcElements[elementsIdx+1].Text(),"(") ||
				cmpStr(srcElements[elementsIdx+1].Text(),"."){
		parseSubroutineCall()
	}else{
		copyFromTokFile(1) //indentifier
	}
	destElament = destElament.Parent()
}


func Parse(path string) {
	srcDoc := etree.NewDocument()
	if err := srcDoc.ReadFromFile(path); err != nil {
		panic(err)
	}
	newDoc = etree.NewDocument()
	isFrstStat = true

	newFileName := strings.Replace(path,"T.xml",".xml",-1)
	wfile, err := os.Create(newFileName)
	if err != nil {
		panic(err)
	}
	defer wfile.Close()

	srcElements =  srcDoc.Element.ChildElements()[0].ChildElements()
	elementsIdx = 0;

	parseClass()
	elementsIdx = 0;
	newDoc.Indent(2)
	newDoc.WriteTo(wfile)
}
