package acceptance

import "io/ioutil"

type IntegrationConfig struct {
	Host           string `json:"bosh_host"`
	SSHUsername    string `json:"bosh_ssh_username"`
	SSHPrivateKey  string `json:"bosh_ssh_private_key"`
	TimeoutMinutes int    `json:"timeout_in_minutes"`
}

func (i IntegrationConfig) SSHPrivateKeyPath() (string, error) {
	file, err := ioutil.TempFile("", "b-drats")
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(file.Name(), []byte(i.SSHPrivateKey), 0400)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
