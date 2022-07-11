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
	proj, err := parseProj(projData)
	if err != nil {
		t.Fatal(err)
	}
	if len(proj.tasksRecords) != 4 {
		t.Fatalf("wrong tasksRecord cnt")
	}
	if *proj.getSingleVal(directiveCurrency) != "usd" {
		t.Fatalf("wrong currency")
	}
	//fmt.Println(proj)
}
func TestParsing2(t *testing.T) {
	project := mustNoError(t, projData)
	if len(project.Tasks) != 4 {
		t.Fatalf("wrong tasksRecord cnt")
	}
	if project.Currency != Usd {
		t.Fatalf("wrong currency")
	}
	//fmt.Println(project)
}

func mustBeError(t *testing.T, s string) {
	project, err := ProjectFromString(s)
	fmt.Println(project)
	fmt.Println(err)
	if err == nil {
		t.Fatalf("should be error")
	}
}

func mustNoError(t *testing.T, s string) Project {
	project, err := ProjectFromString(s)
	project.Calculate()
	fmt.Println(project)
	fmt.Println(err)
	if err != nil {
		t.Fatalf("should be error")
	}
	return project
}

func TestWrongCurrency(t *testing.T) {
	mustBeError(t, `currency wrong`)
}

func TestWrongTimeUnit(t *testing.T) {
	mustBeError(t, `time_unit wrong`)
}

func TestWrongAcceptancePercent1(t *testing.T) {
	mustBeError(t, `acceptance_percent wrong`)
}
func TestWrongAcceptancePercent2(t *testing.T) {
	mustBeError(t, `acceptance_percent 123`)
}
func TestWrongRiskVal1(t *testing.T) {
	mustBeError(t, `risks be=0.5`)
}
func TestWrongRiskVal2(t *testing.T) {
	mustBeError(t, `risks be=abc`)
}
func TestWrongDirective(t *testing.T) {
	mustBeError(t, `
wrong 123
wrong1 hello`)
}

func TestRepeatingDirective(t *testing.T) {
	mustBeError(t, `
currency usd
currency eur`)
}
func TestWrongRiskName(t *testing.T) {
	mustBeError(t, `
risks aaa=10
tasks
a|b|risks=wrong`)
}
func TestDefaultRisksApply(t *testing.T) {
	mustNoError(t, `
team be=1
tasks
a|b|be=1 risks=low
`)
}
func TestWrongResourceName(t *testing.T) {
	mustBeError(t, `
team be=1
tasks
a|b|zz=1`)
}
func TestNoTeam(t *testing.T) {
	mustNoError(t, `
time_unit day
desired_duration 1mth
tasks
a|b|be=1 fe=2 risks=low
`)
}

//func TestParsing3(t *testing.T) {
//	GenerateExcel(ProjectFromString(projData), "../Book3.xlsx")
//}
