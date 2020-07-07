// Avital Parasha 3183766696
// Tali Kushelevsky 207822339
package Tokenizing

import (
	"bufio"
	etree "etree-master"
	"os"
	"regexp"
	"strings"
)

var keywords = []string{"class","constructor","function","method","field","static","var","int","char","boolean","void","true","false","null","this","let","do","if","else","while","return"}
var isSymbol = "{}()[].,;+-*/&<>=~"

func isIntegerConstant (value string) bool{
	matched, _ := regexp.MatchString("^\\d*\\.{0,1}\\d+$", value)
	return matched
}
func isStringConstant (value string) bool{
	matched, _ := regexp.MatchString("^\"[^\\\"]*\"$", value)
	return matched
}

func isIdentifier(value string) bool {
	matched, _ := regexp.MatchString("^[^\\d][\\d\\w\\_]*", value)
	return matched
}

func isKeyword (value string) bool{
	for _, v := range keywords {
		if v == value {
			return true
		}
	}
	return false
}
func GetNodeName(str string) string{
	if isKeyword(str){
		return "keyword"
	}else if  strings.Contains(isSymbol,str){
		return "symbol"
	}else if isIntegerConstant(str){
		return "integerConstant"
	}else if isStringConstant(str){
		return "stringConstant"
	}else if isIdentifier(str){
		return "identifier"
	}
	return ""
}


func GetContent(content string) string {

	return strings.Replace(content,"\"","",-1)
}

var regularExpr = regexp.MustCompile(`(?:\/\/.*|\/\*|\*\/|\<|\>|\.|#|&|\,|:|\||\*|\(|\)|=|\{|\}|\(|\)|\[|\]|\.|\;|\+|\-|\*|\/|\&|\|\|\=|\~|\"[^\"]*\"|\d+\.{0,1}\d*|\s|\n|\w+)?`)



func CreateXML(path string)   {
	doc := etree.NewDocument()
	tokens := doc.CreateElement("tokens")

	rfile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer rfile.Close()
	newFileName := strings.Replace(path,".jack","T.xml",-1)
	wfile, err := os.Create(newFileName)
	if err != nil {
		panic(err)
	}

	defer wfile.Close()


	scanner := bufio.NewScanner(rfile)
	for scanner.Scan() {
		var line = strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line,"/*") &&  !strings.HasPrefix(line,"*")  && !strings.HasPrefix(line,"//"){
			if strings.Contains(line,"//"){
				line = strings.Split(line, "//")[0]
			}
			var items = regularExpr.FindAllStringSubmatch(line, -1)
			for _, item := range items {
				var currentItem =  strings.TrimSpace(item[0])
				if currentItem != ""  {
					var elem = tokens.CreateElement(GetNodeName(currentItem))
					elem.SetText( GetContent(currentItem))
				}
			}
		}

	}

	doc.Indent(2)
	doc.WriteTo(wfile)

}