package estimatorium

// 1. each directive can go at most one time

type directiveParsed struct {
	directive string
	value     string
	values    map[string]string
}

type projParsed struct {
	directives   map[string]directiveParsed // each directive can go at most one time
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
	//parse(line string) (directiveParsed, error)
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

func ProjectFromString(projData string) Project {

}
