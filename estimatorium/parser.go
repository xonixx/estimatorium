package estimatorium

import (
	"regexp"
	"strconv"
	"strings"
)

// TODO validate mandatory directives present

type directiveVals struct {
	value  string
	values map[string]string
	directiveDef
}

type projParsed struct {
	directives   map[string]directiveVals // each directive can go at most one time
	tasksRecords []taskRecord
}

func (projParsed projParsed) getSingleVal(directive string) *string {
	v, exists := projParsed.directives[directive]
	if !exists {
		return nil
	}
	if v.directiveType == DtSingleValue {
		return &v.value
	} else {
		panic(v.directiveType)
	}
}
func (projParsed projParsed) getKVPairs(directive string) *map[string]string {
	v, exists := projParsed.directives[directive]
	if !exists {
		return nil
	}
	if v.directiveType == DtKeyVal {
		return &v.values
	} else {
		panic(v.directiveType)
	}
}

type taskRecord struct {
	category string
	title    string
	efforts  map[string]float32
	risk     string
}

type directiveType int8

const (
	DtSingleValue directiveType = iota
	DtKeyVal
)

type directiveDef struct {
	name          string
	directiveType directiveType
}

var directives = []directiveDef{
	{name: "project", directiveType: DtSingleValue},
	{name: "author", directiveType: DtSingleValue},
	{name: "currency", directiveType: DtSingleValue},
	{name: "time_unit", directiveType: DtSingleValue},
	{name: "acceptance_percent", directiveType: DtSingleValue},
	{name: "risks", directiveType: DtKeyVal},
	{name: "rates", directiveType: DtKeyVal},
	{name: "formula", directiveType: DtKeyVal},
	{name: "desired_duration", directiveType: DtSingleValue},
	{name: "team", directiveType: DtKeyVal},
}

var spaceRe = regexp.MustCompile("[ \t]+")

type ProjectParseError struct {
	errors []string
}

func (ppe *ProjectParseError) hasErrors() bool {
	return len(ppe.errors) > 0
}
func (ppe *ProjectParseError) addError(err string) {
	ppe.errors = append(ppe.errors, err)
}
func (ppe *ProjectParseError) Error() string {
	return strings.Join(ppe.errors, "\n")
}

func (ppe *ProjectParseError) addOtherError(err error) {
	if err == nil {
		return
	}
	if parseError, ok := err.(*ProjectParseError); ok {
		ppe.errors = append(ppe.errors, parseError.errors...)
	} else {
		ppe.addError(err.Error())
	}
}

