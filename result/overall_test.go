package result

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
)

func TestOverall_AddOK(t *testing.T) {
	overall := Overall{}
	overall.Add(0, "test ok")

	if overall.oks != 1 {
		t.Fatalf("expected 1, got %d", overall.oks)
	}

	expectedOutputs := []string{"[OK] test ok"}
	if !reflect.DeepEqual(overall.Outputs, expectedOutputs) {
		t.Fatalf("expected %v, got %v", expectedOutputs, overall.Outputs)
	}
}

func TestOverall_AddWarning(t *testing.T) {
	overall := Overall{}
	overall.Add(1, "test warning")

	if overall.warnings != 1 {
		t.Fatalf("expected 1, got %d", overall.warnings)
	}

	expectedOutputs := []string{"[WARNING] test warning"}
	if !reflect.DeepEqual(overall.Outputs, expectedOutputs) {
		t.Fatalf("expected %v, got %v", expectedOutputs, overall.Outputs)
	}
}

func TestOverall_AddCritical(t *testing.T) {
	overall := Overall{}
	overall.Add(2, "test critical")

	if overall.criticals != 1 {
		t.Fatalf("expected 1, got %d", overall.criticals)
	}

	expectedOutputs := []string{"[CRITICAL] test critical"}
	if !reflect.DeepEqual(overall.Outputs, expectedOutputs) {
		t.Fatalf("expected %v, got %v", expectedOutputs, overall.Outputs)
	}
}

func TestOverall_AddUnknown(t *testing.T) {
	overall := Overall{}
	overall.Add(3, "test unknown")

	if overall.unknowns != 1 {
		t.Fatalf("expected 1, got %d", overall.unknowns)
	}

	expectedOutputs := []string{"[UNKNOWN] test unknown"}
	if !reflect.DeepEqual(overall.Outputs, expectedOutputs) {
		t.Fatalf("expected %v, got %v", expectedOutputs, overall.Outputs)
	}
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
			actual:          Overall{oks: 1, stateSetExplicitly: true},
			expectedSummary: "states: ok=1",
			expectedStatus:  0,
		},
		{
			actual:          Overall{criticals: 2, oks: 1, warnings: 2, unknowns: 1, stateSetExplicitly: true},
			expectedSummary: "states: critical=2 unknown=1 warning=2 ok=1",
			expectedStatus:  2,
		},
		{
			actual:          Overall{unknowns: 2, oks: 1, warnings: 2, stateSetExplicitly: true},
			expectedSummary: "states: unknown=2 warning=2 ok=1",
			expectedStatus:  3,
		},
		{
			actual:          Overall{oks: 1, warnings: 2, stateSetExplicitly: true},
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
		if test.expectedStatus != test.actual.GetStatus() {
			t.Fatalf("expected %d, got %d", test.expectedStatus, test.actual.GetStatus())
		}

		if test.expectedSummary != test.actual.GetSummary() {
			t.Fatalf("expected %s, got %s", test.expectedSummary, test.actual.GetSummary())
		}
	}
}

func TestOverall_GetOutput(t *testing.T) {
	var overall Overall

	overall = Overall{}
	overall.Add(0, "First OK")
	overall.Add(0, "Second OK")

	expected := "states: ok=2\n[OK] First OK\n[OK] Second OK\n"

	if expected != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", expected, overall.GetOutput())
	}

	overall = Overall{}
	overall.Add(0, "State OK")

	expected = "states: ok=1\n[OK] State OK\n"

	if expected != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", expected, overall.GetOutput())
	}

	// TODO: compress when only one state
	overall = Overall{}
	overall.Add(0, "First OK")
	overall.Add(2, "Second Critical")
	overall.Summary = "Custom Summary"

	expected = "Custom Summary\n[OK] First OK\n[CRITICAL] Second Critical\n"

	if expected != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", expected, overall.GetOutput())
	}
}

