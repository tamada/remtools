package remtools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const VERSION = "5.0.0"

/*
Context represents the values for running remtools.
*/
type Context struct {
	Patterns       []*regexp.Regexp
	trashBoxPath   string
	followSymlinks bool
}

func NewContext() *Context {
	var context = new(Context)
	context.initialize()
	return context
}

func GetVersion(prog string) string {
	return fmt.Sprintf("%s version %s", prog, VERSION)
}

func (context *Context) IsExitOnError() bool {
	return false
}

func (context *Context) IsFollowSymlink() bool {
	return context.followSymlinks
}

func (context *Context) Path() string {
	abs, _ := filepath.Abs(context.trashBoxPath)

	return filepath.Clean(abs)
}

func appendRegexp(context *Context, regexpString string) {
	pattern, err := regexp.Compile(regexpString)
	if err == nil {
		context.Patterns = append(context.Patterns, pattern)
	}
}

func initializeTrashPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".Trash")
}

func (context *Context) initialize() {
	context.Patterns = []*regexp.Regexp{}
	context.trashBoxPath = initializeTrashPath()
	context.followSymlinks = false
	appendRegexp(context, `.*\~`)
	appendRegexp(context, `\.DS_Store`)
}

func AskToUser(fileName, message string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s: %s (y/N)> ", fileName, message)
		input, _ := reader.ReadString('\n')
		lowerInput := strings.TrimSpace(strings.ToLower(input))
		if strings.HasPrefix(lowerInput, "y") {
			return true
		} else if lowerInput == "" || strings.HasPrefix(lowerInput, "n") {
			return false
		}
		fmt.Printf("%s: Invalid input\n", strings.TrimSpace(input))
	}
}

func (context *Context) Copy(path string) bool {
	return true
}

func (context *Context) Move(path string) bool {
	switch runtime.GOOS {
	case "darwin":
		return context.moveTrashOnMac(path)
	}
	return false
}

func EmptyTrash() bool {
	switch runtime.GOOS {
	case "darwin":
		return emptyTrashOnMac()
	}
	return false
}

func emptyTrashOnMac() bool {
	cmd := exec.Command("osascript", "-e", `tell application "Finder" to empty trash`)
	_, err := cmd.Output()
	if err != nil {
		return false
	}
	return true
}

func (context *Context) moveTrashOnMac(path string) bool {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Finder"
	move POSIX file "%s" to trash
end tell`, path))
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
