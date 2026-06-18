go-check
========

go-check is a Golang library to help with development of monitoring plugins for tools like Icinga.

See the [documentation on pkg.go.dev](https://pkg.go.dev/github.com/NETWAYS/go-check) for more details and examples.

# Usage

## Simple Example

go-check includes everything to quickly create a CLI monitoring plugin:

```go
package main

import (
    "fmt"

    "github.com/NETWAYS/go-check"
)

func main() {
    // Global configuration of the plugin
    config := check.NewConfig()
    config.Name = "check_test"
    config.Readme = `Test Plugin`
    config.Version = "1.0.0"

    // Command line arguments
    _ = config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")

    config.ParseArguments()

    // Handle exit with the desired exit code
    check.Exit(check.OK, fmt.Sprintf("Everything is fine - answer=%d", 42))
    // Output:
    // [OK] - Everything is fine - answer=42
}
```

## Return Codes

The library provides predefined return or exit codes:

```go
check.OK
check.Warning
check.Critical
check.Unknown

// These exit codes implement the Stringer interface
fmt.Println(check.OK)
```

To convert an integer or string into an exit code

```
unknown, err := NewStatus(3)

warning, err := NewStatusFromString("Warning")
```

See also: https://www.monitoring-plugins.org/doc/guidelines.html#AEN74

## Exit

The `Exit` function can be used to cause an exit with the given status code.

```go
check.Exit(check.OK, fmt.Sprintf("Everything is fine - value=%d", 42)) // OK, 0

// With perfdata
check.Exit(check.Critical, "CRITICAL", "|", "percent_packet_loss=100") // CRITICAL, 2
```

`ExitError` can be used to cause an exit with the given error.

```go
err := fmt.Errorf("connection to %s has been timed out", "localhost:12345")

check.ExitError(err)
// UNKNOWN, 3
```

## Timeout Handling

HandleTimeout is a helper for a goroutine, to wait for signals and timeout, and exit with a proper code.

```go
checkPluginTimeoutInSeconds := 10

go check.HandleTimeout(checkPluginTimeoutInSeconds)
```

## Thresholds

Threshold objects represent monitoring plugin thresholds that have methods to evaluate if a given input is within the range.

They can be created with the ParseThreshold parser.

```go
warnThreshold, err := check.ParseThreshold("~:3")

if err != nil {
    // Handle the error
}

if warnThreshold.DoesViolate(3.6) {
    fmt.Println("Not great, not terrible.")
}
```

See also: https://www.monitoring-plugins.org/doc/guidelines.html#THRESHOLDFORMAT

## Performance data

The `Perfdata` object represents monitoring plugin performance data that relates to the actual execution of a host or service check.

```go
var pl perfdata.PerfdataList

pl.Add(&perfdata.Perfdata{
    Label: "process.cpu.percent",
    Value: 25,
    Uom:   "%",
    Warn:  &check.Threshold{Lower: check.NegInf, Upper: 50 },
    Crit:  &check.Threshold{Lower: check.NegInf, Upper: 90},
    Min:   0,
    Max:   100})

fmt.Println(pl.String())
```

See also: https://www.monitoring-plugins.org/doc/guidelines.html#AEN197

## WorstState

The `WorstState` helper can be used to determine the worst exit status from a set of exit states.

```go
allStates = []check.Status{check.OK, check.Critical, check.Warning, check.Unknown}

rc := result.WorstState(allStates...)
```

## Overall and Partial Results

The `Overall` and `PartialResult` objects can be used to represent a simple parent-child relationship.

An `Overall` can contain multiple subchecks. The final exit of the `Overall` will be automatically determined by the worst state of a `PartialResult`.

```go
o := Overall{}
o.Add(0, "Something is OK")

pr := PartialResult{
    Output: "My Subcheck",
}

if err := pr.SetState(check.OK); err != nil {
  fmt.Printf(%s, err)
}

o.AddSubcheck(pr)

fmt.Println(o.GetOutput())

// states: ok=1
// [OK] Something is OK
// \_ [OK] My Subcheck
```

Overall is concurrency-safe.

## Human-readable bytes

`ParseBytes` is a helper that can be used to parse string containering IEC or SI bytes into the number of bytes.

```go
b, err := ParseBytes("2MiB")
// uint64 2 * 1024 * 1024

b, err := ParseBytes("1MB")
// uint64 1000 * 1000
```

`BytesIEC` and `BytesSI` can be used to format a byte value with human-readable string output.

```go
b := convert.BytesIEC(999)

fmt.Println(b)
// "999B"

b := convert.BytesIEC(999 * 1024)

fmt.Println(b)
// "999KiB"

b := convert.BytesIEC(999 * 1024 * 1024 * 1024 * 1024)

fmt.Println(b)
// "999TiB"

b := convert.BytesSI(999)

fmt.Println(b)
// "999B"

b := convert.BytesSI(999 * 1000)

fmt.Println(b)
// "999KB"

b := convert.BytesSI(999 * 1000 * 1000 * 1000 * 1000)

fmt.Println(b)
// "999TB"
```

# Examples

A few plugins using go-check:

* [check_prometheus](https://github.com/NETWAYS/check_prometheus)
* [check_system_basics](https://github.com/NETWAYS/check_system_basics)
* [check_logstash](https://github.com/NETWAYS/check_logstash)
* [check_sentinelone](https://github.com/NETWAYS/check_sentinelone)

# License

Copyright (c) 2020 [NETWAYS GmbH](mailto:info@netways.de)

This library is distributed under the GPL-2.0 or newer license found in the [COPYING](./COPYING)
file.
