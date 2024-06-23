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

	pd := &perfdata.Perfdata{
		Label: "foo",
		Value: perfdata.NewPdvUint64(23),
	}

	check1.Perfdata.Add(pd)

	check2 := result.PartialResult{}

	check2.Output = "Check2"
	err = check2.SetState(check.Warning)

	if err != nil {
		check.ExitError(err)
	}

	pd2 := &perfdata.Perfdata{
		Label: "bar",
		Value: perfdata.NewPdvUint64(42),
	}
	check2.Perfdata.Add(pd2)

	pd3 := &perfdata.Perfdata{
		Label: "foo2 bar",
		Value: perfdata.NewPdvUint64(46),
	}
	check2.Perfdata.Add(pd3)

	overall.AddSubcheck(check1)
	overall.AddSubcheck(check2)

	check.ExitRaw(overall.GetStatus(), overall.GetOutput())
}
