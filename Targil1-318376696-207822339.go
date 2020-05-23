// tali kushelevsky  207822339
//Avital Parasha 318376696

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var IfEq int
var IfGt int
var IfLt int

func check(e error){
	if e != nil {
		panic(e)
	}
}


func openFile(fName string) *os.File{
	file, err := os.Open(fName)
	check(err)
	defer file.Close()
	return file
}

func HACKnot() string{
	return "@SP\nA=M\nA=A-1\nM=!M\n"
}

func HACKneg() string{
	return "@SP\nA=M-1\nD=M\nM=-D\n"
}

func HACKBinLogic(operatorName string) string{
	var op string
	switch operatorName{

	case "and":
		op = "&"
	case "or":
		op = "|"
	}
	return "@SP\nA=M\nA=A-1\nD=M\nA=A-1\nM=M" + op + "D\n@SP\nM=M-1\n"

}
func HACKBinArithmetic(operatorName string) string{
	var op string
	switch operatorName{
	case "add":
		op = "+"
	case "sub":
		op = "-"

	}
	return "@SP\nM=M-1\n\n@SP\nA=M\nD=M\n\n@SP\nM=M-1\n\n@SP\nA=M\nA=M\n\nD=A" + op + "D\n\n@SP\nA=M\nM=D\n\n@SP\nM=M+1\n"
}

func HACKcomp(comp string) string {
	const comp2Vars = "@SP\nM=M-1\nA=M\nD=M\nA=A-1\nA=M\nD=A-D\n@"
	var compLblT, compLblF, jumpT string
	jumpF := "\n0;JEQ\n"
	switch comp {
	case "gt":
		jumpT = "\nD;JGT\n"
		compLblT = "TGT" + strconv.Itoa(IfGt)
		compLblF = "FGT" + strconv.Itoa(IfGt)
		IfGt++
	case "eq":
		jumpT = "\nD;JEQ\n"
		compLblT = "TEQ" + strconv.Itoa(IfEq)
		compLblF = "FEQ" + strconv.Itoa(IfEq)
		IfEq++
	case "lt":
		jumpT = "\nD;JLT\n"
		compLblT = "TLT" + strconv.Itoa(IfLt)
		compLblF = "FLT" + strconv.Itoa(IfLt)
		IfLt++
	}
	return comp2Vars + compLblT + jumpT + "@0\nA=M-1\nM=0\n@" + compLblF + jumpF + "(" + compLblT + ")\n@0\nA=M-1\nM=-1\n(" + compLblF + ")\n"
}


func HACKpop(popParams []string) string{

	const(
		decrAndPop = "@SP\nM=M-1\nA=M\nD=M\n"
		popToDest = "@13\nM=D\n@SP\nA=M\nD=M\n@13\nA=M\nM=D\n"
	)


	loadBaseAddr := "@SP\nM=M-1\n@" + popParams[2] + "\nD=A\n"

	switch popParams[1] {
	case "this":
		return loadBaseAddr + "@THIS\nD=M+D\n" + popToDest
	case "that":
		return loadBaseAddr + "@THAT\nD=M+D\n" + popToDest
	case "argument":
		return loadBaseAddr + "@ARG\nD=M+D\n" + popToDest
	case "local":
		return loadBaseAddr + "@LCL\nD=M+D\n" + popToDest
	case "temp":
		return loadBaseAddr + "@5\nD=A+D\n" + popToDest
	case "static":
		return loadBaseAddr + "@STATIC\nD=A+D\n" + popToDest
	case "pointer":{
		switch popParams[2] {
		case "0":
			return decrAndPop + "@THIS\nM=D\n"
		case "1":
			return decrAndPop + "@THAT\nM=D\n"
		}
	}
	}
	return ""
}




func HACKpush(pushParams []string) string{

	//pushParams[0] = push always...
	                
	const (
		pushAndIncrement = "@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		pushPointer = "D=M\n@SP\nM=M+1\nA=M-1\nM=D\n"
	)                                             
	loadBaseAddr := "@" + pushParams[2] + "\nD=A\n"
	switch pushParams[1] {

	case "constant":
		return loadBaseAddr + pushAndIncrement
	case "this":
		return loadBaseAddr + "@THIS\nA=M+D\nD=M\n" +  pushAndIncrement
	case "that":
		return loadBaseAddr + "@THAT\nA=M+D\nD=M\n" + pushAndIncrement
	case "argument":
		return loadBaseAddr + "@ARG\nA=M+D\nD=M\n" + pushAndIncrement
	case "local":
		return loadBaseAddr + "@LCL\nA=M+D\nD=M\n" + pushAndIncrement
	case "static":
		return loadBaseAddr + "@STATIC\nA=A+D\nD=M\n" + pushAndIncrement
	case "temp":
		return loadBaseAddr + "@5\nA=A+D\nD=M\n" + pushAndIncrement
	case "pointer":{
		switch pushParams[2] {
		case "0":
			return "@THIS\n" + pushPointer
		case "1":
			return "@THAT\n" + pushPointer
		}
	}
	}
	return ""
}




func convertLine(line string) string{
	switch line {
	case "add","sub":
		return HACKBinArithmetic(line)
	case "and","or":
		     return HACKBinLogic(line)
		case "neg":
		return HACKneg()
	case "eq", "gt", "lt":
		return HACKcomp(line)
	case "not":
		return HACKnot()
	}

	if strings.Contains(line, "push"){
		parts := strings.Split(line, " ")
		return HACKpush(parts)
	}
	if strings.Contains(line, "pop"){
		parts := strings.Split(line, " ")
		return HACKpop(parts)
	}
	return ""
}

func convertVmToAsm(fName string) {
	IfEq = 0
	IfGt = 0
	IfLt = 0
	vmFile, err := os.OpenFile(fName, os.O_RDWR, 0755)
	fmt.Println(fName)
	check(err)
	ASMFileName := fName[0:len(fName)-2] + "asm"
	ASMFile, err := os.OpenFile(ASMFileName, os.O_RDWR|os.O_CREATE, 0755)
	fmt.Println(ASMFile.Name())
	check(err)
	scanner := bufio.NewScanner(vmFile)
	writer := bufio.NewWriter(ASMFile)
	//header for RAM calculations sp starts at 256 ram
	fmt.Fprintln(writer, "@256\nD=A\n@SP\nM=D\n")

	for scanner.Scan() {
		line := scanner.Text()
		res := convertLine(line)
		if res != "" {
			fmt.Fprintln(writer, "//" + line + "\n" +  res)
		}
	}
	writer.Flush()
}






func main() {
	err := filepath.Walk("./", func(path string, info os.FileInfo, errA error) error {
		if filepath.Ext(info.Name()) ==".vm" {
			convertVmToAsm(path)
		}
		return nil
	})
	check(err)
}                

