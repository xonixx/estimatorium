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

type lineParser interface {
	myDirective() string
	parse(line string) (directiveParsed, error)
}

func ProjectFromString(projData string) Project {

}
