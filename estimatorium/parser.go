package estimatorium

import (
	"regexp"
	"strconv"
	"strings"
)

// 1. each directive can go at most one time

type directiveVals struct {
	value  string
	values map[string]string
}

type projParsed struct {
	directives   map[string]directiveVals // each directive can go at most one time
	tasksRecords []taskRecord
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
	{"project", DtSingleValue},
	{"author", DtSingleValue},
	{"currency", DtSingleValue},
	{"time_unit", DtSingleValue},
	{"acceptance_percent", DtSingleValue},
	{"risks", DtKeyVal},
	{"rates", DtKeyVal},
	{"formula", DtKeyVal},
	{"desired_duration", DtSingleValue},
	{"team", DtKeyVal},
}

var spaceRe = regexp.MustCompile("[ \t]+")

func ProjectFromString(projData string) Project {
	return Project{} // TODO
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

// TODO wrong directive error
// TODO handle comment explicitly
func parseProj(projData string) projParsed {
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
			for _, directive := range directives {
				if parts[0] == directive.name {
					if directive.directiveType == DtSingleValue {
						projParsed.directives[directive.name] = directiveVals{value: strings.TrimSpace(parts[1])}
					} else if directive.directiveType == DtKeyVal {
						projParsed.directives[directive.name] = directiveVals{values: parseKeyValPairs(parts[1])}
					}
				}
			}
		} else if mode == pmTasks {
			taskParts := strings.Split(line, "|")
			if len(taskParts) != 3 {
				panic("task should have format: cat | title | efforts") // TODO convert to error
			}
			keyValPairs := parseKeyValPairs(taskParts[2])
			delete(keyValPairs, "risks")
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
	return projParsed
}
