package estimatorium

import (
	"regexp"
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

var spaceRe = regexp.MustCompile("[ \t]")

func ProjectFromString(projData string) Project {
	return Project{} // TODO
}

func parseProj(projData string) projParsed {
	projParsed := projParsed{}
	var parsedDirectives map[string]directiveVals
	lines := strings.Split(projData, "\n")
	for _, line := range lines {
		parts := spaceRe.Split(line, 2)
		if parts[0] == "tasks" {
			break
		}
		for _, directive := range directives {
			if parts[0] == directive.name {
				if directive.directiveType == DtSingleValue {
					parsedDirectives[directive.name] = directiveVals{value: strings.TrimSpace(parts[1])}
				} else if directive.directiveType == DtKeyVal {
					valParts := spaceRe.Split(parts[1], -1)
					values := map[string]string{}
					for _, valPart := range valParts {
						keyVal := strings.SplitN(valPart, "=", 2)
						values[keyVal[0]] = values[keyVal[1]]
					}
					parsedDirectives[directive.name] = directiveVals{values: values}
				}
			}
		}
	}
	projParsed.directives = parsedDirectives
	return projParsed
}
