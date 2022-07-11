package estimatorium

import (
	"fmt"
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
	id            string
	resourceProps map[string]string
}

type taskRecord struct {
	category  string
	title     string
	taskProps map[string]string
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
	directiveDesiredDuration   = newDirectiveDef("desired_duration", DtSingleValue)
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
func (ppe *ProjectParseError) addError(error string) {
	ppe.errors = append(ppe.errors, error)
}
func (ppe *ProjectParseError) addErrorf(errorF string, args ...any) {
	ppe.addError(fmt.Sprintf(errorF, args...))
}
func (ppe *ProjectParseError) intOrAddError(v string, errorF string, args ...any) int {
	intVal, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		ppe.addErrorf(errorF, args...)
	}
	return int(intVal)
}
func (ppe *ProjectParseError) floatOrAddErrorf(v string, errorF string, args ...any) float64 {
	float, err := strconv.ParseFloat(v, 32)
	if err != nil {
		ppe.addErrorf(errorF, args...)
	}
	return float
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

const risksKey = "risks"

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
			proj.AcceptancePercent = float
		}
	}

	{
		risks := projParsed.getKVPairs(directiveRisks)
		if risks != nil {
			proj.Risks = map[string]float64{}
			for k, v := range *risks {
				float, err := strconv.ParseFloat(v, 32)
				if err != nil || float < 1 {
					errors.addErrorf("Wrong risk value for %s: %s", k, v)
				}
				proj.Risks[k] = float
			}
		} else {
			// apply default risks
			proj.Risks = StandardRisks()
		}
	}

	for _, r := range projParsed.team {
		title := r.resourceProps["title"]
		resourceId := r.id
		if title == "" {
			title = standardResourceTypes[resourceId]
		}
		var cnt int
		if cntStr, exists := r.resourceProps["cnt"]; exists {
			cnt = errors.intOrAddError(cntStr, "Wrong team count value for %s: %s", resourceId, cntStr)
			if cnt < 0 {
				errors.addError(fmt.Sprintf("Team count must be >= 0 for %s: %s", resourceId, cntStr))
			}
		}
		var rate float64
		if rateStr, exists := r.resourceProps["rate"]; exists {
			rate = errors.floatOrAddErrorf(rateStr, "Wrong rate value for %s: %s", resourceId, rateStr)
			if rate < 0 {
				errors.addError(fmt.Sprintf("Rate must be >= 0 for %s: %s", resourceId, rateStr))
			}
		}
		proj.Team = append(proj.Team, Resource{
			Id:      resourceId,
			Title:   title,
			Rate:    rate,
			Count:   cnt,
			Formula: r.resourceProps["formula"],
		})
	}

	for _, taskRecord := range projParsed.tasksRecords {
		risk := taskRecord.taskProps[risksKey]
		if risk != "" {
			if _, exists := proj.Risks[risk]; !exists {
				errors.addError("Wrong risks name: " + risk)
			}
		}
		efforts := map[string]float64{}
		for k, v := range taskRecord.taskProps {
			if k != risksKey {
				effort := errors.floatOrAddErrorf(v, "Wrong effort for task %s|%s for resource %s: %s", taskRecord.category, taskRecord.title, k, v)
				if effort < 0 {
					errors.addErrorf("Effort should be >= 0 for task %s|%s for resource %s: %s", taskRecord.category, taskRecord.title, k, v)
				}
				efforts[k] = effort
			}
		}
		for k := range efforts {
			if proj.ResourceById(k) == nil {
				errors.addError("Wrong resource name in efforts: " + k)
			}
		}
		proj.Tasks = append(proj.Tasks, Task{
			Category: taskRecord.category,
			Title:    taskRecord.title,
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
	pmTeam
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
		} else if line == "tasks" {
			mode = pmTasks
			continue
		} else if line == "team" {
			mode = pmTeam
			continue
		}
		parts := spaceRe.Split(line, 2)
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
			projParsed.tasksRecords = append(projParsed.tasksRecords, taskRecord{
				category:  strings.TrimSpace(taskParts[0]),
				title:     strings.TrimSpace(taskParts[1]),
				taskProps: parseKeyValPairs(taskParts[2]),
			})
		} else if mode == pmTeam {
			keyValPairs := map[string]string{}
			if len(parts) > 1 {
				keyValPairs = parseKeyValPairs(parts[1])
			}
			projParsed.team = append(projParsed.team, resourceRecord{
				id:            parts[0],
				resourceProps: keyValPairs,
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
