package estimatorium

// TODO auto-fit to width https://github.com/qax-os/excelize/issues/92
import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"math"
	"strings"
	"unicode/utf8"
)

type excelGenerator struct {
	colZ, rowZ          int    // current pos 0-based
	sheet               string // current sheet
	f                   *excelize.File
	currencyStyleId     int
	currencyBoldStyleId int
	headerStyleId       int
	valueStyleId        int
	taskNameStyleId     int
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
func (exc *excelGenerator) setVal(val interface{}, styles ...int) {
	checkErr(exc.f.SetCellStyle(exc.sheet, exc.currentCell(), exc.currentCell(), getCellStyle(exc, styles)))
	checkErr(exc.f.SetCellValue(exc.sheet, exc.currentCell(), val))
}
func getCellStyle(exc *excelGenerator, styles []int) int {
	if styles == nil {
		styles = []int{exc.valueStyleId}
	}
	return styles[0]
}
func (exc *excelGenerator) setValAndNext(val interface{}, styles ...int) {
	exc.setVal(val, styles...)
	exc.next()
}
func (exc *excelGenerator) setFormulaAndNext(formula string, styles ...int) {
	checkErr(exc.f.SetCellStyle(exc.sheet, exc.currentCell(), exc.currentCell(), getCellStyle(exc, styles)))
	checkErr(exc.f.SetCellFormula(exc.sheet, exc.currentCell(), formula))
	//fmt.Println(exc.f.GetCellFormula(exc.sheet, exc.currentCell()))
	exc.next()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (exc *excelGenerator) mergeNext(mergeCnt int) {
	cell0 := exc.currentCell()
	for i := 0; i < mergeCnt; i++ {
		exc.next()
	}
	checkErr(exc.f.MergeCell(exc.sheet, cell0, exc.currentCell()))
}

func (exc *excelGenerator) currentCell() string {
	return exc.currentCellAbs(false)
}
func (exc *excelGenerator) currentCellAbs(abs bool) string {
	name, err := excelize.CoordinatesToCellName(exc.colZ+1, exc.rowZ+1, abs)
	checkErr(err)
	return name
}

func newStyle(file *excelize.File, style *excelize.Style) int {
	styleId, err := file.NewStyle(style)
	checkErr(err)
	return styleId
}
func newExcelGenerator(currency Currency) *excelGenerator {
	file := excelize.NewFile()
	fmtCode := "[$" + currency.Symbol() + "]#,##0"
	borders := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}
	return &excelGenerator{f: file, sheet: "Sheet1",
		currencyStyleId:     newStyle(file, &excelize.Style{CustomNumFmt: &fmtCode, Border: borders}),
		currencyBoldStyleId: newStyle(file, &excelize.Style{CustomNumFmt: &fmtCode, Border: borders, Font: &excelize.Font{Bold: true}}),
		valueStyleId:        newStyle(file, &excelize.Style{Border: borders}),
		taskNameStyleId:     newStyle(file, &excelize.Style{Border: borders, Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#93c47d"}}}),
		headerStyleId: newStyle(file, &excelize.Style{
			Font:      &excelize.Font{Bold: true, Color: "#ffffff"},
			Alignment: &excelize.Alignment{Horizontal: "center"},
			Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#091e42"}},
		})}
}

func GenerateExcel(project Project, fileName string) {
	exc := newExcelGenerator(project.Currency)
	taskTableInfo := generateTasksTable(exc, project)
	exc.cr()
	parametersTableInfo := parametersTableInfo{}
	if project.AcceptancePercent > 0 {
		parametersTableInfo = generateParametersTable(exc, project.AcceptancePercent)
		exc.cr()
	}
	costsTableInfo := generateCostsTable(exc, project, taskTableInfo, parametersTableInfo)
	exc.cr()
	generateDurationsTable(exc, project, costsTableInfo)

	autoFixColWidths(exc)

	checkErr(exc.f.SaveAs(fileName))
}

const ColTotal = "Total"

func autoFixColWidths(exc *excelGenerator) {
	// Autofit all columns according to their text content
	cols, err := exc.f.GetCols(exc.sheet)
	checkErr(err)
	largestWidthMap := map[int]float64{}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
			if rowCell == ColTotal {
				largestWidth = int(math.Max(float64(largestWidth), 11))
			}
		}
		largestWidthMap[idx] = float64(largestWidth)
	}

	// hack for merged
	if largestWidthMap[1] > 17 {
		largestWidthMap[1] -= largestWidthMap[2]
		largestWidthMap[1] -= 1
	}
	largestWidthMap[1] = (largestWidthMap[1] + largestWidthMap[2]) / 2
	largestWidthMap[2] = largestWidthMap[1]

	for idx := range cols {
		name, err := excelize.ColumnNumberToName(idx + 1)
		checkErr(err)
		checkErr(exc.f.SetColWidth(exc.sheet, name, name, largestWidthMap[idx]))
	}
}

type tasksTableInfo struct {
	cellRanges         map[string]*cellRange
	cellRangesWithRisk map[string]*cellRange
}

func generateTasksTable(exc *excelGenerator, project Project) tasksTableInfo {
	generateTasksTableHeader(exc, project)

	res := tasksTableInfo{cellRanges: map[string]*cellRange{}, cellRangesWithRisk: map[string]*cellRange{}}

	riskLabels := RiskLabels(project.Risks)

	startCatCell := ""
	endCatCell := ""
	currCat := ""

	for i, t := range project.Tasks {
		exc.setVal(t.Category)

		if i == 0 {
			startCatCell = exc.currentCell()
			endCatCell = exc.currentCell()
			currCat = t.Category
		} else if currCat != t.Category {
			fmt.Printf("merging: %s, %s\n", startCatCell, endCatCell)
			checkErr(exc.f.MergeCell(exc.sheet, startCatCell, endCatCell))
			currCat = t.Category
			startCatCell = exc.currentCell()
			endCatCell = exc.currentCell()
		} else {
			endCatCell = exc.currentCell()
		}

		exc.next()

		exc.setVal(t.Title, exc.taskNameStyleId)
		exc.mergeNext(1)
		exc.next()
		v := map[string]string{}
		teamExcludingDerived := project.TeamExcludingDerived()
		for _, r := range teamExcludingDerived {
			if i == 0 {
				res.cellRanges[r.Id] = &cellRange{hCell: exc.currentCell()}
			} else if i == len(project.Tasks)-1 {
				res.cellRanges[r.Id].vCell = exc.currentCell()
			}
			v[r.Id] = exc.currentCell()
			exc.setValAndNext(t.Work[r.Id])
		}
		riskCell := exc.currentCell()

		dv := excelize.NewDataValidation(true)
		dv.Sqref = riskCell + ":" + riskCell
		checkErr(dv.SetDropList(riskLabels))
		checkErr(exc.f.AddDataValidation(exc.sheet, dv))

		exc.setValAndNext(t.Risk)
		for _, r := range teamExcludingDerived {
			if i == 0 {
				res.cellRangesWithRisk[r.Id] = &cellRange{hCell: exc.currentCell()}
			} else if i == len(project.Tasks)-1 {
				res.cellRangesWithRisk[r.Id].vCell = exc.currentCell()
			}
			//exc.setVal(t.Work[r.Id]) // TODO
			exc.setFormulaAndNext(risksFormula(project.Risks, v[r.Id], riskCell))
			//fmt.Println(exc.f.GetCellFormula(exc.sheet, exc.currentCell()))
		}
		exc.cr()
	}
	fmt.Printf("merging: %s, %s\n", startCatCell, endCatCell)
	checkErr(exc.f.MergeCell(exc.sheet, startCatCell, endCatCell))

	return res
}

type parametersTableInfo struct {
	acceptancePercentCell string
}

func generateParametersTable(exc *excelGenerator, acceptancePercent float32) parametersTableInfo {
	res := parametersTableInfo{}
	exc.setValAndNext("Cleanup & acceptance")
	res.acceptancePercentCell = exc.currentCellAbs(true)
	exc.setValAndNext(fmt.Sprintf("%.1f%%", acceptancePercent))
	exc.cr()
	return res
}

type resourceCostsCells struct {
	effortsCell, effortsWithRisksCell, countCell string
}

type costsTableInfo struct {
	costsData map[string]*resourceCostsCells
}

func generateCostsTable(exc *excelGenerator, project Project, tasksTableInfo tasksTableInfo, parametersTableInfo parametersTableInfo) costsTableInfo {
	res := costsTableInfo{costsData: map[string]*resourceCostsCells{}}
	generateCostsTableHeader(exc, project)
	effortsRange := cellRange{}
	effortsWithRiskRange := cellRange{}
	totalsRange := cellRange{}
	for i, r := range project.Team {
		exc.setValAndNext(r.Title, exc.headerStyleId)
		isFirst := i == 0
		isLast := i == len(project.Team)-1
		if isFirst {
			effortsRange.hCell = exc.currentCell()
		} else if isLast {
			effortsRange.vCell = exc.currentCell()
		}
		res.costsData[r.Id] = &resourceCostsCells{effortsCell: exc.currentCell()}
		var effortsFormula string
		if r.Formula == "" {
			effortsFormula = tasksTableInfo.cellRanges[r.Id].sumFormula()
		} else {
			formula := r.Formula
			for _, r1 := range project.TeamExcludingDerived() {
				formula = strings.Replace(formula, r1.Id, "SUM("+tasksTableInfo.cellRanges[r1.Id].String()+")", -1)
			}
			effortsFormula = formula
		}
		if project.AcceptancePercent > 0 {
			effortsFormula += "*(1+" + parametersTableInfo.acceptancePercentCell + ")"
		}
		exc.setFormulaAndNext(effortsFormula)
		if isFirst {
			effortsWithRiskRange.hCell = exc.currentCell()
		} else if isLast {
			effortsWithRiskRange.vCell = exc.currentCell()
		}
		effortsWithRisksCell := exc.currentCell()
		res.costsData[r.Id].effortsWithRisksCell = effortsWithRisksCell
		var effortsWithRisksFormula string
		if r.Formula == "" {
			effortsWithRisksFormula = tasksTableInfo.cellRangesWithRisk[r.Id].sumFormula()
		} else {
			formula := r.Formula
			for _, r1 := range project.TeamExcludingDerived() {
				formula = strings.Replace(formula, r1.Id, "SUM("+tasksTableInfo.cellRangesWithRisk[r1.Id].String()+")", -1)
			}
			effortsWithRisksFormula = formula
		}
		if project.AcceptancePercent > 0 {
			effortsWithRisksFormula += "*(1+" + parametersTableInfo.acceptancePercentCell + ")"
		}
		exc.setFormulaAndNext(effortsWithRisksFormula)
		rateCell := exc.currentCell()
		exc.setValAndNext(r.Rate, exc.currencyStyleId)
		res.costsData[r.Id].countCell = exc.currentCell()
		exc.setValAndNext(r.Count)
		if isFirst {
			totalsRange.hCell = exc.currentCell()
		} else if isLast {
			totalsRange.vCell = exc.currentCell()
		}
		exc.setFormulaAndNext(fmt.Sprintf("%d*%s*%s",
			project.TimeUnit.ToHours(), effortsWithRisksCell, rateCell), exc.currencyStyleId)
		exc.cr()
	}
	exc.setValAndNext("Sum", exc.headerStyleId)
	exc.setFormulaAndNext(effortsRange.sumFormula())
	exc.setFormulaAndNext(effortsWithRiskRange.sumFormula())
	exc.setValAndNext("")
	exc.setValAndNext("")
	exc.setFormulaAndNext(totalsRange.sumFormula(), exc.currencyBoldStyleId)
	exc.cr()
	return res
}

func generateCostsTableHeader(exc *excelGenerator, project Project) {
	generateHeader(exc, []headerCell{
		{title: ""},
		{title: fmt.Sprintf("Efforts (%v)", project.TimeUnit)},
		{title: "With Risk"},
		{title: "Rate"},
		{title: "Team"},
		{title: ColTotal},
	})
}

func generateDurationsTable(exc *excelGenerator, project Project, costsTableInfo costsTableInfo) {
	generateDurationsTableHeader(exc)

	exc.setValAndNext("Duration", exc.headerStyleId)

	exc.setFormulaAndNext(durationFormula(project, costsTableInfo, func(cells *resourceCostsCells) string {
		return cells.effortsCell
	}))
	exc.setValAndNext("Months")
	exc.cr()
	exc.setValAndNext("With risks", exc.headerStyleId)
	exc.setFormulaAndNext(durationFormula(project, costsTableInfo, func(cells *resourceCostsCells) string {
		return cells.effortsWithRisksCell
	}))
	exc.setValAndNext("Months")
	exc.cr()
}

func durationFormula(project Project, costsTableInfo costsTableInfo, f func(*resourceCostsCells) string) string {
	var sb strings.Builder
	sb.WriteString("ROUND(MAX(")
	resources := project.TeamExcludingDerived()
	for i, r := range resources {
		cells := costsTableInfo.costsData[r.Id]
		sb.WriteString(f(cells))
		sb.WriteString("/")
		sb.WriteString(cells.countCell)
		if i < len(resources)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")*")
	sb.WriteString(fmt.Sprintf("%d", project.TimeUnit.ToHours()))
	sb.WriteString("/8/21,1)")
	return sb.String()
}

func generateDurationsTableHeader(exc *excelGenerator) {
	generateHeader(exc, []headerCell{
		{title: ""},
		{title: "Timeframe draft", mergedCells: 1},
	})
}

func risksFormula(risks map[string]float32, valCell string, risksCell string) string {
	// =ROUNDUP(D6*SWITCH($F6,"",1, "Low", 1.1, "Medium", 1.5, "High", 2, "Extreme", 5))
	var sb strings.Builder
	sb.WriteString("ROUNDUP(")
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

	teamExcludingDerived := project.TeamExcludingDerived()

	for _, r := range teamExcludingDerived {
		cols = append(cols, headerCell{title: r.Title})
	}

	cols = append(cols, headerCell{title: "Risks"})

	for _, r := range teamExcludingDerived {
		cols = append(cols, headerCell{title: r.Title})
	}

	generateHeader(exc, cols)
}

type cellRange struct {
	hCell, vCell string
}

func (cellRange cellRange) String() string {
	return cellRange.hCell + ":" + cellRange.vCell
}
func (cellRange cellRange) sumFormula() string {
	return fmt.Sprintf("SUM(%s)", cellRange)
}

type headerCell struct {
	mergedCells int
	title       string
}

func generateHeader(exc *excelGenerator, columns []headerCell) {
	exc.generateHeader(exc.headerStyleId, columns)
}
func (exc *excelGenerator) generateHeader(styleId int, columns []headerCell) {
	for _, col := range columns {
		exc.setVal(col.title, styleId)
		if col.mergedCells > 0 {
			exc.mergeNext(col.mergedCells)
		}
		exc.next()
	}
	exc.cr()
}
