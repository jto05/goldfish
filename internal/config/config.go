package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Server ServerConfig `yaml:"server"`
		API    APIConfig    `yaml:"api"`
	}

	ServerConfig struct {
		JavaPath string `yaml:"java_path"`
		JarPath  string `yaml:"jar_path"`
		ModsDir  string `yaml:"mods_dir"`
	}
	APIConfig struct {
		ListenAddr string `yaml:"listen_addr"`
	}
)

// Load takes the filepath of a config.yaml file and decodes it for dashboard
// configuration
func Load(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode file input stream into cfg
	var cfg Config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
