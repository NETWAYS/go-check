# Introdcution
In the following, several different examples of Check Plugins will be
presetend to give some hints as to how this libary should be used.

# Minimal check plugin

Lets start here with minimalistig example

```go
package main

import (
	"github.com/NETWAYS/go-check"
)

func test_something(input string) (check.Status, string) {
  // This would be the place to do the actual test
	return check.OK, input
}

func main() {
	// The core test
	return_code, return_text := test_something("foo")

	// Handle exit with the desired exit code
  // check.Exit will prefix the output with a [OK] for a return_code 0, [WARNING] for 1, etc.
	check.Exit(return_code, return_text)
}
```

Here the `go-check` library is effectively only used for some types (name the `Status` type,
symbolizing one of the four possible results) and the `Exit` function.


# Example with basic config
```go
package main

import (
	"github.com/NETWAYS/go-check"
)

func test_something(foo string) (check.Status, string) {
  // This would be the place to do the actual test
	return check.OK, foo
}

func main() {
	// Global configuration of the plugin
	config := check.NewConfig()
	config.Name = "check_test"
	config.Readme = `Test Plugin`
	config.Version = "1.0.0"

	// Command line arguments
	foo := config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")

	config.ParseArguments()

	// The core test
	return_code, return_text := test_something(*foo)

	// Handle exit with the desired exit code
	check.Exit(return_code, return_text)
	// Output:
	// [OK] - Everything is fine - answer=42
}
```

This slightly extented example adds some basic configuration functionality.
Let us assume, that the above content is in a file `main.go` in a module
`go-check-plugin`.

```
# go build
# ./go-check-plugin --help
Usage of check_test

Test Plugin

Arguments:
  -H, --hostname string   Hostname to check (default "localhost")
  -t, --timeout int       Abort the check after n seconds (default 30)
  -d, --debug             Enable debug mode
  -v, --verbose           Enable verbose mode
  -V, --version           Print version and exit
[UNKNOWN] - pflag: help requested (*errors.errorString)
```

# Example with multiple tests

Now on to something more complex (and realistic).
Monitoring Plugins often execute several different
tests; a very common example would be the monitoring of filesystems
where each filesystem is tested individually and where each failed test
is cause for an alarm.

The `go-check` library has some useful structures to map this problem
into code and keep it relatively simple.

Each individual test result can be put into a `PartialResult`
and the `Overall` objects collects and combines them:

```go
package main

import (
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
)

func test_1(foo string) (check.Status, string) {
	return check.OK, foo
}

func test_2(foo string) (check.Status, string) {
	return check.OK, foo
}

func main() {
	// Global configuration of the plugin
	config := check.NewConfig()
	config.Name = "check_test"
	config.Readme = `Test Plugin`
	config.Version = "1.0.0"

	// Command line arguments
	foo := config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")

	config.ParseArguments()

	// Create the overall object
	var overall result.Overall

	// Test 1
	return_code1, return_text1 := test_1(*foo)
	overall.Add(return_code1, return_text1)

  // Test 2
	return_code2, return_text2 := test_2(*foo)
	overall.Add(return_code2, return_text2)

	// Handle exit
	check.Exit(overall.GetStatus(), overall.GetOutput())
}
```

This gives us:
```
# ./go-check-plugin
[OK] - states: ok=2
\_ [OK] localhost
\_ [OK] localhost
```

The `Overall` object is, by design, a singleton and must only exist once per
Monitoring Plugin.
