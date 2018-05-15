package runner

import (
	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/acceptance"
	"time"
)

type Config struct {
	BOSH          BOSHConfig
	BBRBinaryPath string
	ArtifactPath  string
	Timeout       time.Duration
}

type BOSHConfig struct {
	Host              string
	SSHUsername       string
	SSHPrivateKeyPath string
	Client            string
	ClientSecret      string
	CACertPath        string
}

func NewConfig(integrationConfig acceptance.IntegrationConfig, bbrBinaryPath, artifactDirPath string) (Config, error) {
	privateKeyPath, err := integrationConfig.SSHPrivateKeyPath()
	if err != nil {
		return Config{}, err
	}

	caCertPath, err := integrationConfig.CACertPath()
	if err != nil {
		return Config{}, err
	}

	timeout := 30 * time.Minute
	if integrationConfig.TimeoutMinutes != 0 {
		timeout = time.Duration(integrationConfig.TimeoutMinutes) * time.Minute
	}

	return Config{
		BOSH: BOSHConfig{
			Host:              integrationConfig.Host,
			SSHUsername:       integrationConfig.SSHUsername,
			SSHPrivateKeyPath: privateKeyPath,
			Client:            integrationConfig.BOSHClient,
			ClientSecret:      integrationConfig.BOSHClientSecret,
			CACertPath:        caCertPath,
		},
		BBRBinaryPath: bbrBinaryPath,
		ArtifactPath:  artifactDirPath,
		Timeout:       timeout,
	}, nil
}
