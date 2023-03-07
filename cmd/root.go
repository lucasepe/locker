package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/internal/text"
	"github.com/lucasepe/subcommands"
)

const (
	banner = `┬  ┌─┐┌─┐┬┌─┌─┐┬─┐
│  │ ││  ├┴┐├┤ ├┬┘
┴─┘└─┘└─┘┴ ┴└─┘┴└─`
	summary     = "Store secrets on your local file system."
	maxFileSize = 1024 * 128
)

func Run(ver, bld string) error {
	err := os.MkdirAll(app.Dir(), os.ModePerm)
	if err != nil {
		return err
	}

	cli := subcommands.New(flag.CommandLine, app.Name)
	cli.Banner = fmt.Sprintf("%s\n%s\n", banner, summary)
	if len(app.MasterPassword()) == 0 {
		cli.Banner = fmt.Sprintf("%s\n> %s\n",
			cli.Banner, app.ErrUnsetMasterPassword)
	}
	cli.Register(cli.HelpCommand(), "")

	cli.Register(newCmdGet(), "")
	cli.Register(newCmdPut(), "")
	cli.Register(newCmdList(), "")
	cli.Register(newCmdInfo(ver, bld), "")
	cli.Register(newCmdDelete(), "")
	cli.Register(newCmdImport(), "")

	flag.Parse()

	return cli.Execute()
}

func grabContent(fs *flag.FlagSet) []byte {
	var reader io.Reader

	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return []byte{}
	}

	if (info.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		reader = io.LimitReader(bufio.NewReader(os.Stdin), maxFileSize)
	} else if fs.NArg() > 0 {
		reader = strings.NewReader(fs.Arg(0))
	}

	if reader == nil {
		return []byte{}
	}

	dat, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return []byte{}
	}

	if !text.IsText(dat) {
		fmt.Fprintln(os.Stderr, "warn: content must be text")
		return []byte{}
	}

	return dat
}
