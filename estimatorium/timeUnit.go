package estimatorium

type TimeUnit uint8

const (
	TimeUnitUnknown TimeUnit = iota
	Hr
	Day
)

var timeUnit2Str = map[TimeUnit]string{
	TimeUnitUnknown: "TimeUnitUnknown", Hr: "hr", Day: "day",
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
