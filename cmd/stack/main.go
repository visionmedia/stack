package main

import "github.com/visionmedia/stack/pkg/logger/interactive"
import "github.com/visionmedia/stack/pkg/logger/plain"
import "github.com/visionmedia/stack/pkg/logger/tty"
import "github.com/visionmedia/stack/pkg/provisioner"
import "github.com/visionmedia/docopt-go"
import "path/filepath"
import "os/signal"
import "os/user"
import "syscall"
import "fmt"
import "os"

const Usage = `
  Usage:
    stack [--list] [--no-color] [--verbose] <file>
    stack -h | --help
    stack --version

  Options:
    -C, --no-color   output with color disabled
    -l, --list       output commit status
    -V, --verbose    output command stdio
    -h, --help       output help information
    -v, --version    output version

`

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args, err := docopt.Parse(Usage, nil, true, provisioner.Version, false)
	check(err)

	u, err := user.Current()
	check(err)

	file := args["<file>"].(string)
	f, err := os.Open(file)
	check(err)

	path := filepath.Join(u.HomeDir, ".provision.log")
	p := provisioner.New(f, path)
	p.DryRun = args["--list"].(bool)
	p.Verbose = args["--verbose"].(bool)

	switch {
	case args["--no-color"].(bool):
		p.Log = plain_logger.New(os.Stdout)
	case p.Verbose:
		p.Log = tty_logger.New(os.Stdout)
	default:
		p.Log = interactive_logger.New(os.Stdout)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	hide()
	go func() {
		<-ch
		show()
		os.Exit(1)
	}()

	p.Run()
	show()
}

func show() {
	fmt.Printf("\033[?25h\n")
}

func hide() {
	fmt.Printf("\033[?25l\n")
}
