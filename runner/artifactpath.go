package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var tmpMatcher = regexp.MustCompile(`: stdout \| (/tmp/\S+)\s`)
var tmpMatcherExistingJumpbox = regexp.MustCompile(`(/tmp/\S+)\s`)

func newArtifactPath(config Config) commonArtifactPath {
	if config.Jumpbox == nil {
		return newLocalArtifactPath()
	}
	return newJumpboxArtifactPath(config)
}

type commonArtifactPath interface {
	path() string
	cleanup()
	firstMatch(substr string) string
}

func newLocalArtifactPath() commonArtifactPath {
	return localArtifactPath{dir: GinkgoT().TempDir()}
}

type localArtifactPath struct {
	dir string
}

func (l localArtifactPath) path() string {
	return l.dir

}

func (l localArtifactPath) cleanup() {
	// Cleanup is automatic
}

func (l localArtifactPath) firstMatch(substr string) string {
	entries, err := os.ReadDir(l.dir)
	Expect(err).NotTo(HaveOccurred())
	for _, e := range entries {
		if strings.Contains(e.Name(), substr) {
			return filepath.Join(l.dir, e.Name())
		}
	}

	Fail(fmt.Sprintf("No match for %q in directory %q", substr, l.dir))
	return "" // Unreachable
}

func newJumpboxArtifactPath(config Config) commonArtifactPath {
	By("creating an artifact directory on the jumpbox")
	session := config.Jumpbox.Run("make temporary directory", config, "mktemp", "-d")
	if config.Jumpbox.HostIsSet() {
		tmpMatcher = tmpMatcherExistingJumpbox
	}
	matches := tmpMatcher.FindStringSubmatch(string(session.Out.Contents()))
	if len(matches) == 2 {
		return jumpboxArtifactPath{
			config: config,
			dir:    matches[1],
		}
	}

	Fail(fmt.Sprintf("failed to extract dir name from: %s", string(session.Out.Contents())))
	return jumpboxArtifactPath{} // Unreachable
}

type jumpboxArtifactPath struct {
	dir    string
	config Config
}

func (l jumpboxArtifactPath) path() string {
	return l.dir

}

func (l jumpboxArtifactPath) cleanup() {
	By("cleaning up the artifact directory on the jumpbox")
	l.config.Jumpbox.Run("cleanup", l.config, "rm", "-rf", l.dir)
}

func (l jumpboxArtifactPath) firstMatch(substr string) string {
	By("listing the contents of the artifact directory on the jumpbox")
	session := l.config.Jumpbox.Run("list", l.config, "find", l.dir, "-maxdepth 1")
	if l.config.Jumpbox.HostIsSet() {
		tmpMatcher = tmpMatcherExistingJumpbox
	}
	contents := tmpMatcher.FindAllStringSubmatch(string(session.Out.Contents()), -1)
	for _, m := range contents {
		if m[1] == l.dir {
			continue // skip the parent directory
		}

		if strings.Contains(m[1], substr) {
			return m[1]
		}
	}

	Fail(fmt.Sprintf("No match for %q in remote jumpbox directory %q", substr, l.dir))
	return "" // Unreachable
}
