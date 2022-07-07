package estimatorium

import "sort"

// TODO validate Project model: 1. correct resources in tasks 2. correct risks etc.

type Project struct {
	Name              string
	TimeUnit          TimeUnit
	Currency          Currency
	AcceptancePercent float32 // "Cleanup & acceptance" parameter
	Team              []Resource
	Risks             map[string]float32
	Tasks             []Task
}

func (p Project) TeamExcludingDerived() []Resource {
	res := []Resource{}
	for _, r := range p.Team {
		if r.Formula == "" {
			res = append(res, r)
		}
	}
	return res
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

var timeUnitStr2Val = map[string]TimeUnit{
	"hr": Hr, "day": Day,
}

func TimeUnitFromString(tu string) TimeUnit {
	return timeUnitStr2Val[tu]
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

var currency2Symbol = map[Currency]string{
	Usd: "$", Eur: "€",
}

func (c Currency) Symbol() string {
	return currency2Symbol[c]
}

var currencyStr2Val = map[string]Currency{
	"USD": Usd, "EUR": Eur,
}

func CurrencyFromString(curr string) Currency {
	return currencyStr2Val[curr]
}

type Resource struct {
	Id      string
	Title   string
	Rate    float64
	Count   int
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
		"low":     1.1,
		"medium":  1.5,
		"high":    2,
		"extreme": 5,
	}
}

func RiskLabels(risks map[string]float32) []string {
	var keys []string
	for k := range risks {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		ki := keys[i]
		kj := keys[j]
		return risks[ki] < risks[kj]
	})
	return keys
}

var standardResourceTypes = map[string]string{
	"fe":    "Front dev",
	"be":    "Back dev",
	"mob":   "Mob dev",
	"ios":   "iOS dev",
	"droid": "Android dev",
	"do":    "DevOps",
	"pm":    "Project Manager",
	"ba":    "Business Analyst",
	"qa":    "QA Engineer",
	"ds":    "UI Designer",
}
