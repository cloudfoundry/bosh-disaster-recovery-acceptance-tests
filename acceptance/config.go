package acceptance

import (
	"os"
)

type IntegrationConfig struct {
	Host                 string `json:"bosh_host"`
	TimeoutMinutes       int64  `json:"timeout_in_minutes"`
	BOSHClient           string `json:"bosh_client"`
	BOSHClientSecret     string `json:"bosh_client_secret"`
	BOSHCACert           string `json:"bosh_ca_cert"`
	BOSHAgentEndpoint    string `json:"bosh_agent_endpoint"`
	BOSHAgentCertificate string `json:"bosh_agent_certificate"`
	CredhubClient        string `json:"credhub_client"`
	CredhubClientSecret  string `json:"credhub_client_secret"`
	CredhubCACert        string `json:"credhub_ca_cert"`
	CredhubServer        string `json:"credhub_server"`
	DeploymentVMType     string `json:"deployment_vm_type"`
	DeploymentNetwork    string `json:"deployment_network"`
	DeploymentAZ         string `json:"deployment_az"`
	StemcellSrc          string `json:"stemcell_src"`
	JumpboxHost          string `json:"jumpbox_host"`
	JumpboxUser          string `json:"jumpbox_user"`
	JumpboxPrivKey       string `json:"jumpbox_privkey"`
	JumpboxPubKey        string `json:"jumpbox_pubkey"`
}

func (i IntegrationConfig) CACertPath() (string, error) {
	return i.writeSecretToFile(i.BOSHCACert)
}

func (i IntegrationConfig) writeSecretToFile(contents string) (string, error) {
	file, err := os.CreateTemp("", "b-drats")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(file.Name(), []byte(contents), 0400)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
