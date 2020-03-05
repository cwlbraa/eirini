package eirini

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
)

const (
	//Environment Variable Names
	EnvDownloadURL        = "DOWNLOAD_URL"
	EnvBuildpacks         = "BUILDPACKS"
	EnvDropletUploadURL   = "DROPLET_UPLOAD_URL"
	EnvAppID              = "APP_ID"
	EnvStagingGUID        = "STAGING_GUID"
	EnvCompletionCallback = "COMPLETION_CALLBACK"
	EnvEiriniAddress      = "EIRINI_ADDRESS"

	EnvPodName              = "POD_NAME"
	EnvCFInstanceIP         = "CF_INSTANCE_IP"
	EnvCFInstanceInternalIP = "CF_INSTANCE_INTERNAL_IP"
	EnvCFInstanceAddr       = "CF_INSTANCE_ADDR"
	EnvCFInstancePort       = "CF_INSTANCE_PORT"
	EnvCFInstancePorts      = "CF_INSTANCE_PORTS"

	RecipeBuildPacksDir    = "/var/lib/buildpacks"
	RecipeBuildPacksName   = "recipe-buildpacks"
	RecipeWorkspaceDir     = "/recipe_workspace"
	RecipeWorkspaceName    = "recipe-workspace"
	RecipeOutputName       = "staging-output"
	RecipeOutputLocation   = "/out"
	RecipePacksBuilderPath = "/packs/builder"

	AppMetricsEmissionIntervalInSecs = 15

	CertsMountPath  = "/etc/config/certs"
	CertsVolumeName = "certs-volume"

	CACertName       = "internal-ca-cert"
	CCAPICertName    = "cc-server-crt"
	CCAPIKeyName     = "cc-server-crt-key"
	EiriniClientCert = "eirini-client-crt"
	EiriniClientKey  = "eirini-client-crt-key"
)

type Config struct {
	Properties Properties `yaml:"opi"`
}

type KubeConfig struct {
	Namespace  string `yaml:"app_namespace"`
	ConfigPath string `yaml:"kube_config_path"`
}

type Properties struct {
	ClientCAPath   string `yaml:"client_ca_path"`
	ServerCertPath string `yaml:"server_cert_path"`
	ServerKeyPath  string `yaml:"server_key_path"`
	TLSPort        int    `yaml:"tls_port"`

	CCUploaderSecretName string `yaml:"cc_uploader_secret_name"`
	CCUploaderCertPath   string `yaml:"cc_uploader_cert_path"`
	CCUploaderKeyPath    string `yaml:"cc_uploader_key_path"`

	ClientCertsSecretName string `yaml:"client_certs_secret_name"`
	ClientCertPath        string `yaml:"client_cert_path"`
	ClientKeyPath         string `yaml:"client_key_path"`

	CACertSecretName string `yaml:"ca_cert_secret_name"`
	CACertPath       string `yaml:"ca_cert_path"`

	RegistryAddress                  string `yaml:"registry_address"`
	RegistrySecretName               string `yaml:"registry_secret_name"`
	EiriniAddress                    string `yaml:"eirini_address"`
	DownloaderImage                  string `yaml:"downloader_image"`
	UploaderImage                    string `yaml:"uploader_image"`
	ExecutorImage                    string `yaml:"executor_image"`
	AppMetricsEmissionIntervalInSecs int    `yaml:"app_metrics_emission_interval_in_secs"`

	CCCertPath string `yaml:"cc_cert_path"`
	CCKeyPath  string `yaml:"cc_key_path"`
	CCCAPath   string `yaml:"cc_ca_path"`

	RootfsVersion string `yaml:"rootfs_version"`
	DiskLimitMB   int64  `yaml:"disk_limit_mb"`

	KubeConfig `yaml:",inline"`
}

type EventReporterConfig struct {
	CcInternalAPI string `yaml:"cc_internal_api"`
	CCCertPath    string `yaml:"cc_cert_path"`
	CCKeyPath     string `yaml:"cc_key_path"`
	CCCAPath      string `yaml:"cc_ca_path"`

	KubeConfig `yaml:",inline"`
}

type RouteEmitterConfig struct {
	NatsPassword string `yaml:"nats_password"`
	NatsIP       string `yaml:"nats_ip"`
	NatsPort     int    `yaml:"nats_port"`

	KubeConfig `yaml:",inline"`
}

type MetricsCollectorConfig struct {
	LoggregatorAddress  string `yaml:"loggregator_address"`
	LoggregatorCertPath string `yaml:"loggergator_cert_path"`
	LoggregatorKeyPath  string `yaml:"loggregator_key_path"`
	LoggregatorCAPath   string `yaml:"loggregator_ca_path"`

	AppMetricsEmissionIntervalInSecs int `yaml:"app_metrics_emission_interval_in_secs"`

	KubeConfig `yaml:",inline"`
}

type StagingReporterConfig struct {
	EiriniCertPath string `yaml:"eirini_cert_path"`
	EiriniKeyPath  string `yaml:"eirini_key_path"`
	CAPath         string `yaml:"ca_path"`

	KubeConfig `yaml:",inline"`
}

//go:generate counterfeiter . Stager
type Stager interface {
	Stage(string, cf.StagingRequest) error
	CompleteStaging(*models.TaskCallbackResponse) error
}

type StagerConfig struct {
	EiriniAddress   string
	DownloaderImage string
	UploaderImage   string
	ExecutorImage   string
}

//go:generate counterfeiter . Extractor
type Extractor interface {
	Extract(src, targetDir string) error
}

//go:generate counterfeiter . Bifrost
type Bifrost interface {
	Transfer(ctx context.Context, request cf.DesireLRPRequest) error
	List(ctx context.Context) ([]*models.DesiredLRPSchedulingInfo, error)
	Update(ctx context.Context, update cf.UpdateDesiredLRPRequest) error
	Stop(ctx context.Context, identifier opi.LRPIdentifier) error
	StopInstance(ctx context.Context, identifier opi.LRPIdentifier, index uint) error
	GetApp(ctx context.Context, identifier opi.LRPIdentifier) (*models.DesiredLRP, error)
	GetInstances(ctx context.Context, identifier opi.LRPIdentifier) ([]*cf.Instance, error)
}

func GetInternalServiceName(appName string) string {
	//Prefix service as the appName could start with numerical characters, which is not allowed
	return fmt.Sprintf("cf-%s", appName)
}

func GetInternalHeadlessServiceName(appName string) string {
	//Prefix service as the appName could start with numerical characters, which is not allowed
	return fmt.Sprintf("cf-%s-headless", appName)
}
