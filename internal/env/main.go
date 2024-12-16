package env

import (
	"os"
	"path/filepath"
)

func LoadEnvVars() {
	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(exec)
	os.Setenv("JENKINS_HOME", path)
	os.Setenv("WORKSPACE", filepath.Join(path, "workspace"))

}
