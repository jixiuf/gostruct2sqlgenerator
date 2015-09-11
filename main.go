package main

import (
	"fmt"
	"github.com/0studio/gostruct2sqlgenerator/generator"
	"os"
	"path/filepath"
	"strings"
)

// give abc.txt return abc
// give a/b/abc.txt return abc
func getFileName(fileName string) string {
	fileName = filepath.Base(fileName)
	idx := strings.Index(fileName, ".")
	if idx != -1 {
		return fileName[:idx]
	}
	return fileName
}

// go run /main.go example/example_1.go
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("please give a go struct defintion file as params like this : %s\n go_struct.go", os.Args[0])
		return
	}
	goStructFile := os.Args[1]
	srcDir := filepath.Dir(goStructFile)
	if !strings.HasSuffix(goStructFile, ".go") {
		fmt.Printf("the first param must be a go source file ,and some struct are defined there\n")
		return
	}
	structDescriptionList, property := generator.ParseStructFile(goStructFile)
	if len(structDescriptionList) == 0 {
		fmt.Println("no struct found in ", goStructFile)
		return
	}
	outputF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("%s_sub.go", getFileName(goStructFile))), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer outputF.Close()
	outputF.WriteString("package ")
	outputF.WriteString(property.PackageName)
	outputF.WriteString(" // do not edit this file ,this is generated by tools(https://github.com/0studio/gostruct2sqlgenerator)\n\n")
	outputF.WriteString("import (\n")
	outputF.WriteString("    \"fmt\"\n")
	outputF.WriteString("    \"time\"\n")
	outputF.WriteString(")\n\n")
	outputF.WriteString("var ___importTime time.Time\n\n")
	outputF.WriteString("func bool2int(b bool) int {\n")
	outputF.WriteString("    if b {\n")
	outputF.WriteString("        return 1\n")
	outputF.WriteString("    } else {\n")
	outputF.WriteString("        return 0\n")
	outputF.WriteString("    }\n")
	outputF.WriteString("}\n")
	for _, sd := range structDescriptionList {
		outputF.WriteString(sd.GenerateInsert())
		outputF.WriteString("\n")
		outputF.WriteString(sd.GenerateCreateTableFunc())
		outputF.WriteString("\n")
	}

	sqlF, err := os.OpenFile(filepath.Join(srcDir, fmt.Sprintf("%s_create_table.sql", getFileName(goStructFile))), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sqlF.Close()
	for _, sd := range structDescriptionList {
		sql, err := sd.GenerateCreateTableSql()
		if err != nil {
			fmt.Println(err)
			continue
		}

		sqlF.WriteString(sql)
		sqlF.WriteString("\n")
	}
}
