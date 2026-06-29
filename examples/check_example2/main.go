package main

import (
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

func main() {
	defer check.CatchPanic()

	var overall result.Overall

	check1 := result.NewPartialResult()

	check1.SetOutput("Check1")
	check1.SetState(check.OK)

	check1.AddPerfdata(&check.Perfdata{
		Label: "foo",
		Value: 23,
	})

	overall.AddSubcheck(check1)

	check2 := result.NewPartialResult()

	check2.SetOutput("Check2")
	check2.SetState(check.Warning)

	check2.AddPerfdata(&check.Perfdata{
		Label: "bar",
		Value: 42,
	})
	check2.AddPerfdata(&check.Perfdata{
		Label: "foo2 bar",
		Value: 46,
	})

	overall.AddSubcheck(check2)

	overall.Add(check.Warning, "Check3")
	overall.Add(check.OK, "Check4")

	check.Exit(overall.GetStatus(), overall.GetOutput())
}
