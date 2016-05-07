package common

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

// Config defines the configuration struct for the application
// This will typically be populated by a TOML file
type Config struct {
	sourceFile string
	Core       struct {
		SiteTitle       string
		SiteCompanyName string
		SiteDomainName  string
	}
	Logging struct {
		Enabled    bool
		EnableHTTP bool
		Level      string
		Path       string
	}
	Database struct {
		Type     string
		Address  string
		Username string
		Password string
	}
	Webserver struct {
		Address             string
		HTTPPort            int
		HTTPSPort           int
		TLSCertFile         string
		TLSKeyFile          string
		RedirectHTTPToHTTPS bool
	}
}

// NewEmptyConfig creates an empty configuration. Used for testing.
func NewEmptyConfig() *Config {
	return &Config{}
}

// NewConfig parses a TOML file into a Config struct
func NewConfig(configFile string) (conf *Config, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	if configFile == "" {
		return nil, errors.New("Empty configuration file path")
	}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var con Config
	if err := toml.Unmarshal(buf, &con); err != nil {
		return nil, err
	}
	con.sourceFile = configFile
	return setSensibleDefaults(&con)
}

// Use this function to set any kind of defaults based on a type's zero value.
func setSensibleDefaults(c *Config) (*Config, error) {
	// Anything not set here implies its zero value is the default

	// Core
	c.Core.SiteTitle = setStringOrDefault(c.Core.SiteTitle, "Application Title")

	// Logging
	c.Logging.Level = setStringOrDefault(c.Logging.Level, "notice")
	c.Logging.Path = setStringOrDefault(c.Logging.Path, "logs/application.log")

	// Database
	c.Database.Type = setStringOrDefault(c.Database.Type, "sqlite")
	c.Database.Address = setStringOrDefault(c.Database.Address, "database.sqlite3")

	// Webserver
	c.Webserver.HTTPPort = setIntOrDefault(c.Webserver.HTTPPort, 8080)
	c.Webserver.HTTPSPort = setIntOrDefault(c.Webserver.HTTPSPort, 1443)
	return c, nil
}

// Given string s, if it is empty, return v else return s.
func setStringOrDefault(s, v string) string {
	if s == "" {
		return v
	}
	return s
}

// Given integer s, if it is 0, return v else return s.
func setIntOrDefault(s, v int) int {
	if s == 0 {
		return v
	}
	return s
}
