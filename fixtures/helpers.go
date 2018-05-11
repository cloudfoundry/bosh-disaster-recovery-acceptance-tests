package fixtures

import (
	"path"
	"runtime"
)

func Path(relativePath string) string {
	return path.Join(currentTestDir(), "../fixtures", relativePath)
}

func currentTestDir() string {
	_, filePath, _, _ := runtime.Caller(1)
	return path.Dir(filePath)
}
