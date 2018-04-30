package runner

import (
	"io"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunCommandSuccessfullyWithFailureMessage(commandDescription, cmd string, args ...string) *gexec.Session {
	session := runCommandWithStream(commandDescription, GinkgoWriter, GinkgoWriter, cmd, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+commandDescription)
	return session
}

func runCommandWithStream(commandDescription string, stdout, stderr io.Writer, cmd string, args ...string) *gexec.Session {
	cmdToRunArgs := strings.Join(args, " ")
	cmdToRun := cmd + " " + cmdToRunArgs

	command := exec.Command("bash", "-c", cmdToRun)
	session, err := gexec.Start(command, stdout, stderr)

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(), "Command timed out: "+commandDescription)
	return session
}
