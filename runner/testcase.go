package runner

type TestCase interface {
	Name() string
	BeforeBackup()
	AfterBackup()
	AfterRestore()
	Cleanup()
}