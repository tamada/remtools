package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	flag "github.com/ogier/pflag"
	"github.com/tamada/remtools"
)

func help(prog string) string {
	return fmt.Sprintf(`%s [OPTIONS]
OPTIONS
    -a, --all            print hidden files.
    -l, --long-format    print long format.

    -h, --help           print this message.
    -V, --version        print version.`, prog)
}

type options struct {
	allFlag     bool
	longFormat  bool
	helpFlag    bool
	versionFlag bool
	args        []string
}

func buildFlagSet(args []string) (*flag.FlagSet, *options) {
	var options = options{}
	var flag = flag.NewFlagSet("rem", flag.ContinueOnError)
	flag.Usage = func() { fmt.Println(help(args[0])) }
	flag.BoolVarP(&options.allFlag, "all", "a", false, "print hidden file.")
	flag.BoolVarP(&options.longFormat, "long-format", "l", false, "print long format.")
	flag.BoolVarP(&options.versionFlag, "version", "V", false, "print version.")
	flag.BoolVarP(&options.helpFlag, "help", "h", false, "print this message.")
	return flag, &options
}

func parse(args []string) (*options, error) {
	flags, opts := buildFlagSet(args)
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	opts.args = flags.Args()[1:]
	return opts, nil
}

func isHelpOrVersionFlag(opts *options) bool {
	return opts.helpFlag || opts.versionFlag
}

func printVersionAndOrHelp(args []string, opts *options) int {
	if opts.versionFlag {
		fmt.Println(remtools.GetVersion(args[0]))
	}
	if opts.helpFlag {
		fmt.Println(help(args[0]))
	}
	return 0
}

func perform(opts *options) int {
	context := remtools.NewContext()
	options := createOptions(opts)
	cmd := exec.Command("ls", context.Path())
	if options != "" {
		cmd = exec.Command("ls", options, context.Path())
	}
	// fmt.Printf("lsrem: ls %s %s\n", options, context.Path())
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("lsrem: %s\n", err.Error())
		return 1
	}
	fmt.Println(strings.TrimSpace(string(output)))
	return 0
}

func createOptions(opts *options) string {
	options := ""
	if opts.allFlag {
		options = options + "A"
	}
	if opts.longFormat {
		options = options + "l"
	}
	if len(options) > 0 {
		options = "-" + options
	}
	return options
}

func goMain(args []string) int {
	opts, err := parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if isHelpOrVersionFlag(opts) {
		return printVersionAndOrHelp(args, opts)
	}
	return perform(opts)
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