func ExampleOverall_Add() {
	overall := Overall{}
	overall.Add(check.OK, "One element is good")
	overall.Add(check.Critical, "The other is critical")

	fmt.Printf("%#v\n", overall)
	// Output: result.Overall{oks:1, warnings:0, criticals:1, unknowns:0, Summary:"", stateSetExplicitly:true, Outputs:[]string{"[OK] One element is good", "[CRITICAL] The other is critical"}, PartialResults:[]result.PartialResult(nil)}
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
		Output:   "Subcheck1 Test",
		Perfdata: pd_list,
	}

	subcheck.SetState(check.OK)

	overall.AddSubcheck(subcheck)
	overall.Add(0, "bla")

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
		Output:   "Subcheck1 Test",
		Perfdata: pd_list,
	}

	subcheck.SetState(check.OK)

	subcheck2 := PartialResult{
		Output:   "Subcheck2 Test",
		Perfdata: pd_list2,
	}

	subcheck2.SetState(check.Warning)

	overall.AddSubcheck(subcheck)
	overall.AddSubcheck(subcheck2)

	resString := overall.GetOutput()
	//nolint:lll
	expectedString := `states: warning=1 ok=1
\_ [OK] Subcheck1 Test
\_ [WARNING] Subcheck2 Test
|pd_test=5s pd_test2=1099511627776kB;@3.14:7036874417766;549755813887:1208925819614629174706176;;18446744073709551615 kl;jr2if;l2rkjasdf=5m asdf=18446744073709551615B
`

	if expectedString != resString {
		t.Fatalf("expected %s, got %s", expectedString, resString)
	}

	if check.Warning != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", check.Warning, overall.GetStatus())
	}
}

func TestOverall_withSubchecks_Simple_Output(t *testing.T) {
	var overall Overall

	subcheck2 := PartialResult{
		Output: "SubSubcheck",
	}

	subcheck2.SetState(check.OK)

	subcheck := PartialResult{
		Output: "PartialResult",
	}

	subcheck.SetState(check.OK)

	subcheck.PartialResults = append(subcheck.PartialResults, subcheck2)

	overall.AddSubcheck(subcheck)

	output := overall.GetOutput()

	resString := `states: ok=1
\_ [OK] PartialResult
    \_ [OK] SubSubcheck
`

	if output != resString {
		t.Fatalf("expected %s, got %s", output, resString)
	}
}

func TestOverall_withSubchecks_Perfdata(t *testing.T) {
	var overall Overall

	subcheck2 := PartialResult{
		Output: "SubSubcheck",
	}

	subcheck2.SetState(check.OK)

	subcheck := PartialResult{
		Output: "PartialResult",
	}

	subcheck.SetState(check.OK)

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
	subcheck.PartialResults = append(subcheck.PartialResults, subcheck2)

	overall.AddSubcheck(subcheck)

	res := `states: ok=1
\_ [OK] PartialResult
    \_ [OK] SubSubcheck
|foo=3 bar=300%
`

	if res != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", res, overall.GetOutput())
	}

	if 0 != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", 0, overall.GetStatus())
	}
}

func TestOverall_withSubchecks_PartialResult(t *testing.T) {
	var overall Overall

	subcheck3 := PartialResult{
		Output: "SubSubSubcheck",
	}

	subcheck3.SetState(check.Critical)

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
	perf3 := perfdata.Perfdata{
		Label: "baz",
		Value: 23,
		Uom:   "B",
	}

	subcheck3.Perfdata.Add(&perf3)
	subcheck2.Perfdata.Add(&perf1)
	subcheck2.Perfdata.Add(&perf2)
	subcheck2.PartialResults = append(subcheck.PartialResults, subcheck3)
	subcheck.PartialResults = append(subcheck.PartialResults, subcheck2)

	overall.AddSubcheck(subcheck)

	res := `states: critical=1
\_ [CRITICAL] PartialResult
    \_ [CRITICAL] SubSubcheck
        \_ [CRITICAL] SubSubSubcheck
|foo=3 bar=300% baz=23B
`

	if res != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", res, overall.GetOutput())
	}

	if check.Critical != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", 2, overall.GetStatus())
	}
}

