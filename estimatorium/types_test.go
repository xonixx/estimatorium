package estimatorium

import (
	"fmt"
	"testing"
)

func TestParseDurationCorrect(t *testing.T) {
	checkParseCorrect(t, "10mth", Duration{10, Month})
	checkParseCorrect(t, "3 days", Duration{3, Day})
	checkParseCorrect(t, ".5weeks", Duration{.5, Week})
}
func TestParseDurationIncorrect(t *testing.T) {
	checkParseIncorrect(t, "aaa")
	checkParseIncorrect(t, "zz day")
	checkParseIncorrect(t, "10bbb")
	checkParseIncorrect(t, ".5 ccc dd")
}

func checkParseCorrect(t *testing.T, v string, expected Duration) {
	duration, err := ParseDuration(v)
	if err != nil {
		t.Fatal(err)
	}
	if duration != expected {
		t.Fatalf("wrong value")
	}
}
func checkParseIncorrect(t *testing.T, v string) {
	_, err := ParseDuration(v)
	fmt.Println(err)
	if err == nil {
		t.Fatal("should error")
	}
}