func ProjectFromString(projData string) (Project, error) {
	proj := Project{}
	errors := &ProjectParseError{}

	projParsed, err := parseProj(projData)
	errors.addOtherError(err)

	{
		name := projParsed.getSingleVal("project")
		if name != nil {
			proj.Name = *name
		}
	}

	{
		timeUnit := projParsed.getSingleVal("time_unit")
		if timeUnit != nil {
			proj.TimeUnit = TimeUnitFromString(*timeUnit)
			if proj.TimeUnit == TimeUnitUnknown {
				errors.addError("Unknown time_unit: " + *timeUnit)
			}
		}
	}

	{
		currency := projParsed.getSingleVal("currency")
		if currency != nil {
			proj.Currency = CurrencyFromString(*currency)
			if proj.Currency == CurrencyUnknown {
				errors.addError("Unknown currency: " + *currency)
			}
		}
	}

	{
		acceptancePercent := projParsed.getSingleVal("acceptance_percent")
		if acceptancePercent != nil {
			float, err := strconv.ParseFloat(*acceptancePercent, 32)
			if err != nil || float < 0 || float > 100 {
				errors.addError("Wrong acceptance_percent: " + *acceptancePercent)
			}
			proj.AcceptancePercent = float32(float)
		}
	}

	{
		risks := projParsed.getKVPairs("risks")
		if risks != nil {
			proj.Risks = map[string]float32{}
			for k, v := range *risks {
				float, err := strconv.ParseFloat(v, 32)
				if err != nil || float < 1 {
					errors.addError("Wrong risk value for " + k + ": " + v)
				}
				proj.Risks[k] = float32(float)
			}
		}
	}

	ratesM := map[string]float64{}
	formulaM := map[string]string{}
	teamM := map[string]int{}

	{
		rates := projParsed.getKVPairs("rates")
		if rates != nil {
			for k, v := range *rates {
				float, err := strconv.ParseFloat(v, 32)
				if err != nil || float < 0 {
					errors.addError("Wrong rate value for " + k + ": " + v)
				}
				ratesM[k] = float
			}
		}
	}

	{
		formula := projParsed.getKVPairs("formula")
		if formula != nil {
			for k, v := range *formula {
				formulaM[k] = v
			}
		}
	}

	{
		team := projParsed.getKVPairs("team")
		if team != nil {
			for k, v := range *team {
				intVal, err := strconv.ParseInt(v, 10, 32)
				if err != nil || intVal < 0 {
					errors.addError("Wrong team count value for " + k + ": " + v)
				}
				teamM[k] = int(intVal)
			}
		}
	}

	resourcesM := map[string]*Resource{}
	for rId, cnt := range teamM {
		resourcesM[rId] = &Resource{Id: rId, Count: cnt, Title: standardResourceTypes[rId]}
	}

	for rId, rate := range ratesM {
		resourcesM[rId].Rate = rate
	}

	for rId, formula := range formulaM {
		resourcesM[rId].Formula = formula
	}

	for _, resource := range resourcesM {
		proj.Team = append(proj.Team, *resource)
	}

	for _, tasksRecord := range projParsed.tasksRecords {
		proj.Tasks = append(proj.Tasks, Task{
			Category: tasksRecord.category,
			Title:    tasksRecord.title,
			Risk:     tasksRecord.risk,
			Work:     tasksRecord.efforts,
		})
	}

	// TODO desired duration

	if !errors.hasErrors() {
		return proj, nil
	}
	return proj, errors
}

type parseMode int

const (
	pmDirectives parseMode = iota
	pmTasks
)

func parseKeyValPairs(str string) map[string]string {
	str = strings.TrimSpace(str)
	values := map[string]string{}
	valParts := spaceRe.Split(str, -1)
	for _, valPart := range valParts {
		keyVal := strings.SplitN(valPart, "=", 2)
		values[keyVal[0]] = keyVal[1]
	}
	return values
}

func parseProj(projData string) (projParsed, error) {
	errors := &ProjectParseError{}
	projParsed := projParsed{
		directives:   map[string]directiveVals{},
		tasksRecords: []taskRecord{},
	}
	lines := strings.Split(projData, "\n")
	mode := pmDirectives
linesLoop:
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Index(line, "#") == 0 {
			continue
		}
		parts := spaceRe.Split(line, 2)
		if parts[0] == "tasks" {
			mode = pmTasks
			continue
		}
		if mode == pmDirectives {
			for _, directive := range directives {
				if parts[0] == directive.name {
					if _, exists := projParsed.directives[directive.name]; exists {
						errors.addError("Duplicating directive: " + directive.name)
					} else if directive.directiveType == DtSingleValue {
						projParsed.directives[directive.name] = directiveVals{directiveDef: directive, value: strings.TrimSpace(parts[1])}
					} else if directive.directiveType == DtKeyVal {
						projParsed.directives[directive.name] = directiveVals{directiveDef: directive, values: parseKeyValPairs(parts[1])}
					}
					continue linesLoop
				}
			}
			errors.addError("Unknown directive: " + parts[0])
		} else if mode == pmTasks {
			taskParts := strings.Split(line, "|")
			if len(taskParts) != 3 {
				panic("task should have format: cat | title | efforts") // TODO convert to error
			}
			keyValPairs := parseKeyValPairs(taskParts[2])
			efforts := map[string]float32{}
			for k, v := range keyValPairs {
				if k != "risks" {
					float, err := strconv.ParseFloat(v, 32)
					checkErr(err)
					efforts[k] = float32(float)
				}
			}
			projParsed.tasksRecords = append(projParsed.tasksRecords, taskRecord{
				category: strings.TrimSpace(taskParts[0]),
				title:    strings.TrimSpace(taskParts[1]),
				efforts:  efforts,
				risk:     keyValPairs["risks"],
			})
		} else {
			panic("Unknown mode")
		}
	}
	if !errors.hasErrors() {
		return projParsed, nil
	}
	return projParsed, errors
}
