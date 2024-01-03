
// Code generated by go generate; DO NOT EDIT.
package param

import (
	"time"
)

type config struct {
	Cache struct {
		DataLocation string
		EnableVoms bool
		ExportLocation string
		Port int
		XRootDPrefix string
	}
	Client struct {
		DisableHttpProxy bool
		DisableProxyFallback bool
		MinimumDownloadSpeed int
		SlowTransferRampupTime int
		SlowTransferWindow int
		StoppedTransferTimeout int
	}
	ConfigDir string
	Debug bool
	Director struct {
		CacheResponseHostnames []string
		DefaultResponse string
		GeoIPLocation string
		MaxMindKeyFile string
		OriginResponseHostnames []string
	}
	DisableHttpProxy bool
	DisableProxyFallback bool
	Federation struct {
		DirectorUrl string
		DiscoveryUrl string
		JwkUrl string
		NamespaceUrl string
		RegistryUrl string
		TopologyNamespaceUrl string
		TopologyReloadInterval time.Duration
	}
	GeoIPOverrides interface{}
	Issuer struct {
		AuthenticationSource string
		AuthorizationTemplates interface{}
		GroupFile string
		GroupRequirements []string
		GroupSource string
		OIDCAuthenticationRequirements interface{}
		OIDCAuthenticationUserClaim string
		QDLLocation string
		ScitokensServerLocation string
		TomcatLocation string
	}
	IssuerKey string
	Logging struct {
		Level string
		LogLocation string
	}
	MinimumDownloadSpeed int
	Monitoring struct {
		AggregatePrefixes []string
		DataLocation string
		MetricAuthorization bool
		PortHigher int
		PortLower int
		TokenExpiresIn time.Duration
		TokenRefreshInterval time.Duration
	}
	OIDC struct {
		AuthorizationEndpoint string
		ClientID string
		ClientIDFile string
		ClientRedirectHostname string
		ClientSecretFile string
		DeviceAuthEndpoint string
		Issuer string
		TokenEndpoint string
		UserInfoEndpoint string
	}
	Origin struct {
		EnableCmsd bool
		EnableDirListing bool
		EnableIssuer bool
		EnableUI bool
		EnableVoms bool
		ExportVolume string
		Mode string
		Multiuser bool
		NamespacePrefix string
		S3AccessKeyfile string
		S3Bucket string
		S3Region string
		S3SecretKeyfile string
		S3ServiceName string
		S3ServiceUrl string
		ScitokensDefaultUser string
		ScitokensMapSubject bool
		ScitokensNameMapFile string
		ScitokensRestrictedPaths []string
		ScitokensUsernameClaim string
		SelfTest bool
		Url string
		XRootDPrefix string
	}
	Plugin struct {
		Token string
	}
	Registry struct {
		AdminUsers []string
		DbLocation string
		Institutions interface{}
		RequireKeyChaining bool
	}
	Server struct {
		EnableUI bool
		ExternalWebUrl string
		Hostname string
		IssuerHostname string
		IssuerJwks string
		IssuerPort int
		IssuerUrl string
		Modules []string
		SessionSecretFile string
		TLSCACertificateDirectory string
		TLSCACertificateFile string
		TLSCAKey string
		TLSCertificate string
		TLSKey string
		UIActivationCodeFile string
		UIPasswordFile string
		WebHost string
		WebPort int
	}
	StagePlugin struct {
		Hook bool
		MountPrefix string
		OriginPrefix string
		ShadowOriginPrefix string
	}
	TLSSkipVerify bool
	Transport struct {
		DialerKeepAlive time.Duration
		DialerTimeout time.Duration
		ExpectContinueTimeout time.Duration
		IdleConnTimeout time.Duration
		MaxIdleConns int
		ResponseHeaderTimeout time.Duration
		TLSHandshakeTimeout time.Duration
	}
	Xrootd struct {
		Authfile string
		DetailedMonitoringHost string
		LocalMonitoringHost string
		MacaroonsKeyFile string
		ManagerHost string
		Mount string
		Port int
		RobotsTxtFile string
		RunLocation string
		ScitokensConfig string
		Sitename string
		SummaryMonitoringHost string
	}
}


