package testcases

import (
  "fmt"

  "code.cloudfoundry.org/credhub-cli/credhub"
  "code.cloudfoundry.org/credhub-cli/credhub/auth"
  "code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
  . "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
  . "github.com/onsi/gomega"
)

const CredentialName = "some-password"
const CredentialValue = "some-value"

type CredhubTestcase struct {}

func (t CredhubTestcase) Name() string {
  return "credhub_testcase"
}

func (t CredhubTestcase) BeforeBackup(config Config) {
  credhubClient, err := t.credhubClient(config)
  Expect(err).ToNot(HaveOccurred())

  _, err = credhubClient.SetPassword(CredentialName, values.Password(CredentialValue))
  Expect(err).ToNot(HaveOccurred())
}

func (t CredhubTestcase) AfterBackup(config Config) {
  credhubClient, err := t.credhubClient(config)
  Expect(err).ToNot(HaveOccurred())

  err = credhubClient.Delete(CredentialName)
  Expect(err).ToNot(HaveOccurred())
}

func (t CredhubTestcase) AfterRestore(config Config) {
  credhubClient, err := t.credhubClient(config)
  Expect(err).ToNot(HaveOccurred())

  credential, err := credhubClient.GetLatestPassword(CredentialName)
  Expect(err).ToNot(HaveOccurred())

  Expect(string(credential.Value)).To(Equal(CredentialValue))
}

func (t CredhubTestcase) Cleanup(config Config) {
  credhubClient, err := t.credhubClient(config)
  Expect(err).ToNot(HaveOccurred())

  err = credhubClient.Delete(CredentialName)
  Expect(err).ToNot(HaveOccurred())
}

func (t CredhubTestcase) credhubClient(config Config) (*credhub.CredHub, error) {
  return credhub.New(fmt.Sprintf("https://%s:8844", config.BOSH.Host),
    credhub.CaCerts(config.Credhub.CA),
    credhub.Auth(
      auth.UaaClientCredentials(config.Credhub.Client, config.Credhub.ClientSecret),
    ))
}
