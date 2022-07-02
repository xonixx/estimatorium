package estimatorium

type Project struct {
	TimeUnit TimeUnit
	Currency Currency
	Team     map[string]Resource
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
