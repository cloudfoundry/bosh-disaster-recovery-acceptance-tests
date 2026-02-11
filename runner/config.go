package runner

import (
	"time"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/acceptance"
)

type Config struct {
	BOSH          BOSHConfig
	Credhub       CredhubConfig
	BBRBinaryPath string
	ArtifactPath  string
	StemcellSrc   string
	StemcellOs    string
	Timeout       time.Duration
	Jumpbox       *jumpbox
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

type jumpbox struct {
	host    string
	privkey string
	user    string
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

	var jb *jumpbox
	if integrationConfig.JumpboxHost != "" && integrationConfig.JumpboxPrivKey != "" && integrationConfig.JumpboxUser != "" {
		jb = &jumpbox{
			privkey: integrationConfig.JumpboxPrivKey,
			host:    integrationConfig.JumpboxHost,
			user:    integrationConfig.JumpboxUser,
		}
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
		StemcellOs:    integrationConfig.StemcellOs,
		BBRBinaryPath: bbrBinaryPath,
		ArtifactPath:  artifactDirPath,
		Timeout:       timeout,
		Jumpbox:       jb,
	}, nil
}
