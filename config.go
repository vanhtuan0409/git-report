package gitreport

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	FilterEmail      []string `yaml:"emails"`
	Repos            []string `yaml:"repositories"`
	DefaultTimeRange int      `yaml:"default_time_range"`
}

func (c *Config) ToString() string {
	sb := new(strings.Builder)
	fmt.Fprintf(sb, "Emails:\n")
	for _, email := range c.FilterEmail {
		fmt.Fprintf(sb, "  - %s\n", email)
	}

	fmt.Fprintf(sb, "Repositories:\n")
	for _, repo := range c.Repos {
		fmt.Fprintf(sb, "  - %s\n", repo)
	}

	fmt.Fprintf(sb, "Default time range: %d\n", c.DefaultTimeRange)
	return sb.String()
}

func GetDefaultConfigPath() string {
	path, _ := ResolvePath("~/.greport/config.yml")
	return path
}

func ReadConfigFromFile(filePath string) (*Config, error) {
	configContent, err := ioutil.ReadFile(filePath)

	config := new(Config)
	// cannot read config file
	if err != nil {
		config, err = CreateDefaultConfig(filePath)
		if err != nil {
			return nil, err
		}
		setDefaultConfig(config)
		return config, nil
	}

	err = yaml.Unmarshal(configContent, config)
	if err != nil {
		return nil, err
	}
	setDefaultConfig(config)

	return config, nil
}

func CreateDefaultConfig(filePath string) (*Config, error) {
	directoryPath := filepath.Dir(filePath)
	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	defaultConfig := Config{
		FilterEmail:      []string{},
		Repos:            []string{},
		DefaultTimeRange: 7,
	}
	out, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(out)
	return &defaultConfig, err
}

func setDefaultConfig(config *Config) {
	if config.DefaultTimeRange == 0 {
		config.DefaultTimeRange = 7
	}
	if len(config.Repos) == 0 {
		pwd, err := os.Getwd()
		if err == nil {
			config.Repos = []string{pwd}
		}
	}
}
