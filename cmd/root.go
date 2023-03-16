package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/locker/internal/text"
	"github.com/lucasepe/subcommands"
	"github.com/lucasepe/xdg"
)

const (
	EnvSecret = "LOCKER_SECRET"

	banner = `┬  ┌─┐┌─┐┬┌─┌─┐┬─┐
│  │ ││  ├┴┐├┤ ├┬┘
┴─┘└─┘└─┘┴ ┴└─┘┴└─`
	summary      = "Store secrets on your local file system."
	maxFileSize  = 1024 * 128
	appName      = "Locker"
	appLowerName = "locker"
)

var (
	ErrUnsetMasterSecret = fmt.Errorf(
		"specify a master secret setting the env var: %s", EnvSecret)
)

func Run(ver, bld string) error {
	err := os.MkdirAll(AppDir(), os.ModePerm)
	if err != nil {
		return err
	}

	cli := subcommands.New(flag.CommandLine, appLowerName)
	cli.Banner = fmt.Sprintf("%s\n%s\n", banner, summary)
	if _, err := getMasterSecret(); err != nil {
		cli.Banner = fmt.Sprintf("%s\n> %s\n", cli.Banner, err)
	}
	cli.Register(cli.HelpCommand(), "")

	cli.Register(newCmdGet(), "")
	cli.Register(newCmdPut(), "")
	cli.Register(newCmdList(), "")
	cli.Register(newCmdInfo(ver, bld), "")
	cli.Register(newCmdDelete(), "")
	cli.Register(newCmdImport(), "")
	cli.Register(newCmdTotp(), "")

	flag.Parse()

	return cli.Execute()
}

func AppDir() string {
	return filepath.Join(xdg.ConfigDir(), appName)
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

func getMasterSecret() (string, error) {
	mp := os.Getenv(EnvSecret)
	if len(mp) == 0 {
		return "", ErrUnsetMasterSecret
	}

	return mp, nil
}
