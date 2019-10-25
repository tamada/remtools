package remtools

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type DarwinContext struct {
	trashBoxPath string
}

func (context *DarwinContext) Path() string {
	return context.trashBoxPath
}

func (context *DarwinContext) InitPath() {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".Trash")
	context.trashBoxPath = path
}

func (context *DarwinContext) Move(path string) bool {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Finder"
	move POSIX file "%s" to trash
end tell`, path))
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
		return false
	}
	return true
}

func (context *DarwinContext) Copy(path string) bool {
	toPath := filepath.Join(context.Path(), filepath.Base(path))
	return copy(path, toPath)
}

func (context *DarwinContext) EmptyTrash() bool {
	cmd := exec.Command("osascript", "-e", `tell application "Finder" to empty trash`)
	_, err := cmd.Output()
	if err != nil {
		return false
	}
	return true
}
