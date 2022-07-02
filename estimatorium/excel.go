package estimatorium

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
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

func (exc *excelGenerator) cellName( /*abs ...bool*/ ) string {
	name, err := excelize.CoordinatesToCellName(exc.colZ+1, exc.rowZ+1 /*, abs...*/)
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

	startCatCell := ""
	endCatCell := ""
	currCat := ""

	for i, t := range project.Tasks {
		exc.setVal(t.Category)

		if i == 0 {
			startCatCell = exc.cellName()
			endCatCell = exc.cellName()
			currCat = t.Category
		} else if currCat != t.Category {
			fmt.Printf("merging: %s, %s\n", startCatCell, endCatCell)
			err := exc.f.MergeCell(exc.sheet, startCatCell, endCatCell)
			if err != nil {
				panic(err)
			}
			currCat = t.Category
			startCatCell = exc.cellName()
			endCatCell = exc.cellName()
		} else {
			endCatCell = exc.cellName()
		}

		exc.next()

		exc.setVal(t.Title)
		c1 := exc.cellName()
		exc.next()
		exc.f.MergeCell(exc.sheet, c1, exc.cellName())
		exc.next()
		v := map[string]string{}
		for _, r := range project.Team {
			v[r.Id] = exc.cellName()
			exc.setVal(t.Work[r.Id])
			exc.next()
		}
		riskCell := exc.cellName()
		exc.setVal(t.Risk)
		exc.next()
		for _, r := range project.Team {
			//exc.setVal(t.Work[r.Id]) // TODO
			err := exc.f.SetCellFormula(exc.sheet, exc.cellName(), risksFormula(project.Risks, v[r.Id], riskCell))
			//fmt.Println(exc.f.GetCellFormula(exc.sheet, exc.cellName()))
			if err != nil {
				panic(err)
			}
			exc.next()
		}
		exc.cr()
	}
	fmt.Printf("merging: %s, %s\n", startCatCell, endCatCell)
	err := exc.f.MergeCell(exc.sheet, startCatCell, endCatCell)
	if err != nil {
		panic(err)
	}
}

func risksFormula(risks map[string]float32, valCell string, risksCell string) string {
	// =ROUNDUP(D6*SWITCH($F6,"",1, "Low", 1.1, "Medium", 1.5, "High", 2, "Extreme", 5))
	var sb strings.Builder
	sb.WriteString("=ROUNDUP(")
	sb.WriteString(valCell)
	sb.WriteString("*_xlfn.SWITCH(")
	sb.WriteString(risksCell)
	sb.WriteString(",\"\",1")
	for k, v := range risks {
		sb.WriteString(",\"")
		sb.WriteString(k)
		sb.WriteString("\",")
		sb.WriteString(fmt.Sprintf("%f", v))
	}
	sb.WriteString("))")
	//fmt.Println(sb.String())
	return sb.String()
}

func generateTasksTableHeader(exc *excelGenerator, project Project) {
	cell0 := exc.cellName()

	exc.setVal("Feature")
	exc.next()
	exc.setVal("Story")
	c1 := exc.cellName()
	exc.next()
	exc.f.MergeCell(exc.sheet, c1, exc.cellName())
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
