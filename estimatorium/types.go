package estimatorium

import (
	"sort"
)

// TODO validate Project model: 1. correct resources in tasks 2. correct risks etc.

type Project struct {
	Name              string
	TimeUnit          TimeUnit
	Currency          Currency
	AcceptancePercent float32 // "Cleanup & acceptance" parameter
	Team              []Resource
	DesiredDuration   *DesiredDuration
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

func (p Project) ResourceById(rId string) *Resource {
	for _, r := range p.Team {
		if r.Id == rId {
			return &r
		}
	}
	return nil
}

type DesiredDuration struct {
	duration float32
	unit     DesiredDurationTimeUnit
}

type DesiredDurationTimeUnit int

const (
	DDTUUnknown DesiredDurationTimeUnit = iota
	DDTUMonths
	DDTUWeeks
)

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
