package config

type Config struct {
	App                 App
	Database            Database
	Redis               Redis
	Minio               Minio
	Mail                Mail
	RolePathPermissions RolePathPermissions `mapstructure:"role-path-permissions"`
	Jwt                 Jwt
}

type App struct {
	Name    string
	Version string
	Port    string
	SSL     bool
	Cert    string
	Key     string
	CORS    CORS
}

type CORS struct {
	AllowOrigins     []string `mapstructure:"allow-origins"`
	AllowOriginRegex []string `mapstructure:"allow-origin-regex"`
	AllowMethods     []string `mapstructure:"allow-methods"`
	AllowHeaders     []string `mapstructure:"allow-headers"`
	ExposeHeaders    []string `mapstructure:"expose-headers"`
	AllowCredentials bool     `mapstructure:"allow-credentials"`
	MaxAge           int      `mapstructure:"max-age"`
}

type Database struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

type Redis struct {
	Database   int
	Host       string
	Port       string
	Password   string
	TimeToLive int `mapstructure:"time-to-live"`
}

type Minio struct {
	Endpoint  string
	UseSSL    bool   `mapstructure:"useSSL"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Bucket    string
}

type Mail struct {
	Host       string
	Port       int
	User       string
	Password   string
	SSL        bool
	Timeout    int
	SenderName string `mapstructure:"sender-name"`
}

type RolePathPermissions struct {
	Permissions map[string][]string `mapstructure:"permissions"`
}

type Jwt struct {
	Secret     string
	Expiration int64
}
