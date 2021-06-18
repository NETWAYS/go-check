go-check
========

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/NETWAYS/go-check?label=version)](https://github.com/NETWAYS/go-check/releases)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/NETWAYS/go-check)
[![Test Status](https://github.com/NETWAYS/go-check/workflows/Go/badge.svg)](https://github.com/NETWAYS/go-check/actions?query=workflow%3AGo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/NETWAYS/go-check)
![GitHub](https://img.shields.io/github/license/NETWAYS/go-check?color=green)

go-check is a library to help with development of monitoring plugins for tools like Icinga.

## Usage

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
}
```

```
OK - Everything is fine - answer=42
would exit with code 0
```

See the [documentation on pkg.go.dev](https://pkg.go.dev/github.com/NETWAYS/go-check) for more details and examples.

## Plugins

A few plugins using go-check:

* [check_cloud_aws](https://github.com/NETWAYS/check_cloud_aws)
* [check_cloud_azure](https://github.com/NETWAYS/check_cloud_azure)
* [check_cloud_gcp](https://github.com/NETWAYS/check_cloud_gcp)
* [check_sentinelone](https://github.com/NETWAYS/check_sentinelone)
* [check_sophos_central](https://github.com/NETWAYS/check_sophos_central)

## License

Copyright (c) 2020 [NETWAYS GmbH](mailto:info@netways.de)

This library is distributed under the GPL-2.0 or newer license found in the [COPYING](./COPYING)
file.
