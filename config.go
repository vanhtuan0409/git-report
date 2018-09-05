package gitreport

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	FilterEmail []string `yaml:"emails"`
	Repos       []string `yaml:"repositories"`
}

func GetDefaultConfigPath() string {
	path, _ := ResolvePath("~/.greport/config.yml")
	return path
}

func ReadConfigFromFile(filePath string) (*Config, error) {
	configContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, ErrNoFileConfig
	}

	var config Config
	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func CreateDefaultConfig(filePath string) error {
	directoryPath := filepath.Dir(filePath)
	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	defaultConfig := Config{
		FilterEmail: []string{"your_email@domain.com"},
		Repos:       []string{"<path_to_repo>"},
	}
	out, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	_, err = f.Write(out)
	return err
}
