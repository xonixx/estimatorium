package estimatorium

type TimeUnit uint8

const (
	TimeUnitUnknown TimeUnit = iota
	Hr
	Day
	Week
	Month
)

var timeUnit2Str = map[TimeUnit]string{
	TimeUnitUnknown: "TimeUnitUnknown", Hr: "hr", Day: "day", Week: "week", Month: "mth",
}

func (tu TimeUnit) String() string {
	return timeUnit2Str[tu]
}

var timeUnit2Hrs = map[TimeUnit]int{
	Hr: 1, Day: 8,
	Week:  5 * 8,  // 5 working days in week
	Month: 21 * 8, // 21 working days in mth
}

func (tu TimeUnit) ToHours() int {
	return timeUnit2Hrs[tu]
}

var timeUnitStr2Val = map[string]TimeUnit{
	"hr": Hr, "day": Day, "week": Week, "mth": Month,
}

func TimeUnitFromString(tu string) TimeUnit {
	return timeUnitStr2Val[tu]
}
