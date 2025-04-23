package check

import (
	"os"
	"testing"
)

func ExampleConfig() {
	config := NewConfig()
	config.Name = "check_test"
	config.Readme = `Test Plugin`
	config.Version = "1.0.0"

	_ = config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")

	config.ParseArguments()

	// Some checking should be done here

	Exitf(OK, "Everything is fine - answer=%d", 42)

	// Output: [OK] - Everything is fine - answer=42
	// would exit with code 0
}

type ConfigForTesting struct {
	Auth            string `env:"AUTH"`
	Bearer          string `env:"EXAMPLE"`
	OneMoreThanTags string
}

func TestLoadFromEnv(t *testing.T) {
	c := ConfigForTesting{}

	err := os.Setenv("EXAMPLE", "foobar")
	defer os.Unsetenv("EXAMPLE") // just to not create any side effects

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	LoadFromEnv(&c)

	if c.Bearer != "foobar" {
		t.Fatalf("expected %v, got %v", c.Bearer, "foobar")
	}

	if c.Auth != "" {
		t.Fatalf("expected %v, got %v", c.Auth, "")
	}

	if c.OneMoreThanTags != "" {
		t.Fatalf("expected %v, got %v", c.OneMoreThanTags, "")
	}
}
