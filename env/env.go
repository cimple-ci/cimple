package env

func Cimple() map[string]string {
	env := make(map[string]string)
	env["CIMPLE_VERSION"] = "0.0.1"
	return env
}
