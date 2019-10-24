package main

import (
	"fmt"
	"os"

	flag "github.com/ogier/pflag"
	"github.com/tamada/remtools"
)

type options struct {
	inquiryFlag bool
	versionFlag bool
	helpFlag    bool
}

func buildFlagSet() (*flag.FlagSet, *options) {
	var options = options{}
	var flag = flag.NewFlagSet("remrem", flag.ContinueOnError)
	flag.Usage = func() { fmt.Println(getHelpMessage()) }
	flag.BoolVarP(&options.inquiryFlag, "inquiry", "i", false, "inquiry mode")
	flag.BoolVarP(&options.versionFlag, "version", "v", false, "print version")
	flag.BoolVarP(&options.helpFlag, "help", "h", false, "print this message")
	return flag, &options
}

func getHelpMessage() string {
	return `remrem [OPTIONS]
OPTIONS
    -i, --inquiry    inquiry mode.

    -v, --version    print version.
    -h, --help       print this message.`
}

func getVersion(prog string) string {
	return fmt.Sprintf("%s version %s", prog, remtools.VERSION)
}

func performImpl(args []string, opts *options) int {
	if opts.inquiryFlag && !remtools.AskToUser("trash", "empty trash?") {
		return 0
	}
	remtools.EmptyTrash()
	return 0
}

func perform(args []string, opts *options) int {
	if opts.versionFlag {
		fmt.Println(remtools.GetVersion("remrem"))
	}
	if opts.helpFlag {
		fmt.Println(getHelpMessage())
	}
	if opts.versionFlag || opts.helpFlag {
		return 0
	}
	return performImpl(args, opts)
}

func goMain(args []string) int {
	var flagset, opts = buildFlagSet()
	if err := flagset.Parse(args); err != nil {
		fmt.Println(getHelpMessage())
		return 1
	}
	return perform(flagset.Args(), opts)
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
