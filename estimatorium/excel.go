package estimatorium

// TODO auto-fit to width https://github.com/qax-os/excelize/issues/92
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
func (exc *excelGenerator) setValAndNext(val interface{}) {
	exc.setVal(val)
	exc.next()
}
func (exc *excelGenerator) setFormulaAndNext(formula string) {
	checkErr(exc.f.SetCellFormula(exc.sheet, exc.cellName(), formula))
	//fmt.Println(exc.f.GetCellFormula(exc.sheet, exc.cellName()))
	exc.next()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (exc *excelGenerator) mergeNext(mergeCnt int) {
	cell0 := exc.cellName()
	for i := 0; i < mergeCnt; i++ {
		exc.next()
	}
	checkErr(exc.f.MergeCell(exc.sheet, cell0, exc.cellName()))
}

func (exc *excelGenerator) cellName( /*abs ...bool*/ ) string {
	name, err := excelize.CoordinatesToCellName(exc.colZ+1, exc.rowZ+1 /*, abs...*/)
	if err != nil {
		panic(err)
	}
	return name
}
func (exc *excelGenerator) setCellStyle(hCell, vCell string, styleId int) {
	checkErr(exc.f.SetCellStyle(exc.sheet, hCell, vCell, styleId))
}

func newExcelGenerator() *excelGenerator {
	return &excelGenerator{f: excelize.NewFile(), sheet: "Sheet1"}
}

func GenerateExcel(project Project, fileName string) {
	exc := newExcelGenerator()
	taskTableInfo := generateTasksTable(exc, project)
	exc.cr()
	generateCostsTable(exc, project, taskTableInfo)
	//exc.cr()
	//exc.next()
	//exc.setVal(100)
	if err := exc.f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
}

type tasksTableInfo struct {
	cellRanges         map[string]*cellRange
	cellRangesWithRisk map[string]*cellRange
}

func generateTasksTable(exc *excelGenerator, project Project) tasksTableInfo {
	generateTasksTableHeader(exc, project)

	res := tasksTableInfo{cellRanges: map[string]*cellRange{}, cellRangesWithRisk: map[string]*cellRange{}}

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
			checkErr(exc.f.MergeCell(exc.sheet, startCatCell, endCatCell))
			currCat = t.Category
			startCatCell = exc.cellName()
			endCatCell = exc.cellName()
		} else {
			endCatCell = exc.cellName()
		}

		exc.next()

		exc.setVal(t.Title)
		exc.mergeNext(1)
		exc.next()
		v := map[string]string{}
		for _, r := range project.Team {
			if i == 0 {
				res.cellRanges[r.Id] = &cellRange{hCell: exc.cellName()}
			} else if i == len(project.Tasks)-1 {
				res.cellRanges[r.Id].vCell = exc.cellName()
			}
			v[r.Id] = exc.cellName()
			exc.setValAndNext(t.Work[r.Id])
		}
		riskCell := exc.cellName()
		exc.setValAndNext(t.Risk)
		for _, r := range project.Team {
			if i == 0 {
				res.cellRangesWithRisk[r.Id] = &cellRange{hCell: exc.cellName()}
			} else if i == len(project.Tasks)-1 {
				res.cellRangesWithRisk[r.Id].vCell = exc.cellName()
			}
			//exc.setVal(t.Work[r.Id]) // TODO
			exc.setFormulaAndNext(risksFormula(project.Risks, v[r.Id], riskCell))
			//fmt.Println(exc.f.GetCellFormula(exc.sheet, exc.cellName()))
		}
		exc.cr()
	}
	fmt.Printf("merging: %s, %s\n", startCatCell, endCatCell)
	checkErr(exc.f.MergeCell(exc.sheet, startCatCell, endCatCell))

	return res
}

func generateCostsTable(exc *excelGenerator, project Project, tasksTableInfo tasksTableInfo) {
	generateCostsTableHeader(exc, project)
	for _, r := range project.Team {
		exc.setCellStyle(exc.cellName(), exc.cellName(), headerStyle(exc))
		exc.setValAndNext(r.Title)
		exc.setFormulaAndNext(tasksTableInfo.cellRanges[r.Id].sumFormula())
		exc.setFormulaAndNext(tasksTableInfo.cellRangesWithRisk[r.Id].sumFormula())
		exc.setValAndNext(r.Rate) // TODO fmt $
		exc.setValAndNext(r.Count)
		exc.setValAndNext("TODO")
		exc.cr()
	}
}

func generateCostsTableHeader(exc *excelGenerator, project Project) {
	generateHeader(exc, []headerCell{
		{title: ""},
		{title: fmt.Sprintf("Efforts (%v)", project.TimeUnit)},
		{title: "With Risk"},
		{title: "Rate"},
		{title: "Team"},
		{title: "Total"},
	})
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
	cols := []headerCell{
		{title: "Feature"},
		{title: "Story", mergedCells: 1},
	}

	for _, r := range project.Team {
		cols = append(cols, headerCell{title: r.Title})
	}

	cols = append(cols, headerCell{title: "Risks"})

	for _, r := range project.Team {
		cols = append(cols, headerCell{title: r.Title})
	}

	generateHeader(exc, cols)
}

type cellRange struct {
	hCell, vCell string
}

func (cellRange cellRange) sumFormula() string {
	return fmt.Sprintf("=SUM(%s:%s)", cellRange.hCell, cellRange.vCell)
}

type headerCell struct {
	mergedCells int
	title       string
}

func generateHeader(exc *excelGenerator, columns []headerCell) {
	exc.generateHeader(headerStyle(exc), columns)
}
func headerStyle(exc *excelGenerator) int {
	styleId, err := exc.f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#ffffff"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#091e42"}},
	})
	checkErr(err)
	return styleId
}
func (exc *excelGenerator) generateHeader(styleId int, columns []headerCell) {
	cell0 := exc.cellName()
	for _, col := range columns {
		exc.setVal(col.title)
		if col.mergedCells > 0 {
			exc.mergeNext(col.mergedCells)
		}
		exc.next()
	}
	exc.prev()
	cell1 := exc.cellName()
	exc.setCellStyle(cell0, cell1, styleId)
	exc.cr()
}
