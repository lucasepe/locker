package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func newCmdInfo(ver, bld string) *cmdInfo {
	return &cmdInfo{
		appVersion: ver,
		appBuild:   bld,
	}
}

type cmdInfo struct {
	appVersion string
	appBuild   string
}

func (*cmdInfo) Name() string { return "info" }

func (*cmdInfo) Synopsis() string {
	return "Print build information and list all existing lockers."
}

func (*cmdInfo) Usage() string {
	return fmt.Sprintf("%s info", appLowerName)
}

func (c *cmdInfo) SetFlags(fs *flag.FlagSet) {}

func (p *cmdInfo) Execute(fs *flag.FlagSet) error {
	fmt.Fprintf(fs.Output(),
		"%s %s (build: %s) <https://github.com/lucasepe/locker>\n", appName, p.appVersion, p.appBuild)
	archives, err := p.listStores()
	if err != nil {
		return err
	}
	if len(archives) == 0 {
		return nil
	}

	fmt.Fprintf(fs.Output(), "\nExisting lockers:\n\n")

	for _, v := range archives {
		fmt.Fprintf(fs.Output(), " - %s\n", v)
	}

	return nil
}

func (c *cmdInfo) listStores() (map[string]string, error) {
	dir := AppDir()
	fp, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	res := map[string]string{}

	list, err := fp.Readdirnames(0)
	if err != nil {
		return res, err
	}

	for _, name := range list {
		if !strings.HasSuffix(name, ".db") {
			continue
		}

		key := strings.TrimSuffix(name, filepath.Ext(name))
		res[key] = filepath.Join(dir, name)
	}

	return res, nil
}
