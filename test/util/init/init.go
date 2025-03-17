package init

import (
	"os"
	"path"
	"runtime"
)

func init() {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled
	dir := path.Join(path.Dir(filename), "..", "..", "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
