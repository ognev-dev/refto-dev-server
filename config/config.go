package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DevEnv     = "dev"
	TestEnv    = "test"
	ReleaseEnv = "release"
)

type Config struct {
	AppEnv string `yaml:"app_env"`
	DB     struct {
		Addr       string
		User       string
		Password   string
		Name       string
		LogQueries bool `yaml:"log_queries"`
	}

	Server struct {
		Port        string
		Host        string
		ApiBasePath string `yaml:"api_base_path"`
	}

	Dir struct {
		Data   string
		Logs   string
	}
}

var conf *Config

// Get returns config from .config.yaml
func Get() *Config {
	if conf != nil {
		return conf
	}

	name := ".config.yaml"
	yamlConf, err := ioutil.ReadFile(name)
	if err != nil {
		msg := ".config.yaml missing from working directory"
		wd, err := os.Getwd()
		if err == nil {
			msg += fmt.Sprintf(" (%s)", wd)
		}
		panic(msg)
	}
	err = yaml.Unmarshal(yamlConf, &conf)
	if err != nil {
		panic(err)
	}

	return conf
}

func Reload() *Config {
	conf = nil
	return Get()
}

func IsDevEnv() bool {
	return conf.AppEnv == DevEnv
}

func IsTestEnv() bool {
	return conf.AppEnv == TestEnv
}

func IsReleaseEnv() bool {
	return conf.AppEnv == ReleaseEnv
}
