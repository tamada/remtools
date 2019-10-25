package remtools

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestPeripheral(t *testing.T) {
	if GetVersion("remtools") != fmt.Sprintf("remtools version %s", VERSION) {
		t.Errorf("version string was wrong")
	}
}

func TestNewContext(t *testing.T) {
	context := NewContext()
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".Trash")
	if context.Path() != filepath.Join(home, ".Trash") {
		t.Errorf("trash box path was wrong, wont: %s, got: %s", context.Path(), path)
	}
}

func TestCopy(t *testing.T) {
	context := NewContext()
	path := filepath.Join(context.Path(), "config_test.json")
	if Exists(path) {
		os.Remove(path)
	}
	context.Copy("testdata/config_test.json")
	defer os.Remove(path)
	if !Exists(path) {
		t.Errorf("Copy failed, %s wont exists, got not exists", path)
	}
}

func TestMove(t *testing.T) {
	context := NewContext()
	path := filepath.Join(context.Path(), "config_test2.json")
	if Exists(path) {
		os.Remove(path)
	}
	copy("testdata/config_test.json", "testdata/config_test2.json")
	targetPath, _ := filepath.Abs("testdata/config_test2.json")
	if !context.Move(targetPath) {
		t.Errorf("move failed: %s", targetPath)
	}
	if !Exists(path) {
		t.Errorf("file %s does not exists", path)
	}
}
