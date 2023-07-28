package fixtures

import (
	"path"
	"runtime"
	"time"
)

var EventuallyTimeout = 10 * time.Minute
var EventuallyRetryInterval = 30 * time.Second

func Path(relativePath string) string {
	return path.Join(currentTestDir(), "../fixtures", relativePath)
}

func currentTestDir() string {
	_, filePath, _, _ := runtime.Caller(1)
	return path.Dir(filePath)
}
