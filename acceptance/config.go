package acceptance

import "io/ioutil"

type IntegrationConfig struct {
	Host                string `json:"bosh_host"`
	SSHUsername         string `json:"bosh_ssh_username"`
	SSHPrivateKey       string `json:"bosh_ssh_private_key"`
	TimeoutMinutes      int64  `json:"timeout_in_minutes"`
	BOSHClient          string `json:"bosh_client"`
	BOSHClientSecret    string `json:"bosh_client_secret"`
	CredhubClient       string `json:"credhub_client"`
	CredhubClientSecret string `json:"credhub_client_secret"`
	BOSHCACert          string `json:"bosh_ca_cert"`
	CredhubCACert       string `json:"credhub_ca_cert"`
	DeploymentVMType    string `json:"deployment_vm_type"`
	DeploymentNetwork   string `json:"deployment_network"`
	DeploymentAZ        string `json:"deployment_az"`
	StemcellSrc         string `json:"stemcell_src"`
}

func (i IntegrationConfig) SSHPrivateKeyPath() (string, error) {
	return i.writeSecretToFile(i.SSHPrivateKey)
}

func (i IntegrationConfig) CACertPath() (string, error) {
	return i.writeSecretToFile(i.BOSHCACert)
}

func (i IntegrationConfig) writeSecretToFile(contents string) (string, error) {
	file, err := ioutil.TempFile("", "b-drats")
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(file.Name(), []byte(contents), 0400)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
