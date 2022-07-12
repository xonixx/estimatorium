package main

import (
	"estimatorium/estimatorium"
	"fmt"
	"os"
)

const (
	version = "v0.0.1"
	usage   = "usage: ./estimator proj.txt report.xlsx"
)

func main() {
	args := os.Args

	realArgs := args[1:]

	if len(realArgs) == 0 || len(realArgs) == 1 && (realArgs[0] == "--help" || realArgs[0] == "-h") {
		fmt.Printf("Estimatorium %s\n%s\n", version, usage)
		os.Exit(0)
	} else if len(realArgs) == 0 || len(realArgs) == 1 && (realArgs[0] == "--version" || realArgs[0] == "-v") {
		fmt.Println(version)
		os.Exit(0)
	} else if len(realArgs) == 2 {
		bytes, err := os.ReadFile(realArgs[0])
		if err != nil {
			panic(err)
		}
		projectStr := string(bytes)
		project, err := estimatorium.ProjectFromString(projectStr)
		if err != nil {
			panic(err)
		}
		project.Calculate()
		fmt.Println(project)
		estimatorium.GenerateExcel(project, realArgs[1])
	}
}
