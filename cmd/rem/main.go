package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/ogier/pflag"
	"github.com/tamada/remtools"
)

func help(prog string) string {
	return fmt.Sprintf(`%s [OPTIONS] [DIR...]
OPTIONS
    -a, --all          includes hidden file.
    -d, --dry-run      dry run mode.
    -i, --inquiry      inquiry mode.
    -r, --recursive    recursive mode.
    -v, --verbose      verbose mode.

    -h, --help         print this message and exit.
    -V, --version      print version, and exit.`, prog)
}

type options struct {
	inquiryFlag   bool
	allFlag       bool
	dryRunFlag    bool
	recursiveFlag bool
	verboseFlag   bool
	versionFlag   bool
	helpFlag      bool
	args          []string
}

func buildFlagSet(args []string) (*flag.FlagSet, *options) {
	var options = options{}
	var flag = flag.NewFlagSet("rem", flag.ContinueOnError)
	flag.Usage = func() { fmt.Println(help(args[0])) }
	flag.BoolVarP(&options.allFlag, "all", "a", false, "includes hidden file.")
	flag.BoolVarP(&options.dryRunFlag, "dry-run", "d", false, "dry run mode.")
	flag.BoolVarP(&options.inquiryFlag, "inquiry", "i", false, "inquiry mode.")
	flag.BoolVarP(&options.recursiveFlag, "recursive", "r", false, "recursive mode.")
	flag.BoolVarP(&options.verboseFlag, "verbose", "v", false, "verbose mode.")
	flag.BoolVarP(&options.versionFlag, "version", "V", false, "print version.")
	flag.BoolVarP(&options.helpFlag, "help", "h", false, "print this message.")
	return flag, &options
}

func isHelpOrVersion(opts *options) bool {
	return opts.helpFlag || opts.versionFlag
}

func isSymlinkAndFollowIt(mode os.FileMode, config *remtools.Config) bool {
	if mode&os.ModeSymlink == os.ModeSymlink {
		return config.IsFollowSymlink()
	}
	return false
}

func forceVerbose(event, fileName string, opts *options) {
	fmt.Printf("%-8s    %s\n", event, fileName)
}

func verbose(event, fileName string, opts *options) {
	if opts.verboseFlag {
		forceVerbose(event, fileName, opts)
	}
}

func moveToTrash(name string, opts *options, context remtools.Context) {
	if opts.dryRunFlag {
		verbose("dry run", name, opts)
		return
	}
	if opts.inquiryFlag && !remtools.AskToUser(name, "move to trashbox?") {
		forceVerbose("decline", name, opts)
		return
	}
	if context.Move(name) {
		if opts.verboseFlag {
			verbose("move", name, opts)
		}
	} else {
		forceVerbose("movefail", name, opts)
	}
}

func isRemTarget(file os.FileInfo, opts *options, config *remtools.Config) bool {
	var name = file.Name()
	for _, pattern := range config.Patterns {
		result := pattern.MatchString(name)
		// fmt.Printf("%-7s    %s (%v) %v\n", "regexp", name, pattern, result)
		if result {
			return true
		}
	}
	return false
}

func remEachEntry(dir string, file os.FileInfo, opts *options, context remtools.Context, config *remtools.Config) error {
	fileName := file.Name()
	targetPath := filepath.Join(dir, fileName)
	if opts.recursiveFlag {
		if file.IsDir() {
			return performEachDir(targetPath, opts, context, config)
		} else if isSymlinkAndFollowIt(file.Mode(), config) {
			return performEachDir(targetPath, opts, context, config)
		}
	}
	verbose("check", targetPath, opts)
	if isRemTarget(file, opts, config) {
		moveToTrash(targetPath, opts, context)
	}
	return nil
}

func isTargetFile(file string, opts *options) bool {
	if file == "." || file == ".." {
		return false
	}
	if strings.HasPrefix(file, ".") {
		return opts.allFlag
	}
	return true
}

func performEachDir(arg string, opts *options, context remtools.Context, config *remtools.Config) error {
	verbose("readdir", arg, opts)
	files, err := ioutil.ReadDir(arg)
	if err != nil {
		return err
	}
	for _, file := range files {
		if isTargetFile(file.Name(), opts) {
			remEachEntry(arg, file, opts, context, config)
		}
	}

	return nil
}

var errlist = []error{}

func perform(opts *options) int {
	context := remtools.NewContext()
	config := new(remtools.Config)
	config.Initialize()
	for _, arg := range opts.args {
		err := performEachDir(arg, opts, context, config)
		if err != nil {
			errlist = append(errlist, err)
			if config.IsExitOnError() {
				return 1
			}
		}
	}
	return 0
}

func printHelpAndOrVersion(args []string, opts *options) int {
	if opts.versionFlag {
		fmt.Println(remtools.GetVersion(args[0]))
	}
	if opts.helpFlag {
		fmt.Println(help(args[0]))
	}
	return 0
}

func parse(args []string) (*options, error) {
	flags, opts := buildFlagSet(args)

	if err := flags.Parse(args); err != nil {
		fmt.Println(help(args[0]))
		return nil, err
	}
	opts.args = flags.Args()[1:]
	if len(opts.args) == 0 {
		opts.args = []string{"."}
	}
	return opts, nil
}

func goMain(args []string) int {
	opts, err := parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if isHelpOrVersion(opts) {
		return printHelpAndOrVersion(args, opts)
	}
	return perform(opts)
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
