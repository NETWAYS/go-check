package main

import (
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"github.com/NETWAYS/go-check/result"
)

func main() {
	defer check.CatchPanic()

	var overall result.Overall

	check1 := result.PartialResult{}

	check1.Output = "Check1"
	err := check1.SetState(check.OK)

	if err != nil {
		check.ExitError(err)
	}

	check1.Perfdata.Add(&perfdata.Perfdata{
		Label: "foo",
		Value: 23,
	})

	check2 := result.PartialResult{}

	check2.Output = "Check2"
	err = check2.SetState(check.Warning)

	if err != nil {
		check.ExitError(err)
	}

	check2.Perfdata.Add(&perfdata.Perfdata{
		Label: "bar",
		Value: 42,
	})
	check2.Perfdata.Add(&perfdata.Perfdata{
		Label: "foo2 bar",
		Value: 46,
	})

	overall.AddSubcheck(check1)
	overall.AddSubcheck(check2)

	check.ExitRaw(overall.GetStatus(), overall.GetOutput())
}
