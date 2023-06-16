package result

import (
	"fmt"
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/stretchr/testify/assert"
)

func TestOverall_AddOK(t *testing.T) {
	overall := Overall{}
	overall.Add(0, "test ok")

	assert.Equal(t, 1, overall.oks)
	assert.ElementsMatch(t, overall.Outputs, []string{"[OK] test ok"})
}

func TestOverall_AddWarning(t *testing.T) {
	overall := Overall{}
	overall.Add(1, "test warning")

	assert.Equal(t, 1, overall.warnings)
	assert.ElementsMatch(t, overall.Outputs, []string{"[WARNING] test warning"})
}

func TestOverall_AddCritical(t *testing.T) {
	overall := Overall{}
	overall.Add(2, "test critical")

	assert.Equal(t, 1, overall.criticals)
	assert.ElementsMatch(t, overall.Outputs, []string{"[CRITICAL] test critical"})
}

func TestOverall_AddUnknown(t *testing.T) {
	overall := Overall{}
	overall.Add(3, "test unknown")

	assert.Equal(t, 1, overall.unknowns)
	assert.ElementsMatch(t, overall.Outputs, []string{"[UNKNOWN] test unknown"})
}

func TestOverall_GetStatus_GetSummary(t *testing.T) {
	testcases := []struct {
		actual          Overall
		expectedSummary string
		expectedStatus  int
	}{
		{
			actual:          Overall{},
			expectedSummary: "No status information",
			expectedStatus:  3,
		},
		{
			actual:          Overall{oks: 1, stateSetExplicitely: true},
			expectedSummary: "states: ok=1",
			expectedStatus:  0,
		},
		{
			actual:          Overall{criticals: 2, oks: 1, warnings: 2, unknowns: 1, stateSetExplicitely: true},
			expectedSummary: "states: critical=2 unknown=1 warning=2 ok=1",
			expectedStatus:  2,
		},
		{
			actual:          Overall{unknowns: 2, oks: 1, warnings: 2, stateSetExplicitely: true},
			expectedSummary: "states: unknown=2 warning=2 ok=1",
			expectedStatus:  3,
		},
		{
			actual:          Overall{oks: 1, warnings: 2, stateSetExplicitely: true},
			expectedSummary: "states: warning=2 ok=1",
			expectedStatus:  1,
		},
		{
			actual:          Overall{Summary: "foobar"},
			expectedSummary: "foobar",
			expectedStatus:  3,
		},
	}

	for _, test := range testcases {
		assert.Equal(t, test.expectedStatus, test.actual.GetStatus())
		assert.Equal(t, test.expectedSummary, test.actual.GetSummary())
	}
}

func TestOverall_GetOutput(t *testing.T) {
	var overall Overall

	overall = Overall{}
	overall.Add(0, "First OK")
	overall.Add(0, "Second OK")

	assert.Equal(t, "states: ok=2\n[OK] First OK\n[OK] Second OK\n", overall.GetOutput())

	overall = Overall{}
	overall.Add(0, "State OK")
	// TODO: compress when only one state
	assert.Equal(t, "states: ok=1\n[OK] State OK\n", overall.GetOutput())

	overall = Overall{}
	overall.Add(0, "First OK")
	overall.Add(2, "Second Critical")
	overall.Summary = "Custom Summary"
	assert.Equal(t, "Custom Summary\n[OK] First OK\n[CRITICAL] Second Critical\n", overall.GetOutput())
}

func ExampleOverall_Add() {
	overall := Overall{}
	overall.Add(check.OK, "One element is good")
	overall.Add(check.Critical, "The other is critical")

	fmt.Printf("%#v\n", overall)
	// Output: result.Overall{oks:1, warnings:0, criticals:1, unknowns:0, Summary:"", stateSetExplicitely:true, Outputs:[]string{"[OK] One element is good", "[CRITICAL] The other is critical"}, partialResults:[]result.PartialResult(nil)}
}

