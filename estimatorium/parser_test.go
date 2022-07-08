package estimatorium

import (
	"fmt"
	"testing"
)

const projData = `
project Project Name
author email@example.com

# params currency=usd time_unit=day acceptancePercent=10
currency usd
time_unit day
acceptance_percent 10

risks low=1.1 medium=1.5 high=2
rates be=40 fe=30 qa=20
formula qa=(be+fe)*0.3 pm=fe*0.33

desired_duration 3mth
team be=2 fe=1 qa=1 pm=1

tasks

Initial	|Research 		| be=3 fe=3 risks=low
Initial	|Bootstrap		| be=1 fe=1 risks=low
API		| API task 1	| be=2 
API		| API task 2 	| be=2 
`

func TestParsing1(t *testing.T) {
	proj := parseProj(projData)
	if len(proj.tasksRecords) != 4 {
		t.Fatalf("wrong tasksRecord cnt")
	}
	if *proj.getSingleVal("currency") != "usd" {
		t.Fatalf("wrong currency")
	}
	//fmt.Println(proj)
}
func TestParsing2(t *testing.T) {
	project, err := ProjectFromString(projData)
	if err != nil {
		t.Fatal(err)
	}
	if len(project.Tasks) != 4 {
		t.Fatalf("wrong tasksRecord cnt")
	}
	if project.Currency != Usd {
		t.Fatalf("wrong currency")
	}
	//fmt.Println(project)
}

func TestWrongCurrency(t *testing.T) {
	project, err := ProjectFromString(`
currency wrong
`)
	if err == nil {
		t.Fatalf("should be error")
	}
	fmt.Println(project)
}

//func TestParsing3(t *testing.T) {
//	GenerateExcel(ProjectFromString(projData), "../Book3.xlsx")
//}
