package main

import (
	"estimatorium/estimatorium"
	"fmt"
)

func main() {
	project := estimatorium.Project{
		Name:              "proj name",
		TimeUnit:          estimatorium.Day,
		AcceptancePercent: 10,
		Currency:          estimatorium.Usd,
		Risks:             estimatorium.StandardRisks(),
		Team: []estimatorium.Resource{
			{Id: "fe", Title: "Front dev", Rate: 40, Count: 1},
			{Id: "be", Title: "Back dev", Rate: 50, Count: 2},
			{Id: "qa", Title: "QA", Rate: 35, Formula: "(fe + be)*0.3", Count: 1},
		},
		Tasks: []estimatorium.Task{
			{Category: "Initial", Title: "Task 1", Risk: "low", Work: map[string]float32{"be": 2, "fe": 5}},
			{Category: "Feature 2", Title: "Task 1", Risk: "low", Work: map[string]float32{"be": 2, "fe": 5}},
			{Category: "Feature 2", Title: "Task 1", Risk: "low", Work: map[string]float32{"be": 2, "fe": 5}},
			{Category: "API", Title: "Some Task 2 looooooooooong", Risk: "high", Work: map[string]float32{"be": 5, "fe": 1}},
			{Category: "API", Title: "User (FI users, API users & end-users) management", Risk: "high", Work: map[string]float32{"be": 5, "fe": 1}},
			{Category: "API", Title: "Some Task 3", Work: map[string]float32{"be": 3, "fe": 3}},
			{Category: "API", Title: "Some Task 3", Work: map[string]float32{"be": 3, "fe": 3}},
		},
	}
	fmt.Println(project)
	estimatorium.GenerateExcel(project, "Book2.xlsx")
}
