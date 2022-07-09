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

const (
	WorkingHoursADay   = 8
	WorkingDaysInWeek  = 5
	WorkingDaysInMonth = 21
)

var timeUnit2Hrs = map[TimeUnit]int{
	Hr:    1,
	Day:   WorkingHoursADay,
	Week:  WorkingDaysInWeek * WorkingHoursADay,
	Month: WorkingDaysInMonth * WorkingHoursADay,
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
