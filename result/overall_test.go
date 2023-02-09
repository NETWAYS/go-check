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

	assert.Equal(t, 1, overall.OKs)
	assert.ElementsMatch(t, overall.Outputs, []string{"[OK] test ok"})
}

func TestOverall_AddWarning(t *testing.T) {
	overall := Overall{}
	overall.Add(1, "test warning")

	assert.Equal(t, 1, overall.Warnings)
	assert.ElementsMatch(t, overall.Outputs, []string{"[WARNING] test warning"})
}

func TestOverall_AddCritical(t *testing.T) {
	overall := Overall{}
	overall.Add(2, "test critical")

	assert.Equal(t, 1, overall.Criticals)
	assert.ElementsMatch(t, overall.Outputs, []string{"[CRITICAL] test critical"})
}

func TestOverall_AddUnknown(t *testing.T) {
	overall := Overall{}
	overall.Add(3, "test unknown")

	assert.Equal(t, 1, overall.Unknowns)
	assert.ElementsMatch(t, overall.Outputs, []string{"[UNKNOWN] test unknown"})
}

func TestOverall_GetStatus_GetSummary(t *testing.T) {
	var overall Overall

	overall = Overall{}
	assert.Equal(t, 3, overall.GetStatus())
	assert.Equal(t, "No status information", overall.GetSummary())

	overall = Overall{OKs: 1}
	assert.Equal(t, 0, overall.GetStatus())
	assert.Equal(t, "states: ok=1", overall.GetSummary())

	overall = Overall{Criticals: 2, OKs: 1, Warnings: 2, Unknowns: 1}
	assert.Equal(t, 2, overall.GetStatus())
	assert.Equal(t, "states: critical=2 unknown=1 warning=2 ok=1", overall.GetSummary())

	overall = Overall{Unknowns: 2, OKs: 1, Warnings: 2}
	assert.Equal(t, 3, overall.GetStatus())
	assert.Equal(t, "states: unknown=2 warning=2 ok=1", overall.GetSummary())

	overall = Overall{OKs: 1, Warnings: 2}
	assert.Equal(t, 1, overall.GetStatus())
	assert.Equal(t, "states: warning=2 ok=1", overall.GetSummary())

	overall = Overall{Summary: "foobar"}
	assert.Equal(t, "foobar", overall.GetSummary())
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

	fmt.Println(overall)
	// Output: {1 0 1 0  [[OK] One element is good [CRITICAL] The other is critical] []}
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
	// states: ok=1 ok=1
	// [OK] bla
	// \_ [OK] Subcheck1 Test|pd_test=5s
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
\_ [OK] Subcheck1 Test|pd_test=5s pd_test2=1099511627776kB;@3.14:7036874417766;549755813887:1208925819614629174706176;;18446744073709551615
\_ [WARNING] Subcheck2 Test|kl;jr2if;l2rkjasdf=5m asdf=18446744073709551615B
`
	assert.Equal(t, expectedString, resString)
}

func TestOverall_withSubchecks3(t *testing.T) {
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

func TestOverall_withSubchecks4(t *testing.T) {
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
    \_ [OK] SubSubcheck|foo=3 bar=300%
`

	assert.Equal(t, res, overall.GetOutput())
}
