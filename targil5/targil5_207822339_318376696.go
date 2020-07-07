package main

//avital parasha 318376696
//tali kushelevsky 207822339

import (
Compiler "targil5_207822339_318376696/Compiling"
"targil5_207822339_318376696/Tokenizing"
"os"
"path/filepath"
"strings"
)

func main() {

	root := "./"
	err := filepath.Walk(root, func(path string, info os.FileInfo, errA error) error {
		if filepath.Ext(path) == ".jack" {
			Tokenizing.CreateXML(path)
		}
		if strings.Contains(path, "T.xml") {
			Compiler.Compile(path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}