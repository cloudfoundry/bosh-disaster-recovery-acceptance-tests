package runner

import (
	"io"
	"os/exec"
	"strings"

	"fmt"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunBoshCommandSuccessfullyWithFailureMessage(description string, writer io.Writer, config Config, args ...string) *gexec.Session {
	return RunCommandSuccessfullyWithFailureMessage(description, writer, getBoshBaseCommand(config), args...)
}

func RunCommandSuccessfullyWithFailureMessage(description string, writer io.Writer, cmd string, args ...string) *gexec.Session {
	session := runCommandWithStream(description, writer, cmd, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
	return session
}

func RunBoshCommand(description string, writer io.Writer, config Config, args ...string) *gexec.Session {
	return runCommandWithStream(description, writer, getBoshBaseCommand(config), args...)
}

func runCommandWithStream(description string, writer io.Writer, cmd string, args ...string) *gexec.Session {
	cmdToRunArgs := strings.Join(args, " ")
	cmdToRun := cmd + " " + cmdToRunArgs

	command := exec.Command("bash", "-c", cmdToRun)
	session, err := gexec.Start(command, writer, writer)

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(), "Command timed out: "+description)
	fmt.Fprintln(writer, "")
	return session
}

func getBoshBaseCommand(config Config) string {
	return fmt.Sprintf("bosh-cli "+
		"--environment=%s "+
		"--client=%s "+
		"--client-secret=%s "+
		"--ca-cert=%s ",
		config.BOSH.Host,
		config.BOSH.Client,
		config.BOSH.ClientSecret,
		config.BOSH.CACertPath)
}
