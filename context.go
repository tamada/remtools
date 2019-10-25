package remtools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

const VERSION = "5.0.0"

type Context interface {
	Path() string
	Move(fromPath string) bool
	Copy(fromPath string) bool
	InitPath()
	EmptyTrash() bool
}

func NewContext() Context {
	context := createContext()
	context.InitPath()
	return context
}

func createContext() Context {
	switch runtime.GOOS {
	case "darwin":
		return new(DarwinContext)
	default:
		return new(GeneralContext)
	}
}

func GetVersion(prog string) string {
	return fmt.Sprintf("%s version %s", prog, VERSION)
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

func copy(fromPath, toPath string) bool {
	from, err1 := os.Open(fromPath)
	to, err2 := os.OpenFile(toPath, os.O_CREATE, 0644)
	if err1 != nil || err2 != nil {
		return false
	}
	defer from.Close()
	defer to.Close()
	io.Copy(to, from)
	return true
}
