package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:staticcheck
	. "github.com/onsi/gomega"    //nolint:staticcheck
	"github.com/onsi/gomega/gexec"
)

func RunBoshCommand(description string, config Config, args ...string) *gexec.Session {
	return RunCommandWithStream(description, getBoshBaseCommand(config), args...)
}

func RunBoshCommandSuccessfullyWithFailureMessage(description string, config Config, args ...string) *gexec.Session {
	return RunCommandSuccessfullyWithFailureMessage(description, getBoshBaseCommand(config), args...)
}

func RunCommandInDirectorVM(description string, config Config, cmd string, args ...string) *gexec.Session {
	return runCommandWithStreamInDirectorVM(description, config, cmd, args...)
}

func RunCommandInDirectorVMSuccessfullyWithFailureMessage(description string, config Config, cmd string, args ...string) *gexec.Session {
	session := runCommandWithStreamInDirectorVM(description, config, cmd, args...)
	Expect(session).To(gexec.Exit(0), fmt.Sprintf("Command errored: %s", description))
	return session
}

func RunBBRCommand(description string, config Config, args ...string) *gexec.Session {
	commandWithBoshAllProxy := fmt.Sprintf("%s %s", getBoshAllProxy(config), config.BBRBinaryPath)
	return RunCommandWithStream(description, commandWithBoshAllProxy, args...)
}

func RunBBRCommandSuccessfullyWithFailureMessage(description string, config Config, args ...string) {
	session := RunBBRCommand(description, config, args...)
	Expect(session).To(gexec.Exit(0), fmt.Sprintf("Command errored: %s", description))
}

func RunCommandSuccessfullyWithFailureMessage(description string, cmd string, args ...string) *gexec.Session {
	session := RunCommandWithStream(description, cmd, args...)
	Expect(session).To(gexec.Exit(0), fmt.Sprintf("Command errored: %s", description))
	return session
}

func runCommandWithStreamInDirectorVM(description string, config Config, cmd string, args ...string) *gexec.Session {
	cmdWithArgs := append([]string{cmd}, args...)

	session, err := execCommand(
		"ssh",
		config.BOSH.Host,
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
		"-l", config.BOSH.SSHUsername,
		"-i", config.BOSH.SSHPrivateKeyPath,
		strings.Join(cmdWithArgs, " "),
	)

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(), fmt.Sprintf("Command timed out: %s", description))
	return session
}

func RunCommandWithStream(description string, cmd string, args ...string) *gexec.Session {
	cmdWithArgs := append([]string{cmd}, args...)

	session, err := execCommand("bash", "-c", strings.Join(cmdWithArgs, " "))

	Expect(err).ToNot(HaveOccurred())
	Eventually(session).Should(gexec.Exit(), fmt.Sprintf("Command timed out: %s", description))
	return session
}

func execCommand(cmd string, args ...string) (*gexec.Session, error) {
	session, err := gexec.Start(exec.Command(cmd, args...), GinkgoWriter, GinkgoWriter)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func getBoshAllProxy(config Config) string {
	// ssh+socks5://ubuntu@34.72.88.156:22?private-key=/tmp/tmp.bBURxmHm5j
	if config.Jumpbox != nil && config.Jumpbox.host != "" {
		keyPath, err := jumpboxKeyFile(config.Jumpbox.privkey)
		if err != nil {
			Fail("failed writing jumpbox keyfile")
		}
		return fmt.Sprintf("BOSH_ALL_PROXY=ssh+socks5://%s@%s:22?private-key=%s", config.Jumpbox.user, config.Jumpbox.host, keyPath)
	}
	return ""
}
func getBoshBaseCommand(config Config) string {
	return fmt.Sprintf(
		"%s bosh --environment=%s --client=%s --client-secret=%s --ca-cert=%s",
		getBoshAllProxy(config),
		config.BOSH.Host,
		config.BOSH.Client,
		config.BOSH.ClientSecret,
		config.BOSH.CACertPath,
	)
}

func jumpboxKeyFile(privateKey string) (string, error) {
	fileName := filepath.Join(GinkgoT().TempDir(), "key")
	err := os.WriteFile(fileName, []byte(privateKey), 0400)
	return fileName, err
}
