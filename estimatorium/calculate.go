package estimatorium

import "math"

type ProjectCalculationResult struct {
	Team []Resource // team calculated based on desired duration
}

func (p *Project) Calculate() ProjectCalculationResult {
	res := ProjectCalculationResult{}

	if p.DesiredDuration == (Duration{}) {
		return res
	}

	desiredDuration := p.DesiredDuration
	desiredDurationHrs := desiredDuration.ToHours()

	work := map[string]float64{}
	for _, task := range p.Tasks {
		for resId, effort := range task.Work {
			work[resId] += effort * float64(p.TimeUnit.ToHours()) * p.Risks[task.Risk]
		}
	}

	calculatedTeam := []Resource{}
	for _, resource := range p.Team {
		workOfRes := work[resource.Id]
		if workOfRes > 0 {
			cntF := workOfRes / desiredDurationHrs
			cnt := int(math.Ceil(cntF))
			resource1 := resource
			resource1.Count = cnt
			calculatedTeam = append(calculatedTeam, resource1)
		} else if resource.Formula != "" {
			resource1 := resource
			resource1.Count = 1 // TODO
			calculatedTeam = append(calculatedTeam, resource1)
		}
	}
	p.Team = calculatedTeam

	return res
}
