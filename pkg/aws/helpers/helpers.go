package helpers

import (
	"os"
	"os/exec"
)

type OSExecIface interface {
	LookPath(file string) (string, error)
}

type OSIface interface {
	LookupEnv(envVar string) (string, bool)
	Setenv(envVar, value string) error
}

type Helpers struct{}

// LookupEnv looks up an environment variable
func (Helpers) LookupEnv(envVar string) (string, bool) {
	return os.LookupEnv(envVar)
}

// Setenv sets an environment variable
func (Helpers) Setenv(envVar, value string) error {
	return os.Setenv(envVar, value)
}

// LookPath looks up a path
func (Helpers) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
