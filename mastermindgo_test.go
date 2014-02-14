package main

import "testing"

func TestCalcPips(t *testing.T) {
	type test struct {
		code  string
		guess string
		black int
		white int
	}

	testCases := []test{
		test{"bbbb", "gggg", 0, 0},
		test{"bbbg", "gggy", 0, 1},
		test{"ggbb", "yygg", 0, 2},
		test{"bryg", "rbgs", 0, 3},
		test{"bryg", "gbry", 0, 4},
		test{"bbbb", "bggg", 1, 0},
		test{"bbby", "byrr", 1, 1},
		test{"bbry", "brys", 1, 2},
		test{"brsy", "rsby", 1, 3},
		test{"bbbb", "bbgg", 2, 0},
		test{"bbbg", "bbgy", 2, 1},
		test{"bbry", "bbyr", 2, 2},
		test{"bbbs", "bbby", 3, 0},
		test{"brys", "brys", 4, 0}}

	for _, tc := range testCases {
		b, w := calcPips(chars(tc.code), chars(tc.guess))
		if b != tc.black || w != tc.white {
			t.Errorf("Failed calcPips(\"%s\", \"%s\"). Expected (%d, %d), got (%d, %d).",
				tc.code, tc.guess,
				tc.black, tc.white,
				b, w)
		}
	}
}

// run benchmarks with `go test -bench .`
func BenchmarkCalcPips(b *testing.B) {
	tests := []struct{ code, guess chars }{
		{chars("bbbbaa"), chars("ggggaa")},
		{chars("bbbgaa"), chars("gggyaa")},
		{chars("ggbbaa"), chars("yyggaa")},
		{chars("brygaa"), chars("rbgsaa")},
		{chars("brygaa"), chars("gbryaa")},
		{chars("bbbbaa"), chars("bgggaa")},
		{chars("bbbyaa"), chars("byrraa")},
		{chars("bbryaa"), chars("brysaa")},
		{chars("brsyaa"), chars("rsbyaa")},
		{chars("bbbbaa"), chars("bbggaa")},
		{chars("bbbgaa"), chars("bbgyaa")},
		{chars("bbryaa"), chars("bbyraa")},
		{chars("bbbsaa"), chars("bbbyaa")},
		{chars("brysaa"), chars("brysaa")}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range tests {
			calcPips(c.code, c.guess)
		}
	}
}

func TestIndexOf(t *testing.T) {
	type test struct {
		str      string
		val      string
		position int
	}

	testCases := []test{
		test{"12345", "4", 3},
		test{"12435", "1", 0},
		test{"12345", "5", 4},
		test{"12345", "z", -1}}

	for _, tc := range testCases {
		if pos := indexOf(chars(tc.str), []byte(tc.val)[0]); pos != tc.position {
			t.Errorf("Failed indexOf(\"%s\", \"%s\"). Expected %d, got %d.",
				tc.str, tc.val, tc.position, pos)
		}
	}
}
