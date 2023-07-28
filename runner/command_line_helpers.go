package runner

import (
	"io"
	"os/exec"
	"strings"

	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunBoshCommand(description string, config Config, args ...string) *gexec.Session {
	return runCommandWithStream(description, GinkgoWriter, getBoshBaseCommand(config), args...)
}

func RunBoshCommandSuccessfullyWithFailureMessage(description string, config Config, args ...string) *gexec.Session {
	return RunCommandSuccessfullyWithFailureMessage(description, GinkgoWriter, getBoshBaseCommand(config), args...)
}

func RunCommandInDirectorVM(description string, config Config, cmd string, args ...string) *gexec.Session {
	return runCommandWithStreamInDirectorVM(description, GinkgoWriter, config, cmd, args...)
}

func RunCommandInDirectorVMSuccessfullyWithFailureMessage(description string, config Config, cmd string, args ...string) *gexec.Session {
	session := runCommandWithStreamInDirectorVM(description, GinkgoWriter, config, cmd, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
	return session
}

func RunCommandSuccessfullyWithFailureMessage(description string, writer io.Writer, cmd string, args ...string) *gexec.Session {
	session := runCommandWithStream(description, writer, cmd, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
	return session
}

func runCommandWithStreamInDirectorVM(description string, writer io.Writer, config Config, cmd string, args ...string) *gexec.Session {
	cmdToRunArgs := strings.Join(args, " ")
	cmdToRun := cmd + " " + cmdToRunArgs

	command := exec.Command(
		"ssh",
		config.BOSH.Host,
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"-l", config.BOSH.SSHUsername,
		"-i", config.BOSH.SSHPrivateKeyPath,
		cmdToRun,
	)
	session, err := gexec.Start(command, writer, writer)

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(), "Command timed out: "+description)
	fmt.Fprintln(writer, "")
	return session
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
	return fmt.Sprintf("bosh "+
		"--environment=%s "+
		"--client=%s "+
		"--client-secret=%s "+
		"--ca-cert=%s ",
		config.BOSH.Host,
		config.BOSH.Client,
		config.BOSH.ClientSecret,
		config.BOSH.CACertPath)
}
