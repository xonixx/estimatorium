package estimatorium

// TODO validate Project model: 1. correct resources in tasks 2. correct risks etc.

type Project struct {
	TimeUnit TimeUnit
	Currency Currency
	Team     []Resource
	Risks    map[string]float32
	Tasks    []Task
}

type TimeUnit uint8

const (
	Hr TimeUnit = iota
	Day
)

var timeUnit2Str = map[TimeUnit]string{
	Hr: "hr", Day: "day",
}

func (tu TimeUnit) String() string {
	return timeUnit2Str[tu]
}

var timeUnit2Hrs = map[TimeUnit]int{
	Hr: 1, Day: 8,
}

func (tu TimeUnit) ToHours() int {
	return timeUnit2Hrs[tu]
}

type Currency uint8

const (
	Usd Currency = iota
	Eur
)

var currency2Str = map[Currency]string{
	Usd: "USD", Eur: "EUR",
}

func (c Currency) String() string {
	return currency2Str[c]
}

type Resource struct {
	Id      string
	Title   string
	Rate    float64
	Count   uint8
	Formula string
}

type Task struct {
	Category string
	Title    string
	Risk     string
	Work     map[string]float32 // resource -> time units
}

func StandardRisks() map[string]float32 {
	return map[string]float32{
		"low":    1.1,
		"medium": 1.5,
		"high":   2,
	}
}
