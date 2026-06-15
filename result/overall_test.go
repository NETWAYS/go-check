package result

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
)

func TestOverall_NewPartialResult(t *testing.T) {
	actual := NewPartialResult()

	if actual.String() != "[UNKNOWN] " {
		t.Fatalf("expected '[UNKNOWN] ', got %s", actual.String())
	}
}

func TestOverall_AddOK(t *testing.T) {
	overall := Overall{}
	overall.Add(0, "test ok")

	counts := overall.getStatusCount()

	if counts.OK != 1 && counts.Critical != 0 && counts.Warning != 0 && counts.Unknown != 0 {
		t.Fatalf("expected 1, got %d", counts.OK)
	}

	expectedOutput := "states: ok=1\n\\_ [OK] test ok\n"
	if !reflect.DeepEqual(overall.GetOutput(), expectedOutput) {
		t.Fatalf("expected \n%q\n, got \n%q\n", expectedOutput, overall.GetOutput())
	}
}

func TestOverall_AddWarning(t *testing.T) {
	overall := Overall{}
	overall.Add(1, "test warning")

	counts := overall.getStatusCount()

	if counts.OK != 0 && counts.Critical != 0 && counts.Warning != 1 && counts.Unknown != 0 {
		t.Fatalf("expected 1, got %d", counts.Warning)
	}

	expectedOutput := "test warning\n\\_ [WARNING] test warning\n"
	if !reflect.DeepEqual(overall.GetOutput(), expectedOutput) {
		t.Fatalf("expected %q\n, got %q", expectedOutput, overall.GetOutput())
	}
}

func TestOverall_AddCritical(t *testing.T) {
	overall := Overall{}
	overall.Add(2, "test critical")

	counts := overall.getStatusCount()

	if counts.OK != 0 && counts.Critical != 1 && counts.Warning != 0 && counts.Unknown != 0 {
		t.Fatalf("expected 1, got %d", counts.Critical)
	}

	expectedOutputs := "test critical\n\\_ [CRITICAL] test critical\n"
	if !reflect.DeepEqual(overall.GetOutput(), expectedOutputs) {
		t.Fatalf("expected %q, got %q", expectedOutputs, overall.GetOutput())
	}
}

func TestOverall_AddUnknown(t *testing.T) {
	overall := Overall{}
	overall.Add(3, "test unknown")

	counts := overall.getStatusCount()

	if counts.OK != 0 && counts.Critical != 0 && counts.Warning != 0 && counts.Unknown != 1 {
		t.Fatalf("expected 1, got %d", counts.Unknown)
	}

	expectedOutputs := "test unknown\n\\_ [UNKNOWN] test unknown\n"
	if !reflect.DeepEqual(overall.GetOutput(), expectedOutputs) {
		t.Fatalf("expected %q, got %q", expectedOutputs, overall.GetOutput())
	}
}

func TestOverall_GetOutput(t *testing.T) {
	var overall Overall

	overall = Overall{}
	overall.Add(0, "First OK")
	overall.Add(0, "Second OK")

	expected := "states: ok=2\n\\_ [OK] First OK\n\\_ [OK] Second OK\n"

	if expected != overall.GetOutput() {
		t.Fatalf("expected %q, got %q", expected, overall.GetOutput())
	}

	overall = Overall{}
	overall.Add(0, "State OK")

	expected = "states: ok=1\n\\_ [OK] State OK\n"

	if expected != overall.GetOutput() {
		t.Fatalf("expected %q, got %q", expected, overall.GetOutput())
	}

	// TODO: compress when only one state
	overall = Overall{}
	overall.Add(0, "First OK")
	overall.Add(2, "Second Critical")
	overall.OKSummary = "Custom Summary"

	expected = "Second Critical\n\\_ [OK] First OK\n\\_ [CRITICAL] Second Critical\n"

	if expected != overall.GetOutput() {
		t.Fatalf("expected %q, got %q", expected, overall.GetOutput())
	}
}

func ExampleOverall_Add() {
	overall := Overall{}
	overall.Add(check.OK, "One element is good")
	overall.Add(check.Critical, "The other is critical")

	fmt.Printf("%s", overall.GetOutput())
}

func ExampleOverall_GetOutput() {
	overall := Overall{}
	overall.Add(check.OK, "One element is good")
	overall.Add(check.Critical, "The other is critical")

	fmt.Println(overall.GetOutput())
	// Output:
	// The other is critical
	// \_ [OK] One element is good
	// \_ [CRITICAL] The other is critical
}

func ExampleOverall_GetStatus() {
	overall := Overall{}
	overall.Add(check.OK, "One element is good")
	overall.Add(check.Critical, "The other is critical")

	fmt.Println(overall.GetStatus())
	// Output: CRITICAL
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
	overall.Add(check.OK, "bla")

	fmt.Println(overall.GetOutput())
	// Output:
	// states: ok=2
	// \_ [OK] Subcheck1 Test
	// \_ [OK] bla
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

	expectedString := `Subcheck2 Test
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

	res := `SubSubSubcheck
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

	resultString := "Check2\n\\_ [OK] Check1\n\\_ [WARNING] Check2\n|foo=23 bar=42 'foo2 bar'=46\n"

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
