package estimatorium

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// TODO validate Project model: 1. correct resources in tasks 2. correct risks etc.

type Project struct {
	Name              string
	Author            string
	TimeUnit          TimeUnit
	Currency          Currency
	AcceptancePercent float64 // "Cleanup & acceptance" parameter
	Team              []Resource
	DesiredDuration   Duration // This will be treated as including risks
	Risks             map[string]float64
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
func (p Project) TeamAsMap() map[string]Resource {
	res := map[string]Resource{}
	for _, resource := range p.Team {
		res[resource.Id] = resource
	}
	return res
}
func (p Project) ResourceById(rId string) *Resource {
	for _, r := range p.Team {
		if r.Id == rId {
			return &r
		}
	}
	return nil
}

type Duration struct {
	duration float64
	unit     TimeUnit
}

func (d Duration) ToHours() float64 {
	return float64(d.unit.ToHours()) * d.duration
}

// ParseDuration should parse "10mth", "3 days", ".5weeks"
// should not parse "aaa", "zz day", "10bbb", ".5 ccc dd"
func ParseDuration(str string) (Duration, error) {
	str = strings.TrimSpace(str)
	for k, timeUnit := range timeUnitStr2Val {
		for _, attemptUnit := range []string{k, k + "s"} {
			if strings.HasSuffix(str, attemptUnit) {
				val := str[:len(str)-len(attemptUnit)]
				val = strings.TrimSpace(val)
				float, err := strconv.ParseFloat(val, 32)
				if err != nil {
					return Duration{}, err
				}
				return Duration{
					duration: float,
					unit:     timeUnit,
				}, nil
			}
		}
	}
	return Duration{}, errors.New("unknown time unit")
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
	Work     map[string]float64 // resource -> time units
}

func StandardRisks() map[string]float64 {
	return map[string]float64{
		"low":     1.1,
		"medium":  1.5,
		"high":    2,
		"extreme": 5,
	}
}

func RiskLabels(risks map[string]float64) []string {
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
