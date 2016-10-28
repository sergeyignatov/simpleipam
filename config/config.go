package config

import (
	"fmt"
	"github.com/sergeyignatov/simpleipam/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Subnets map[string]Subnet `yaml:"subnets"`
	DataDir string            `yaml:"datadir"`
}

type Subnet struct {
	Start   string `yaml:"start"`
	End     string `yaml:"end"`
	Gateway string `yaml:"gateway"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func SaveClient(cl *common.Client, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create config file: %v", err)
	}
	data, err := yaml.Marshal(cl)
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write configuration: %v", err)
	}

	f.Close()
	return nil
}
func DeleteClient(path string) error {
	err := os.Remove(path)
	return err
}
func LoadClient(path string) (*common.Client, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c common.Client
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
