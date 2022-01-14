package result

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestOverall_withSubchecks(t *testing.T) {
	var overall Overall

	example_perfdata := perfdata.Perfdata{Label: "pd_test", Value: 5, Uom: "s"}
	pd_list := perfdata.PerfdataList{}
	pd_list.Add(&example_perfdata)
	subcheck := Subcheck{
		State: check.OK,
		Output: "Subcheck1 Test",
		Perfdata: pd_list,
	}

	overall.AddSubcheck(subcheck)
	overall.AddOK("bla")

	fmt.Println(overall.GetOutput())
	
	// Output: [OK] bla
    // |- [OK] Subcheck1 Test|pd_test=5s;;;;
}
