package remtools

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/go-homedir"
)

/*
Config represents the parameters for running remtools.
*/
type Config struct {
	Patterns       []*regexp.Regexp
	followSymlinks bool
	exitOnError    bool
}

type configFile struct {
	Patterns       []string `json:"patterns"`
	FollowSymlinks bool     `json:"follow-symlink"`
	ExitOnError    bool     `json:"exit-on-error"`
}

func NewConfig() *Config {
	return &Config{Patterns: []*regexp.Regexp{}, followSymlinks: false, exitOnError: false}
}

func (config *Config) IsExitOnError() bool {
	return config.exitOnError
}

func (config *Config) IsFollowSymlink() bool {
	return config.followSymlinks
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (config *Config) Initialize() {
	home, err := homedir.Dir()
	if err != nil {
		if !config.ReadConfig(filepath.Join(home, ".remtools.json")) {
			config.defaultConfig()
		}
	} else {
		config.defaultConfig()
	}
}

func (config *Config) ReadConfig(configPath string) bool {
	if !Exists(configPath) {
		return false
	}
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return false
	}
	configFile := configFile{}
	if err := json.Unmarshal(bytes, &configFile); err != nil {
		return false
	}
	copyItems(configFile, config)
	return true
}

func (config *Config) defaultConfig() {
	config.defaultPatterns()
}

func (config *Config) defaultPatterns() {
	appendRegexp(config, `.*\~$`)
	appendRegexp(config, `.*\.bak$`)
	appendRegexp(config, `^\.DS_Store$`)
}

func copyItems(from configFile, to *Config) {
	to.followSymlinks = from.FollowSymlinks
	to.exitOnError = from.ExitOnError
	for _, pattern := range from.Patterns {
		appendRegexp(to, pattern)
	}
}

func appendRegexp(conf *Config, regexpString string) {
	pattern, err := regexp.Compile(regexpString)
	if err == nil {
		conf.Patterns = append(conf.Patterns, pattern)
	}
}
