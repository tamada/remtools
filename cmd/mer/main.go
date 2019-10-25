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
	var flag = flag.NewFlagSet("mer", flag.ContinueOnError)
	flag.Usage = func() { fmt.Println(getHelpMessage("mer")) }
	flag.BoolVarP(&options.inquiryFlag, "inquiry", "i", false, "inquiry mode")
	flag.BoolVarP(&options.versionFlag, "version", "v", false, "print version")
	flag.BoolVarP(&options.helpFlag, "help", "h", false, "print this message")
	return flag, &options
}

func getHelpMessage(prog string) string {
	return fmt.Sprintf(`%s [OPTIONS]
OPTIONS
    -i, --inquiry    inquiry mode.
    -h, --help       print this message.
    -V, --Version    print version.`, prog)
}

func performImpl(args []string, opts *options) int {
	return 0
}

func perform(args []string, opts *options) int {
	if opts.versionFlag {
		fmt.Printf(remtools.GetVersion("mer"))
	}
	if opts.helpFlag {
		fmt.Printf(getHelpMessage("mer"))
	}
	if opts.versionFlag || opts.helpFlag {
		return 0
	}
	return performImpl(args, opts)
}

func goMain(args []string) int {
	var flagset, opts = buildFlagSet()
	if err := flagset.Parse(args); err != nil {
		fmt.Println(getHelpMessage("mer"))
		return 1
	}
	return perform(flagset.Args(), opts)
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
