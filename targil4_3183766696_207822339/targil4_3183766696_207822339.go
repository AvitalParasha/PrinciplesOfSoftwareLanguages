// Avital Parasha 3183766696
// Tali Kushelevsky 207822339
package main

import (
	"ex4/Parsing"
	"ex4/Tokenizing"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	root := "./"
	err := filepath.Walk(root, func(path string, info os.FileInfo, errA error) error {
		if filepath.Ext(path)==".jack"  {
			Tokenizing.CreateXML(path)
		}
		if strings.Contains(path,"T.xml") {
			Parsing.Parse(path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}