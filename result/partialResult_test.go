package result

import (
	"sync"
	"testing"

	"github.com/NETWAYS/go-check"
)

func TestOverall_NewPartialResult(t *testing.T) {
	actual := NewPartialResult()

	if actual.String() != "[UNKNOWN] " {
		t.Fatalf("expected '[UNKNOWN] ', got %s", actual.String())
	}
}

func TestPartialResult_SetGet_WithRace(t *testing.T) {
	pr := NewPartialResult()

	var wg sync.WaitGroup

	for range 3 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			pr.SetState(check.Critical)
		}()
		go func() {
			defer wg.Done()
			_ = pr.GetStatus()
		}()
	}
	wg.Wait()
}

func TestPartialResult_AddSubcheck_WithRace(t *testing.T) {
	parent := NewPartialResult()
	parent.SetOutput("unittest")

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			child := NewPartialResult()
			child.SetState(check.OK)
			parent.AddSubcheck(child)
		}()
	}
	wg.Wait()
}

func TestPartialResult_SetDefaultState_WithRace(t *testing.T) {
	pr := NewPartialResult()

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			example_perfdata := check.Perfdata{Label: "pd_test", Value: 5, Uom: "s"}
			pr.AddPerfdata(&example_perfdata)

			pr.SetDefaultState(check.Warning)
		}()
	}
	wg.Wait()
}
