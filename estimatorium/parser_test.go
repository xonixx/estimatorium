package estimatorium

import (
	"fmt"
	"testing"
)

const projData = `
project Project Name
author email@example.com

currency usd
time_unit day
acceptance_percent 10

risks low=1.1 medium=1.5 high=2

#desired_duration 1mth

team
b  cnt=1 rate=80 title=Blockchain
be cnt=2 rate=40
fe cnt=1 rate=30
qa cnt=1 rate=20 formula=(be+fe)*0.3
pm cnt=1 rate=50 formula=fe*0.33

tasks
Initial	|Research 		| be=3 fe=3  risks=low
Initial	|Bootstrap		| be=1 fe=10 risks=medium
API		| API task 1	| be=20      risks=high
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
func TestDesiredDuration(t *testing.T) {
	project := mustNoError(t, `
time_unit day
desired_duration 1mth
tasks
a|b|be=35 fe=2 risks=low
`)
	teamAsMap := project.TeamAsMap()
	if teamAsMap["be"].Count != 2 {
		t.Fatalf("be must be 2")
	}
	if teamAsMap["fe"].Count != 1 {
		t.Fatalf("fe must be 1")
	}
}

func TestDesiredDurationWithDerived(t *testing.T) {
	project := mustNoError(t, `
time_unit day
desired_duration 1mth
formula qa=be
tasks
a|b|be=35 fe=2 risks=low
`)
	teamAsMap := project.TeamAsMap()
	if teamAsMap["be"].Count != 2 {
		t.Fatalf("be must be 2")
	}
	if teamAsMap["fe"].Count != 1 {
		t.Fatalf("fe must be 1")
	}
	if teamAsMap["qa"].Count != 1 {
		t.Fatalf("qa must be 1")
	}
}

//func TestParsing3(t *testing.T) {
//	GenerateExcel(ProjectFromString(projData), "../Book3.xlsx")
//}