type configWithType struct {
	Server struct {
		Hostname struct { Type string; Value string }
		TLSKey struct { Type string; Value string }
		EnableUI struct { Type string; Value bool }
		UIPasswordFile struct { Type string; Value string }
		TLSCertificate struct { Type string; Value string }
		TLSCACertificateFile struct { Type string; Value string }
		ExternalWebUrl struct { Type string; Value string }
		IssuerUrl struct { Type string; Value string }
		IssuerPort struct { Type string; Value int }
		IssuerJwks struct { Type string; Value string }
		UIActivationCodeFile struct { Type string; Value string }
		TLSCACertificateDirectory struct { Type string; Value string }
		WebHost struct { Type string; Value string }
		IssuerHostname struct { Type string; Value string }
		Modules struct { Type string; Value []string }
		SessionSecretFile struct { Type string; Value string }
		TLSCAKey struct { Type string; Value string }
		WebPort struct { Type string; Value int }
	}
	Issuer struct {
		AuthorizationTemplates struct { Type string; Value interface{} }
		ScitokensServerLocation struct { Type string; Value string }
		QDLLocation struct { Type string; Value string }
		GroupFile struct { Type string; Value string }
		OIDCAuthenticationUserClaim struct { Type string; Value string }
		GroupSource struct { Type string; Value string }
		GroupRequirements struct { Type string; Value []string }
		TomcatLocation struct { Type string; Value string }
		AuthenticationSource struct { Type string; Value string }
		OIDCAuthenticationRequirements struct { Type string; Value interface{} }
	}
	Xrootd struct {
		Port struct { Type string; Value int }
		RobotsTxtFile struct { Type string; Value string }
		MacaroonsKeyFile struct { Type string; Value string }
		Authfile struct { Type string; Value string }
		SummaryMonitoringHost struct { Type string; Value string }
		DetailedMonitoringHost struct { Type string; Value string }
		RunLocation struct { Type string; Value string }
		ScitokensConfig struct { Type string; Value string }
		Mount struct { Type string; Value string }
		ManagerHost struct { Type string; Value string }
		LocalMonitoringHost struct { Type string; Value string }
		Sitename struct { Type string; Value string }
	}
	ConfigDir struct { Type string; Value string }
	IssuerKey struct { Type string; Value string }
	Client struct {
		StoppedTransferTimeout struct { Type string; Value int }
		SlowTransferRampupTime struct { Type string; Value int }
		SlowTransferWindow struct { Type string; Value int }
		DisableHttpProxy struct { Type string; Value bool }
		DisableProxyFallback struct { Type string; Value bool }
		MinimumDownloadSpeed struct { Type string; Value int }
	}
	DisableProxyFallback struct { Type string; Value bool }
	Registry struct {
		DbLocation struct { Type string; Value string }
		RequireKeyChaining struct { Type string; Value bool }
		AdminUsers struct { Type string; Value []string }
		Institutions struct { Type string; Value interface{} }
	}
	Transport struct {
		DialerTimeout struct { Type string; Value time.Duration }
		DialerKeepAlive struct { Type string; Value time.Duration }
		MaxIdleConns struct { Type string; Value int }
		IdleConnTimeout struct { Type string; Value time.Duration }
		TLSHandshakeTimeout struct { Type string; Value time.Duration }
		ExpectContinueTimeout struct { Type string; Value time.Duration }
		ResponseHeaderTimeout struct { Type string; Value time.Duration }
	}
	GeoIPOverrides struct { Type string; Value interface{} }
	MinimumDownloadSpeed struct { Type string; Value int }
	OIDC struct {
		UserInfoEndpoint struct { Type string; Value string }
		AuthorizationEndpoint struct { Type string; Value string }
		ClientIDFile struct { Type string; Value string }
		DeviceAuthEndpoint struct { Type string; Value string }
		TokenEndpoint struct { Type string; Value string }
		ClientRedirectHostname struct { Type string; Value string }
		ClientID struct { Type string; Value string }
		ClientSecretFile struct { Type string; Value string }
		Issuer struct { Type string; Value string }
	}
	Monitoring struct {
		PortHigher struct { Type string; Value int }
		AggregatePrefixes struct { Type string; Value []string }
		TokenExpiresIn struct { Type string; Value time.Duration }
		TokenRefreshInterval struct { Type string; Value time.Duration }
		MetricAuthorization struct { Type string; Value bool }
		DataLocation struct { Type string; Value string }
		PortLower struct { Type string; Value int }
	}
	TLSSkipVerify struct { Type string; Value bool }
	Federation struct {
		RegistryUrl struct { Type string; Value string }
		JwkUrl struct { Type string; Value string }
		DiscoveryUrl struct { Type string; Value string }
		TopologyNamespaceUrl struct { Type string; Value string }
		TopologyReloadInterval struct { Type string; Value time.Duration }
		DirectorUrl struct { Type string; Value string }
		NamespaceUrl struct { Type string; Value string }
	}
	DisableHttpProxy struct { Type string; Value bool }
	Origin struct {
		XRootDPrefix struct { Type string; Value string }
		Mode struct { Type string; Value string }
		S3Region struct { Type string; Value string }
		S3Bucket struct { Type string; Value string }
		EnableIssuer struct { Type string; Value bool }
		ScitokensUsernameClaim struct { Type string; Value string }
		ScitokensNameMapFile struct { Type string; Value string }
		ScitokensDefaultUser struct { Type string; Value string }
		EnableVoms struct { Type string; Value bool }
		S3ServiceName struct { Type string; Value string }
		Url struct { Type string; Value string }
		SelfTest struct { Type string; Value bool }
		ScitokensMapSubject struct { Type string; Value bool }
		EnableDirListing struct { Type string; Value bool }
		S3AccessKeyfile struct { Type string; Value string }
		Multiuser struct { Type string; Value bool }
		EnableCmsd struct { Type string; Value bool }
		ScitokensRestrictedPaths struct { Type string; Value []string }
		S3ServiceUrl struct { Type string; Value string }
		S3SecretKeyfile struct { Type string; Value string }
		ExportVolume struct { Type string; Value string }
		NamespacePrefix struct { Type string; Value string }
		EnableUI struct { Type string; Value bool }
	}
	Plugin struct {
		Token struct { Type string; Value string }
	}
	StagePlugin struct {
		Hook struct { Type string; Value bool }
		MountPrefix struct { Type string; Value string }
		OriginPrefix struct { Type string; Value string }
		ShadowOriginPrefix struct { Type string; Value string }
	}
	Debug struct { Type string; Value bool }
	Logging struct {
		Level struct { Type string; Value string }
		LogLocation struct { Type string; Value string }
	}
	Cache struct {
		DataLocation struct { Type string; Value string }
		ExportLocation struct { Type string; Value string }
		XRootDPrefix struct { Type string; Value string }
		Port struct { Type string; Value int }
		EnableVoms struct { Type string; Value bool }
	}
	Director struct {
		MaxMindKeyFile struct { Type string; Value string }
		GeoIPLocation struct { Type string; Value string }
		DefaultResponse struct { Type string; Value string }
		CacheResponseHostnames struct { Type string; Value []string }
		OriginResponseHostnames struct { Type string; Value []string }
	}
}

