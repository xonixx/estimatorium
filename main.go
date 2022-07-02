package main

import (
	"estimatorium/estimatorium"
	"fmt"
)

func main() {
	/*f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet2")
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}*/

	project := estimatorium.Project{
		TimeUnit: estimatorium.Day,
		Currency: estimatorium.Usd,
		Risks:    estimatorium.StandardRisks(),
		Team: map[string]estimatorium.Resource{
			"fe": {
				Title: "Front dev", Rate: 40, Count: 1,
			},
			"be": {
				Title: "Back dev", Rate: 50, Count: 2,
			},
		},
		Tasks: []estimatorium.Task{
			{Category: "Initial", Title: "Task 1", Risk: "low", Work: map[string]float32{"be": 2, "fe": 5}},
			{Category: "API", Title: "Some Task 2", Risk: "high", Work: map[string]float32{"be": 5, "fe": 1}},
			{Category: "API", Title: "Some Task 3", Work: map[string]float32{"be": 3, "fe": 3}},
		},
	}
	fmt.Println(project)
}
