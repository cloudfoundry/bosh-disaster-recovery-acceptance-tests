package runner

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os"
	"path"
	"strings"
	"text/template"
)

const deployment = "jumpbox"

// jumpbox is a helper to deploy, clean up, and run the BBR command on a jumpbox
// By running the BBR command on a jumpbox (in the same data center as BOSH), we
// speed up the backup and restore significantly because we avoid having to copy
// large data files into Concourse containers, which may not even be located on
// the same continent as BOSH.
type jumpbox struct {
	network string
	vmType  string
	az      string
}

func (j *jumpbox) Run(description string, config Config, args ...string) *gexec.Session {
	if j == nil {
		Fail("jumpbox not initialised")
	}

	params := []string{"-d", deployment, "ssh", "jumpbox", "--", "sudo"}
	params = append(params, strings.Split(strings.Join(args, " "), " ")...)

	return RunBoshCommand(description, config, params...)
}

func (j *jumpbox) RunBBR(description string, config Config, args ...string) *gexec.Session {
	args = append([]string{"/usr/local/bin/bbr"}, args...)
	return j.Run(description, config, args...)
}

func (j *jumpbox) Deploy(config Config) {
	if j == nil {
		GinkgoWriter.Println("Jumpbox not configured")
		return
	}

	By("deploying a jumpbox")
	RunBoshCommandSuccessfullyWithFailureMessage("deploying a jumpbox", config, "-n", "-d", deployment, "deploy", j.manifest())

	By("copying across the BBR command")
	RunBoshCommandSuccessfullyWithFailureMessage("copying across the BBR command", config, "-d", deployment, "scp", config.BBRBinaryPath, "jumpbox:/tmp/bbr")
	By("moving the BBR command into place")
	j.Run("moving the BBR command into place", config, "mv", "/tmp/bbr", "/usr/local/bin/bbr")
	By("Making the BBR command executable")
	j.Run("Making the BBR command executable", config, "chmod", "+x", "/usr/local/bin/bbr")

	keyPath := path.Join("/tmp", path.Base(config.BOSH.SSHPrivateKeyPath))
	By("copying across the BOSH private key")
	RunBoshCommandSuccessfullyWithFailureMessage("copying across the BOSH private key", config, "-d", deployment, "scp", config.BOSH.SSHPrivateKeyPath, fmt.Sprintf("jumpbox:%s", keyPath))
	By("setting the correct permissions on the BOSH private key")
	j.Run("setting the correct permissions on the BOSH private key", config, "chmod", "600", keyPath)
}

func (j *jumpbox) Cleanup(config Config) {
	if j == nil {
		GinkgoWriter.Println("Jumpbox cleanup not needed")
		return
	}

	By("cleaning up the jumpbox")
	RunBoshCommandSuccessfullyWithFailureMessage("cleaning up the jumpbox", config, "-n", "-d", deployment, "delete-deployment")
}

func (j *jumpbox) manifest() string {
	t, err := template.New("manifest").Parse(`
name: jumpbox

releases: []

update:
  canaries: 1
  max_in_flight: 1
  canary_watch_time: 1000-30000
  update_watch_time: 1000-30000

stemcells:
  - alias: default
    os: ubuntu-jammy
    version: latest

instance_groups:
  - name: jumpbox
    azs: [{{.AZ}}]
    instances: 1
    vm_type: {{.VMType}}
    stemcell: default
    networks:
      - name: {{.Network}}
    jobs: []
`)
	Expect(err).NotTo(HaveOccurred())

	fh, err := os.CreateTemp(GinkgoT().TempDir(), "manifest.yml")
	Expect(err).NotTo(HaveOccurred())
	defer fh.Close()

	mapping := map[string]string{
		"AZ":      j.az,
		"VMType":  j.vmType,
		"Network": j.network,
	}

	Expect(t.Execute(fh, mapping)).To(Succeed())

	return fh.Name()
}
