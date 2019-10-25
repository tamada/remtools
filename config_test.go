package remtools

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config.IsExitOnError() {
		t.Errorf("default value of exitOnError is false")
	}
	if config.IsFollowSymlink() {
		t.Errorf("default value of followSymlinks is false")
	}
	if len(config.Patterns) != 0 {
		t.Errorf("invalid data was added to patterns")
	}
}

func TestInitialize(t *testing.T) {
	config := NewConfig()
	config.Initialize()
	if config.IsExitOnError() {
		t.Errorf("default value of exitOnError is false")
	}
	if config.IsFollowSymlink() {
		t.Errorf("default value of followSymlinks is false")
	}
	if len(config.Patterns) != 3 {
		t.Errorf("invalid data was added to patterns")
	}
}

func TestReadConfig(t *testing.T) {
	config := NewConfig()
	if !config.ReadConfig("testdata/config_test.json") {
		t.Errorf("read failed.")
	}
	if !config.IsExitOnError() {
		t.Errorf("default value of exitOnError is true")
	}
	if !config.IsFollowSymlink() {
		t.Errorf("default value of followSymlinks is true")
	}
	if len(config.Patterns) != 2 {
		t.Errorf("pattern building failed.")
	}

}

func TestReadConfigFailed(t *testing.T) {
	config := NewConfig()
	if config.ReadConfig("testdata/notExistFile.json") {
		t.Errorf("successfully read not exist file !?")
	}
	if !config.ReadConfig("testdata/unrelated.json") {
		t.Errorf("read failed.")
	}
}
