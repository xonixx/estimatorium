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
	fmt.Println(parseProj(projData))
}
func TestParsing2(t *testing.T) {
	fmt.Println(ProjectFromString(projData))
}
func TestParsing3(t *testing.T) {
	GenerateExcel(ProjectFromString(projData), "../Book3.xlsx")
}
