package estimatorium

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type excelGenerator struct {
	colZ, rowZ int    // current pos 0-based
	sheet      string // current sheet
	f          *excelize.File
}

func (exc *excelGenerator) next() {
	exc.colZ++
}
func (exc *excelGenerator) prev() {
	exc.colZ--
}
func (exc *excelGenerator) cr() {
	exc.rowZ++
	exc.colZ = 0
}
func (exc *excelGenerator) setVal(val interface{}) {
	err := exc.f.SetCellValue(exc.sheet, exc.cellName(), val)
	if err != nil {
		panic(err)
	}
}

func (exc *excelGenerator) cellName(abs ...bool) string {
	name, err := excelize.CoordinatesToCellName(exc.colZ+1, exc.rowZ+1, abs...)
	if err != nil {
		panic(err)
	}
	return name
}

func newExcelGenerator() *excelGenerator {
	return &excelGenerator{f: excelize.NewFile(), sheet: "Sheet1"}
}

func GenerateExcel(project Project, fileName string) {
	exc := newExcelGenerator()
	generateTasksTable(exc, project)
	//exc.cr()
	//exc.next()
	//exc.setVal(100)
	if err := exc.f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
}

func generateTasksTable(exc *excelGenerator, project Project) {
	generateTasksTableHeader(exc, project)
}

func generateTasksTableHeader(exc *excelGenerator, project Project) {
	cell0 := exc.cellName()

	exc.setVal("Feature")
	exc.next()
	exc.setVal("Story")
	exc.next()
	for _, r := range project.Team {
		exc.setVal(r.Title)
		exc.next()
	}
	exc.setVal("Risks")
	exc.next()
	for _, r := range project.Team {
		exc.setVal(r.Title)
		exc.next()
	}
	exc.prev()
	cell1 := exc.cellName()
	style, _ := exc.f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#ffffff"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#091e42"}},
	})
	exc.f.SetCellStyle(exc.sheet, cell0, cell1, style)
	exc.cr()
}
