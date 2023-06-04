package testing

import (
	"os"
	"path"
	"runtime"
)

var here string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	here = path.Dir(filename)
	root := path.Join(here, "..")
	os.Chdir(root)
}

func Path(name string) string {
	return path.Join(here, name)
}
