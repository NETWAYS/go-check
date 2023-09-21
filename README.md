go-check
========

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/NETWAYS/go-check?label=version)](https://github.com/NETWAYS/go-check/releases)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/NETWAYS/go-check)
[![Test Status](https://github.com/NETWAYS/go-check/workflows/Go/badge.svg)](https://github.com/NETWAYS/go-check/actions?query=workflow%3AGo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/NETWAYS/go-check)
![GitHub](https://img.shields.io/github/license/NETWAYS/go-check?color=green)

go-check is a library to help with development of monitoring plugins for tools like Icinga.

See the [documentation on pkg.go.dev](https://pkg.go.dev/github.com/NETWAYS/go-check) for more details and examples.

# Usage

## Simple Example

```go
package main

import (
	"github.com/NETWAYS/go-check"
)

func main() {
	config := check.NewConfig()
	config.Name = "check_test"
	config.Readme = `Test Plugin`
	config.Version = "1.0.0"

	_ = config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")

	config.ParseArguments()

	// Some checking should be done here, when --help is not passed
	check.Exitf(check.OK, "Everything is fine - answer=%d", 42)
    // Output:
    // OK - Everything is fine - answer=42
}
```

## Exit Codes

```
check.Exitf(OK, "Everything is fine - value=%d", 42) // OK, 0

check.ExitRaw(check.Critical, "CRITICAL", "|", "percent_packet_loss=100") // CRITICAL, 2

err := fmt.Errorf("connection to %s has been timed out", "localhost:12345")

check.ExitError(err) // UNKNOWN, 3
```

## Timeout Handling

```
checkPluginTimeoutInSeconds := 10
go check.HandleTimeout(checkPluginTimeoutInSeconds)
```

## Thresholds

Threshold objects represent monitoring plugin thresholds that have methods to evaluate if a given input is within the range.

They can be created with the ParseThreshold parser.

https://github.com/monitoring-plugins/monitoring-plugin-guidelines/blob/main/definitions/01.range_expressions.md

```
warnThreshold, err := check.ParseThreshold("~:3")

if err != nil {
    return t, err
}

if warnThreshold.DoesViolate(3.6) {
    fmt.Println("Not great, not terrible.")
}
```

## Perfdata

The Perfdata object represents monitoring plugin performance data that relates to the actual execution of a host or service check.

https://github.com/monitoring-plugins/monitoring-plugin-guidelines/blob/main/monitoring_plugins_interface/03.Output.md#performance-data

```
var pl perfdata.PerfdataList

pl.Add(&perfdata.Perfdata{
    Label: "process.cpu.percent",
    Value: 25,
    Uom:   "%",
    Warn:  50,
    Crit:  90,
    Min:   0,
    Max:   100})

fmt.Println(pl.String())
```

## Results

```
allStates = []int{0,2,3,0,1,2}

switch result.WorstState(allStates...) {
case 0:
    rc = check.OK
case 1:
    rc = check.Warning
case 2:
    rc = check.Critical
default:
    rc = check.Unknown
}
```

## Partial Results

```
o := Overall{}
o.Add(0, "Something is OK")

pr := PartialResult{
    State:  check.OK,
    Output: "My Subcheck",
}

o.AddSubcheck(pr)

fmt.Println(o.GetOutput())

// states: ok=1
// [OK] Something is OK
// \_ [OK] My Subcheck
```


# Examples

A few plugins using go-check:

* [check_cloud_aws](https://github.com/NETWAYS/check_cloud_aws)
* [check_logstash](https://github.com/NETWAYS/check_logstash)
* [check_sentinelone](https://github.com/NETWAYS/check_sentinelone)
* [check_sophos_central](https://github.com/NETWAYS/check_sophos_central)

# License

Copyright (c) 2020 [NETWAYS GmbH](mailto:info@netways.de)

This library is distributed under the GPL-2.0 or newer license found in the [COPYING](./COPYING)
file.
