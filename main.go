package main

import (
	"estimatorium/core"
	"fmt"
	"os"
)

const (
	version = "0.0.1"
	usage   = "usage: ./estimatorium proj.txt report.xlsx"
)

func main() {
	args := os.Args

	realArgs := args[1:]

	if len(realArgs) == 0 || len(realArgs) == 1 && (realArgs[0] == "--help" || realArgs[0] == "-h") {
		fmt.Printf("Estimatorium v%s\n%s\n", version, usage)
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
		project, err := core.ProjectFromString(projectStr)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		project.Calculate()
		fmt.Println(project)
		core.GenerateExcel(project, realArgs[1])
	} else {
		fmt.Printf("I don't understand...\n%s\n", usage)
		os.Exit(1)
	}
}
