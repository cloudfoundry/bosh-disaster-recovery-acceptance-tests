package runner

import (
	"os"
	"time"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/acceptance"
)

type Config struct {
	BOSH          BOSHConfig
	Credhub       CredhubConfig
	BBRBinaryPath string
	ArtifactPath  string
	StemcellSrc   string
	Timeout       time.Duration
}

type CloudConfig struct {
	DefaultVMType  string
	DefaultNetwork string
	DefaultAZ      string
}

type BOSHConfig struct {
	Host              string
	SSHUsername       string
	SSHPrivateKeyPath string
	Client            string
	ClientSecret      string
	CACertPath        string
	CloudConfig       CloudConfig
}

type CredhubConfig struct {
	CA           string
	Client       string
	ClientSecret string
	Server       string
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

	if integrationConfig.CredhubServer != "" {
		if err := os.Setenv("CREDHUB_PROXY", integrationConfig.GetCredhubProxy(privateKeyPath)); err != nil {
			return Config{}, err
		}
	}

	timeout := 50 * time.Minute
	if integrationConfig.TimeoutMinutes != 0 {
		timeout = time.Duration(integrationConfig.TimeoutMinutes) * time.Minute
	}

	defaultVMType := "default"
	if integrationConfig.DeploymentVMType != "" {
		defaultVMType = integrationConfig.DeploymentVMType
	}

	defaultNetwork := "default"
	if integrationConfig.DeploymentNetwork != "" {
		defaultNetwork = integrationConfig.DeploymentNetwork
	}

	defaultAZ := "z1"
	if integrationConfig.DeploymentAZ != "" {
		defaultAZ = integrationConfig.DeploymentAZ
	}

	return Config{
		BOSH: BOSHConfig{
			Host:              integrationConfig.Host,
			SSHUsername:       integrationConfig.SSHUsername,
			SSHPrivateKeyPath: privateKeyPath,
			Client:            integrationConfig.BOSHClient,
			ClientSecret:      integrationConfig.BOSHClientSecret,
			CACertPath:        caCertPath,
			CloudConfig: CloudConfig{
				DefaultVMType:  defaultVMType,
				DefaultNetwork: defaultNetwork,
				DefaultAZ:      defaultAZ,
			},
		},
		Credhub: CredhubConfig{
			CA:           integrationConfig.CredhubCACert,
			Client:       integrationConfig.CredhubClient,
			ClientSecret: integrationConfig.CredhubClientSecret,
			Server:       integrationConfig.CredhubServer,
		},
		StemcellSrc:   integrationConfig.StemcellSrc,
		BBRBinaryPath: bbrBinaryPath,
		ArtifactPath:  artifactDirPath,
		Timeout:       timeout,
	}, nil
}
