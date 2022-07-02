package estimatorium

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type excelGenerator struct {
	y, x  int    // current pos 0-based
	sheet string // current sheet
	f     *excelize.File
}

func (exc *excelGenerator) next() {
	exc.x++
}
func (exc *excelGenerator) cr() {
	exc.y++
	exc.x = 0
}
func (exc *excelGenerator) setVal(val interface{}) {
	err := exc.f.SetCellValue(exc.sheet, exc.cellName(), val)
	if err != nil {
		panic(err)
	}
}

func (exc *excelGenerator) cellName(abs ...bool) string {
	name, err := excelize.CoordinatesToCellName(exc.y+1, exc.x+1, abs...)
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
	exc.cr()
	exc.next()
	exc.setVal(100)
	if err := exc.f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
}
