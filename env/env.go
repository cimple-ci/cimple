package env

import "os"
import "strings"

type CimpleEnvironment struct {
	Version string
}

func Cimple() *CimpleEnvironment {
	return &CimpleEnvironment{
		Version: "0.0.1",
	}
}

func EnvironmentVariables() map[string]string {
	vars := make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		vars[pair[0]] = pair[1]
	}

	return vars
}