func TestOverall_withSubchecks_PartialResultStatus(t *testing.T) {
	var overall Overall

	subcheck := PartialResult{
		Output: "Subcheck",
	}

	subcheck.SetState(check.OK)

	subsubcheck := PartialResult{
		Output: "SubSubcheck",
	}

	subsubcheck.SetState(check.Warning)

	subsubsubcheck := PartialResult{
		Output: "SubSubSubcheck",
	}

	subsubsubcheck.SetState(check.Critical)

	subsubcheck.AddSubcheck(subsubsubcheck)
	subcheck.AddSubcheck(subsubcheck)
	overall.AddSubcheck(subcheck)

	res := `states: ok=1
\_ [OK] Subcheck
    \_ [WARNING] SubSubcheck
        \_ [CRITICAL] SubSubSubcheck
`

	if res != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", res, overall.GetOutput())
	}

	if 0 != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", 0, overall.GetStatus())
	}
}

func TestSubchecksPerfdata(t *testing.T) {
	var overall Overall

	check1 := PartialResult{
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

	check1.SetState(check.OK)

	check2 := PartialResult{
		Output: "Check2",
		Perfdata: perfdata.PerfdataList{
			&perfdata.Perfdata{
				Label: "foo2 bar",
				Value: 46,
			},
		},
	}

	check2.SetState(check.Warning)

	overall.AddSubcheck(check1)
	overall.AddSubcheck(check2)

	resultString := "states: warning=1 ok=1\n\\_ [OK] Check1\n\\_ [WARNING] Check2\n|foo=23 bar=42 'foo2 bar'=46\n"

	if resultString != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", resultString, overall.GetOutput())
	}
}

func TestDefaultStates1(t *testing.T) {
	var overall Overall

	subcheck := PartialResult{}

	subcheck.SetDefaultState(check.OK)

	overall.AddSubcheck(subcheck)

	if check.OK != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", check.OK, overall.GetStatus())
	}
}

func TestDefaultStates2(t *testing.T) {
	var overall Overall

	subcheck := PartialResult{}

	overall.AddSubcheck(subcheck)

	if check.Unknown != subcheck.GetStatus() {
		t.Fatalf("expected %d, got %d", check.Unknown, subcheck.GetStatus())
	}

	if check.Unknown != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", check.Unknown, overall.GetStatus())
	}
}

func TestDefaultStates3(t *testing.T) {
	var overall Overall

	subcheck := PartialResult{}
	subcheck.SetDefaultState(check.OK)

	subcheck.SetState(check.Warning)

	overall.AddSubcheck(subcheck)

	if check.Warning != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", check.Warning, overall.GetStatus())
	}
}

func TestOverallOutputWithMultiLayerPartials(t *testing.T) {
	var overall Overall

	subcheck1 := PartialResult{}
	subcheck1.SetState(check.Warning)

	subcheck2 := PartialResult{}

	subcheck2_1 := PartialResult{}
	subcheck2_1.SetState(check.OK)

	subcheck2_2 := PartialResult{}
	subcheck2_2.SetState(check.Critical)

	subcheck2.AddSubcheck(subcheck2_1)
	subcheck2.AddSubcheck(subcheck2_2)

	overall.AddSubcheck(subcheck1)
	overall.AddSubcheck(subcheck2)

	resultString := "states: critical=1 warning=1\n\\_ [WARNING] \n\\_ [CRITICAL] \n    \\_ [OK] \n    \\_ [CRITICAL] \n"

	if resultString != overall.GetOutput() {
		t.Fatalf("expected %s, got %s", resultString, overall.GetOutput())
	}

	if check.Critical != overall.GetStatus() {
		t.Fatalf("expected %d, got %d", check.Critical, overall.GetStatus())
	}
}
