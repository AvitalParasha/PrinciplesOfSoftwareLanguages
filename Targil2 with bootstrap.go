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

var (
	IfGt          int
	IfLt          int
	IfEq		  int
	callLbl       int = -1
)

func check(e error){
	if e != nil {
		panic(e)
	}
}

func HACKpush(pushParams []string) string{

	//pushParams[0] = push always...

	const (
		pushAndIncrement = "@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		pushPointer      = "D=M\n@SP\nM=M+1\nA=M-1\nM=D\n"
	)
	loadBaseAddr := "@" + pushParams[2] + "\nD=A\n"
	switch pushParams[1] {
	case "constant":
		return loadBaseAddr + pushAndIncrement
	case "this":
		return loadBaseAddr + "@THIS\nA=M+D\nD=M\n" + pushAndIncrement
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
	case "pointer":
		{
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
	case "pointer":
		{
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

func HACKnot() string{
	return "@SP\nA=M\nA=A-1\nM=!M\n"
}

func HACKneg() string{
	return "@SP\nA=M-1\nD=M\nM=-D\n"
}

func HACKBinLogic(operatorName string) string{
	var op string
	switch operatorName {
	case "and":	op = "&"
	case "or":	op = "|"
	}
	return "@SP\nA=M\nA=A-1\nD=M\nA=A-1\nM=M" + op + "D\n@SP\nM=M-1\n"
}

func HACKBinArithmetic(operatorName string) string{
	var op string
	switch operatorName{
	case "add": op = "+"
	case "sub":	op = "-"
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

func HACKLabel(labelName string, fileName string) string{
	return "//label\n("+strings.ToUpper(fileName + "." + labelName)+")\n"
}

func HACKGoto(labelName string, fileName string) string{
	return "//goto\n@" + strings.ToUpper(fileName + "." + labelName) + "\n0;JMP\n"
}

func HACKIfGoto(labelName string, fileName string) string{
	return "//if-goto\n@SP\nM=M-1\nA=M\nD=M\n@"+strings.ToUpper(fileName + "." + labelName)+"\nD;JNE\n"
}

func HACKfunction(line []string, fileName string) string {
	funcNameLabeld := strings.ToUpper(line[1])
	return "//label " + funcNameLabeld + "\n(" + funcNameLabeld + ")\n" +
		"//initialize local parameters\n@" + line[2] + "\nD=A\n@" +
		funcNameLabeld + ".End\nD;JEQ\n" +
		"//jump if not k!=0:\n(" + funcNameLabeld + ".Loop)\n@SP\nA=M\nM=0\n@SP\nM=M+1\n@" +
		funcNameLabeld + ".Loop\n" +
		"//jump while k!=0:\nD=D-1;JNE\n" +
		"//finish if k==0:\n(" + funcNameLabeld + ".End)\n"

}

func HACKcall(typePush []string, fileName string) string {
	n, err := strconv.Atoi(typePush[2])
	check(err)
	nPlus5 := strconv.Itoa(n+5)
	uppedName := strings.ToUpper(fileName + "." + typePush[1])
	labelName := strings.ToUpper(typePush[1] + ".ReturnAddress" + strconv.Itoa(callLbl))
	pushREtAdress := "@" + labelName + "\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
	positionArg :=   "@SP\nD=M\n@" + nPlus5 + "\nD=D-A\n@ARG\nM=D\n"
	positionLcl :=   "@SP\nD=M\n@LCL\nM=D\n"
	return "//function call\n" +
		"//PUSH RETURN ADDRESS\n" + pushREtAdress +
		"//PUSH LOCAL\n" + pushForFunc("LCL") +
		"//PUSH ARG\n" + pushForFunc("ARG")+
		"//PUSH THIS\n" + pushForFunc("THIS")+
		"//PUSH THAT\n" + pushForFunc("THAT") +
		"//ARG = SP-5-N\n" + positionArg +
		"//LCL = SP\n" + positionLcl +
		"//GOTO " + uppedName + "\n@" + strings.ToUpper(typePush[1]) + "\n0;JMP\n" +
		"//LABEL RETURN ADDRESS\n(" + labelName + ")\n"
}

func pushForFunc(typePush string) string{
	return "@" + typePush + "\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"

}

func HACKreturn() string {
	lclToTemp    := "@LCL\nD=M\n"
	retTotemp    := "\n@5\nA=D-A\nD=M\n@13\nM=D\n"
	returnValue  := "\n@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n"
	repositionSP := "\n@ARG\nD=M\n@SP\nM=D+1\n"
	restoreThat  := "@LCL\nM=M-1\nA=M\nD=M\n@THAT\nM=D\n"
	restoreThis  := "@LCL\nM=M-1\nA=M\nD=M\n@THIS\nM=D\n"
	restoreArg   := "@LCL\nM=M-1\nA=M\nD=M\n@ARG\nM=D\n"
	restoreLcl   := "@LCL\nM=M-1\nA=M\nD=M\n@LCL\nM=D\n"
	gotoRet      := "@13\nA=M\n0;JMP\n"
	return  "//FRAME=LCL\n" + lclToTemp +
		"//RET = *(FRAME-5)" +  retTotemp +
		"//*ARG = pop()" + returnValue +
		"//SP = ARG+1" + repositionSP  +
		"//RESTORE THAT\n" + restoreThat +
		"//RESTORE THIS\n" + restoreThis +
		"//RESTORE ARG\n" + restoreArg +
		"//RESTORE LCL\n" + restoreLcl +
		"//GOTO RET\n" + gotoRet

}

/*******************************************************************************/

func convertLine(line string, fileName string) string {
	croppedFileName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	switch line {
	case "add", "sub":
		return HACKBinArithmetic(line)
	case "and", "or":
		return HACKBinLogic(line)
	case "neg":
		return HACKneg()
	case "eq", "gt", "lt":
		return HACKcomp(line)
	case "not":
		return HACKnot()
	case "return":
		return HACKreturn()
	}
	if strings.Contains(line, "function"){
		return HACKfunction(strings.Split(line, " "), croppedFileName)
	}
	if strings.Contains(line, "call"){
		callLbl++
		return HACKcall(strings.Split(line, " "), croppedFileName) //had a label
	}
	if strings.Contains(line, "label"){
		return HACKLabel(strings.Split(line, " ")[1], croppedFileName)
	}
	if strings.Contains(line, "if-goto") {
		return HACKIfGoto(strings.Split(line, " ")[1], croppedFileName)
	}
	if strings.Contains(line, "goto") {
		return HACKGoto(strings.Split(line, " ")[1], croppedFileName)
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

func convertVmToAsm(fName string) string{
	fmt.Println(fName)

	IfEq, IfGt, IfLt = 0, 0, 0
	vmFile, err := os.OpenFile(fName, os.O_RDWR, 0755)
	check(err)
	fmt.Println(fName)
	ASMFileName := fName[0:len(fName)-2] + "asm"
	ASMFile, err := os.OpenFile(ASMFileName, os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	fmt.Println(ASMFile.Name())
	scanner := bufio.NewScanner(vmFile)
	writer := bufio.NewWriter(ASMFile)
	//header for RAM calculations, sp starts at 256 ram
	//fmt.Fprintln(writer, "@256\nD=A\n@SP\nM=D\n")
    var rv string
	for scanner.Scan() {
		line := scanner.Text()
		res := convertLine(line, vmFile.Name())
		if res != "" {
			rv +=  "//" + line + "\n" +  res
		}
	}
	writer.Flush()
	return rv
}

func main() {
	fmt.Println("Starts Compiling")
	dir , err := os.Open(".\\")
	if err != nil{
		panic(err)
	}
	dir2, err := os.Getwd()
	fmt.Println("that is dir 2: \n"+dir2)
	var ss [] string
	ss = strings.Split(dir2, "\\")
	fmt.Println("split:\n",ss)
	currentDirName:= ss[len(ss)-1]
	fmt.Println("current:\n",currentDirName)
	f1, _ := os.Create(currentDirName + ".asm")
	fmt.Println("f1:\n",f1.Name())
	filelist, _ := dir.Readdirnames(0)
	fmt.Println("file list:\n",filelist)
	f1.WriteString("@256\nD=A\n@SP\nM=D\n" + convertLine("call Sys.init 0", "")/* + "@END\n0;JMP\n"*/)
	fmt.Println("hi:\n")
	for _, filename := range filelist{
		if filename[len(filename) - 3:] == ".vm"{
			fmt.Println("in if:\n")
			f1.WriteString(strings.ToUpper(convertVmToAsm(filename)))
			fmt.Println("after convrt:\n")
		}
	}
	//f1.WriteString("(END)\n")
	f1.Close()
	fmt.Println("Compiled successfulTaliAvital")
}

