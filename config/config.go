package config

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/BurntSushi/toml"
)

type AuthServerCnf struct {
	URL            string `toml:"url"`
	TimeoutSeconds int    `toml:"timeout_seconds"`
	SkipTLSVerify  bool   `toml:"skip_tls_verify"`
	VerifyTokenUrl string `toml:"verify_token_url"`
}

type Config struct {
	Server struct {
		Host     string
		Port     int
		Env      string
		CertFile string `toml:"cert_file"`
		KeyFile  string `toml:"key_file"`
	}

	Logging struct {
		LogDir   string `toml:"log_dir"`
		LogFile  string `toml:"log_file"`
		LogLevel string `toml:"log_level"`
	}

	Database struct {
		Username string `toml:"username"`
		Password string `toml:"password"`
		Endpoint string `toml:"endpoint"`
		Params   string `toml:"params"`
		DBName   string `toml:"db_name"`
	}

	AuthServer    AuthServerCnf        `toml:"auth-server"`
	WorkloadCollections map[string]string `toml:"workload-collections"`
}

var (
	cfg  *Config
	once sync.Once
)

func loadConfig(path string) {
	var localCfg Config
	if _, err := toml.DecodeFile(path, &localCfg); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	cfg = &localCfg
}

// once.Do is like a singleton, it will cache the result and will return it
// used so that config can be used anywhere without parsing again & again
func InitConfig(path string) {
	once.Do(func() {
		loadConfig(path)
	})
}

func GetConfig() *Config {
	if cfg == nil {
		panic("config not initialized. Call InitConfig(path) before using GetConfig()")
	}
	return cfg
}

func (c *Config) BuildMongoURI() string {
	username := url.QueryEscape(c.Database.Username)
	password := url.QueryEscape(c.Database.Password)
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/?%s",
		username,
		password,
		c.Database.Endpoint,
		c.Database.Params,
	)
}

func (c *Config) GetCollection(workload string) (string, error) {
	if name, ok := c.WorkloadCollections[workload]; ok {
		return name, nil
	}
	return "", fmt.Errorf("no collection found for workload: %s", workload)
}