func ExampleOverall_GetOutput() {
	overall := Overall{}
	overall.Add(check.OK, "One element is good")
	overall.Add(check.Critical, "The other is critical")

	fmt.Println(overall.GetOutput())
	// Output:
	// states: critical=1 ok=1
	// [OK] One element is good
	// [CRITICAL] The other is critical
}

func ExampleOverall_GetStatus() {
	overall := Overall{}
	overall.Add(check.OK, "One element is good")
	overall.Add(check.Critical, "The other is critical")

	fmt.Println(overall.GetStatus())
	// Output: 2
}

func ExampleOverall_withSubchecks() {
	var overall Overall

	example_perfdata := perfdata.Perfdata{Label: "pd_test", Value: 5, Uom: "s"}
	pd_list := perfdata.PerfdataList{}
	pd_list.Add(&example_perfdata)

	subcheck := PartialResult{
		State:    check.OK,
		Output:   "Subcheck1 Test",
		Perfdata: pd_list,
	}

	overall.AddSubcheck(subcheck)
	overall.AddOK("bla")

	fmt.Println(overall.GetOutput())
	// Output:
	// states: ok=1
	// [OK] bla
	// \_ [OK] Subcheck1 Test
	// |pd_test=5s
}

func TestOverall_withEnhancedSubchecks(t *testing.T) {
	var overall Overall

	example_perfdata := perfdata.Perfdata{Label: "pd_test", Value: 5, Uom: "s"}
	example_perfdata2 := perfdata.Perfdata{
		Label: "pd_test2",
		Value: 1099511627776,
		Uom:   "kB",
		Warn:  &check.Threshold{Inside: true, Lower: 3.14, Upper: 0x66666666666},
		Crit:  &check.Threshold{Inside: false, Lower: 07777777777777, Upper: 0xFFFFFFFFFFFFFFFFFFFF},
		Max:   uint64(18446744073709551615),
	}
	example_perfdata3 := perfdata.Perfdata{Label: "kl;jr2if;l2rkjasdf", Value: 5, Uom: "m"}
	example_perfdata4 := perfdata.Perfdata{Label: "asdf", Value: uint64(18446744073709551615), Uom: "B"}

	pd_list := perfdata.PerfdataList{}
	pd_list.Add(&example_perfdata)
	pd_list.Add(&example_perfdata2)

	pd_list2 := perfdata.PerfdataList{}
	pd_list2.Add(&example_perfdata3)
	pd_list2.Add(&example_perfdata4)

	subcheck := PartialResult{
		State:    check.OK,
		Output:   "Subcheck1 Test",
		Perfdata: pd_list,
	}
	subcheck2 := PartialResult{
		State:    check.Warning,
		Output:   "Subcheck2 Test",
		Perfdata: pd_list2,
	}

	overall.AddSubcheck(subcheck)
	overall.AddSubcheck(subcheck2)

	resString := overall.GetOutput()
	//nolint:lll
	expectedString := `states: warning=1 ok=1
\_ [OK] Subcheck1 Test
\_ [WARNING] Subcheck2 Test
|pd_test=5s pd_test2=1099511627776kB;@3.14:7036874417766;549755813887:1208925819614629174706176;;18446744073709551615 kl;jr2if;l2rkjasdf=5m asdf=18446744073709551615B
`
	assert.Equal(t, expectedString, resString)
}

func TestOverall_withSubchecks_Simple_Output(t *testing.T) {
	var overall Overall

	subcheck2 := PartialResult{
		State:  check.OK,
		Output: "SubSubcheck",
	}
	subcheck := PartialResult{
		State:  check.OK,
		Output: "PartialResult",
	}
	subcheck.partialResults = append(subcheck.partialResults, subcheck2)

	overall.AddSubcheck(subcheck)

	output := overall.GetOutput()

	resString := `states: ok=1
\_ [OK] PartialResult
    \_ [OK] SubSubcheck
`

	assert.Equal(t, resString, output)
}

