package utils

import (
	"os"
	"path"
	"runtime"
)

func PathEnable(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}


var runDir string
func GetRuntimeDir(name ...string) string {
	if runDir == "" {
		_, filename, _, _ := runtime.Caller(1)
		runDir, _ = path.Split(path.Dir(filename))
	}

	p := runDir

	if len(name) > 0 {
		p = path.Join(p, name[0])
	}
	return p
}
