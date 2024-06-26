package runner

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunBoshCommand(description string, config Config, args ...string) *gexec.Session {
	return RunCommandWithStream(description, GinkgoWriter, getBoshBaseCommand(config), args...)
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

func RunBBRCommand(description string, config Config, args ...string) *gexec.Session {
	commandWithBoshAllProxy := fmt.Sprintf("%s %s", getBoshAllProxy(config), config.BBRBinaryPath)
	return RunCommandWithStream(description, os.Stdout, commandWithBoshAllProxy, args...)
}

func RunBBRCommandSuccessfullyWithFailureMessage(description string, config Config, args ...string) {
	session := RunBBRCommand(description, config, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
}

func RunCommandSuccessfullyWithFailureMessage(description string, writer io.Writer, cmd string, args ...string) *gexec.Session {
	session := RunCommandWithStream(description, writer, cmd, args...)
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

func RunCommandWithStream(description string, writer io.Writer, cmd string, args ...string) *gexec.Session {
	cmdToRunArgs := strings.Join(args, " ")
	cmdToRun := cmd + " " + cmdToRunArgs

	command := exec.Command("bash", "-c", cmdToRun)
	session, err := gexec.Start(command, writer, writer)

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(), "Command timed out: "+description)
	fmt.Fprintln(writer, "")
	return session
}
func getBoshAllProxy(config Config) string {
	// ssh+socks5://ubuntu@34.72.88.156:22?private-key=/tmp/tmp.bBURxmHm5j
	if config.Jumpbox != nil && config.Jumpbox.host != "" {
		keyPath, err := jumpboxKeyFile(config.Jumpbox.privkey)
		if err != nil {
			Fail("failed writing jumphost keyfile")
		}
		return fmt.Sprintf("BOSH_ALL_PROXY=ssh+socks5://%s@%s:22?private-key=%s", config.Jumpbox.user, config.Jumpbox.host, keyPath)
	}
	return ""
}
func getBoshBaseCommand(config Config) string {

	return fmt.Sprintf(
		"%s "+
			"bosh "+
			"--environment=%s "+
			"--client=%s "+
			"--client-secret=%s "+
			"--ca-cert=%s ",
		getBoshAllProxy(config),
		config.BOSH.Host,
		config.BOSH.Client,
		config.BOSH.ClientSecret,
		config.BOSH.CACertPath)
}

func jumpboxKeyFile(privateKey string) (string, error) {
	d, err := os.MkdirTemp("", "drats")
	if err != nil {
		return "", err
	}
	fname := filepath.Join(d, "key")
	err = os.WriteFile(fname, []byte(privateKey), 0400)
	return fname, err
}