func TestOverall_withSubchecks_Perfdata(t *testing.T) {
	var overall Overall

	subcheck2 := PartialResult{
		State:  check.OK,
		Output: "SubSubcheck",
	}
	subcheck := PartialResult{
		State:  check.OK,
		Output: "PartialResult",
	}

	perf1 := perfdata.Perfdata{
		Label: "foo",
		Value: 3,
	}
	perf2 := perfdata.Perfdata{
		Label: "bar",
		Value: 300,
		Uom:   "%",
	}

	subcheck2.Perfdata.Add(&perf1)
	subcheck2.Perfdata.Add(&perf2)
	subcheck.partialResults = append(subcheck.partialResults, subcheck2)

	overall.AddSubcheck(subcheck)

	res := `states: ok=1
\_ [OK] PartialResult
    \_ [OK] SubSubcheck
|foo=3 bar=300%
`

	assert.Equal(t, res, overall.GetOutput())
	assert.Equal(t, 0, overall.GetStatus())
}

func TestOverall_withSubchecks_PartialResult(t *testing.T) {
	var overall Overall

	subcheck3 := PartialResult{
		State:  check.Critical,
		Output: "SubSubSubcheck",
	}
	subcheck2 := PartialResult{
		Output: "SubSubcheck",
	}
	subcheck := PartialResult{
		Output: "PartialResult",
	}

	perf1 := perfdata.Perfdata{
		Label: "foo",
		Value: 3,
	}
	perf2 := perfdata.Perfdata{
		Label: "bar",
		Value: 300,
		Uom:   "%",
	}

	subcheck2.Perfdata.Add(&perf1)
	subcheck2.Perfdata.Add(&perf2)
	subcheck2.partialResults = append(subcheck.partialResults, subcheck3)
	subcheck.partialResults = append(subcheck.partialResults, subcheck2)

	overall.AddSubcheck(subcheck)

	res := `states: ok=1
\_ [OK] PartialResult
    \_ [OK] SubSubcheck
        \_ [CRITICAL] SubSubSubcheck
|foo=3 bar=300%
`

	assert.Equal(t, res, overall.GetOutput())
	assert.Equal(t, 0, overall.GetStatus())
}

func TestOverall_withSubchecks_PartialResultStatus(t *testing.T) {
	var overall Overall

	subcheck := PartialResult{
		State:  check.OK,
		Output: "Subcheck",
	}
	subsubcheck := PartialResult{
		State:  check.Warning,
		Output: "SubSubcheck",
	}
	subsubsubcheck := PartialResult{
		State:  check.Critical,
		Output: "SubSubSubcheck",
	}

	subsubcheck.AddSubcheck(subsubsubcheck)
	subcheck.AddSubcheck(subsubcheck)
	overall.AddSubcheck(subcheck)

	res := `states: ok=1
\_ [OK] Subcheck
    \_ [WARNING] SubSubcheck
        \_ [CRITICAL] SubSubSubcheck
`
	assert.Equal(t, res, overall.GetOutput())
	assert.Equal(t, 0, overall.GetStatus())
}

func TestSubchecksPerfdata(t *testing.T) {
	var overall Overall

	check1 := PartialResult{
		State:  check.OK,
		Output: "Check1",
		Perfdata: perfdata.PerfdataList{
			&perfdata.Perfdata{
				Label: "foo",
				Value: 23,
			},
			&perfdata.Perfdata{
				Label: "bar",
				Value: 42,
			},
		},
	}
	check2 := PartialResult{
		State:  check.Warning,
		Output: "Check2",
		Perfdata: perfdata.PerfdataList{
			&perfdata.Perfdata{
				Label: "foo2",
				Value: 46,
			},
		},
	}

	overall.AddSubcheck(check1)
	overall.AddSubcheck(check2)

	resultString := "states: warning=1 ok=1\n\\_ [OK] Check1\n\\_ [WARNING] Check2\n|foo=23 bar=42 foo2=46\n"

	assert.Equal(t, resultString, overall.GetOutput())
}
