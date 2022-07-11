package estimatorium

import (
	"regexp"
	"strconv"
	"strings"
)

// TODO validate mandatory directives present
// TODO validate rate(s) absent for resources in tasks
// TODO validate that derived resources are not referred in tasks
// TODO how can we set custom team resources (w titles)

type directiveVals struct {
	value  string
	values map[string]string
	directiveDef
}

type projParsed struct {
	directives   map[string]directiveVals // each directive can go at most one time
	team         []resourceRecord
	tasksRecords []taskRecord
}

type resourceRecord struct {
	id      string
	cnt     int
	rate    float64
	title   string
	formula string
}

type taskRecord struct {
	category string
	title    string
	efforts  map[string]float64
	risk     string
}

func (projParsed projParsed) getSingleVal(directive directiveDef) *string {
	v, exists := projParsed.directives[directive.name]
	if !exists {
		return nil
	}
	if v.directiveType == DtSingleValue {
		return &v.value
	} else {
		panic(v.directiveType)
	}
}

func (projParsed projParsed) getKVPairs(directive directiveDef) *map[string]string {
	v, exists := projParsed.directives[directive.name]
	if !exists {
		return nil
	}
	if v.directiveType == DtKeyVal {
		return &v.values
	} else {
		panic(v.directiveType)
	}
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

var (
	directiveProject           = newDirectiveDef("project", DtSingleValue)
	directiveAuthor            = newDirectiveDef("author", DtSingleValue)
	directiveCurrency          = newDirectiveDef("currency", DtSingleValue)
	directiveTimeUnit          = newDirectiveDef("time_unit", DtSingleValue)
	directiveAcceptancePercent = newDirectiveDef("acceptance_percent", DtSingleValue)
	directiveRisks             = newDirectiveDef("risks", DtKeyVal)
	directiveRates             = newDirectiveDef("rates", DtKeyVal)
	directiveFormula           = newDirectiveDef("formula", DtKeyVal)
	directiveDesiredDuration   = newDirectiveDef("desired_duration", DtSingleValue)
	directiveTeam              = newDirectiveDef("team", DtKeyVal)
)

var directives = map[string]directiveDef{}

func newDirectiveDef(name string, directiveType directiveType) directiveDef {
	d := directiveDef{
		name:          name,
		directiveType: directiveType,
	}
	directives[name] = d
	return d
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
		name := projParsed.getSingleVal(directiveProject)
		if name != nil {
			proj.Name = *name
		}
	}
	{
		author := projParsed.getSingleVal(directiveAuthor)
		if author != nil {
			proj.Author = *author
		}
	}

	{
		timeUnit := projParsed.getSingleVal(directiveTimeUnit)
		if timeUnit != nil {
			proj.TimeUnit = TimeUnitFromString(*timeUnit)
			if proj.TimeUnit == TimeUnitUnknown {
				errors.addError("Unknown time_unit: " + *timeUnit)
			}
		}
	}

	{
		currency := projParsed.getSingleVal(directiveCurrency)
		if currency != nil {
			proj.Currency = CurrencyFromString(*currency)
			if proj.Currency == CurrencyUnknown {
				errors.addError("Unknown currency: " + *currency)
			}
		}
	}

	{
		acceptancePercent := projParsed.getSingleVal(directiveAcceptancePercent)
		if acceptancePercent != nil {
			float, err := strconv.ParseFloat(*acceptancePercent, 32)
			if err != nil || float < 0 || float > 100 {
				errors.addError("Wrong acceptance_percent: " + *acceptancePercent)
			}
			proj.AcceptancePercent = float64(float)
		}
	}

	{
		risks := projParsed.getKVPairs(directiveRisks)
		if risks != nil {
			proj.Risks = map[string]float64{}
			for k, v := range *risks {
				float, err := strconv.ParseFloat(v, 32)
				if err != nil || float < 1 {
					errors.addError("Wrong risk value for " + k + ": " + v)
				}
				proj.Risks[k] = float64(float)
			}
		} else {
			// apply default risks
			proj.Risks = StandardRisks()
		}
	}

	ratesM := map[string]float64{}
	formulaM := map[string]string{}
	teamM := map[string]int{}

	{
		rates := projParsed.getKVPairs(directiveRates)
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
		formula := projParsed.getKVPairs(directiveFormula)
		if formula != nil {
			for k, v := range *formula {
				formulaM[k] = v
			}
		}
	}

	{
		team := projParsed.getKVPairs(directiveTeam)
		if team != nil {
			for k, v := range *team {
				intVal, err := strconv.ParseInt(v, 10, 32)
				if err != nil || intVal < 0 {
					errors.addError("Wrong team count value for " + k + ": " + v)
				}
				teamM[k] = int(intVal)
			}
		} else {
			// use standard team
			for r, _ := range standardResourceTypes {
				teamM[r] = 1
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
		risk := tasksRecord.risk
		if risk != "" {
			if _, exists := proj.Risks[risk]; !exists {
				errors.addError("Wrong risks name: " + risk)
			}
		}
		efforts := tasksRecord.efforts
		for k := range efforts {
			if proj.ResourceById(k) == nil {
				errors.addError("Wrong resource name in efforts: " + k)
			}
		}
		proj.Tasks = append(proj.Tasks, Task{
			Category: tasksRecord.category,
			Title:    tasksRecord.title,
			Risk:     risk,
			Work:     efforts,
		})
	}

	{
		desiredDurationStr := projParsed.getSingleVal(directiveDesiredDuration)
		if desiredDurationStr != nil {
			duration, err := ParseDuration(*desiredDurationStr)
			if err != nil {
				errors.addError("Unable to parse desired duration: " + err.Error())
			}
			proj.DesiredDuration = duration
		}
	}

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
			if directive, found := directives[parts[0]]; found {
				if _, exists := projParsed.directives[directive.name]; exists {
					errors.addError("Duplicating directive: " + directive.name)
				} else if directive.directiveType == DtSingleValue {
					projParsed.directives[directive.name] = directiveVals{directiveDef: directive, value: strings.TrimSpace(parts[1])}
				} else if directive.directiveType == DtKeyVal {
					projParsed.directives[directive.name] = directiveVals{directiveDef: directive, values: parseKeyValPairs(parts[1])}
				}
			} else {
				errors.addError("Unknown directive: " + parts[0])
			}
		} else if mode == pmTasks {
			taskParts := strings.Split(line, "|")
			if len(taskParts) != 3 {
				panic("task should have format: cat | title | efforts") // TODO convert to error
			}
			keyValPairs := parseKeyValPairs(taskParts[2])
			efforts := map[string]float64{}
			for k, v := range keyValPairs {
				if k != "risks" {
					float, err := strconv.ParseFloat(v, 32)
					checkErr(err)
					efforts[k] = float64(float)
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
