package main

import (
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
)

func main() {
	defer check.CatchPanic()

	var overall result.Overall

	check1 := result.NewPartialResult()

	check1.Output = "Check1"
	check1.SetState(check.OK)

	check1.Perfdata.Add(&perfdata.Perfdata{
		Label: "foo",
		Value: 23,
	})

	overall.AddSubcheck(check1)

	check2 := result.NewPartialResult()

	check2.Output = "Check2"
	check2.SetState(check.Warning)

	check2.Perfdata.Add(&perfdata.Perfdata{
		Label: "bar",
		Value: 42,
	})
	check2.Perfdata.Add(&perfdata.Perfdata{
		Label: "foo2 bar",
		Value: 46,
	})

	overall.AddSubcheck(check2)

	check.Exit(overall.GetStatus(), overall.GetOutput())
}
