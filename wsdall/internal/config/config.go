package config

import (
	"io/ioutil"
	"path/filepath"
	"sync/atomic"

	"github.com/wsdall/internal/logger"
	"gopkg.in/yaml.v3"
)

type WsProxyConfig struct {
	ListenAddr string   `yaml:"listen"`
	Router     string   `yaml:"router"`
	Balance    string   `yaml:"balance"`
	Server     []string `yaml:"server"`
	// hashTable  *CHashTable
}

type Config struct {
	WsProxy []*WsProxyConfig `yaml:"ws_proxy"`
}

const (
	// YAML .yaml
	YAML = ".yaml"
	// YML .yml
	YML = ".yml"
)

var (
	gConfig    atomic.Value
	configFile string
)

func (w *WsProxyConfig) init() {
	if w.Balance == "hash" {
		// w.hashTable = newCHashTable()
		// w.hashTable.init(c.Server)
	}
}

func newConfig() *Config {
	c := &Config{}
	c.WsProxy = make([]*WsProxyConfig, 0)
	return c
}

func (c *Config) loadFromFile() error {
	body, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, c)
	if err != nil {
		return err
	}

	for _, c := range c.WsProxy {
		c.init()
	}

	return nil
}

func getConfig() *Config {
	c := gConfig.Load()
	if c == nil {
		return nil
	}
	return c.(*Config)
}

func CheckYamlFormat(path string) bool {
	ext := filepath.Ext(path)
	if ext == YAML || ext == YML {
		return true
	}
	return false
}

func LoadConfig() {
	configFile = "../internal/config/config.yaml"
	configFile, _ = filepath.Abs(configFile)
	if configFile == "" && !CheckYamlFormat(configFile) {
		logger.Errorf("config file %s format error", configFile)
		return
	}
	logger.Infof("load config file from %s\n", configFile)
	c := newConfig()
	err := c.loadFromFile()
	if err != nil {
		logger.Errorf("load config from file failed. err(%v)\n", err)
		return
	}
	gConfig.Store(c)
	for i, data := range c.WsProxy {
		logger.Infof("config: %d : %+v", i, data)
	}
}
