package cmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/locker/internal/clipboard"
	"github.com/lucasepe/locker/internal/kv"
)

const (
	fmtTxt = "txt"
	fmtEnv = "env"
)

func newCmdGet() *cmdGet {
	return &cmdGet{
		namespace: flags.Namespace{},
		keys:      flags.StringList{},
		storeRef: flags.Store{
			BaseDir: AppDir(),
		},
		output: flags.Enum{Choices: []string{fmtEnv, fmtTxt}},
		exportFuncMap: map[string]exportFunc{
			fmtEnv: func(w io.Writer, k, v string) {
				fmt.Fprintf(w, "%s=%s", strings.ToUpper(k), v)
			},
			fmtTxt: func(w io.Writer, k, v string) {
				fmt.Fprintf(w, "%s: %s", k, v)
			},
		},
	}
}

type cmdGet struct {
	namespace     flags.Namespace
	keys          flags.StringList
	storeRef      flags.Store
	output        flags.Enum
	exportFuncMap map[string]exportFunc
}

func (*cmdGet) Name() string { return "get" }
func (*cmdGet) Synopsis() string {
	return "Get one, some or all secrets from a namespace."
}

func (*cmdGet) Usage() string {
	return strings.ReplaceAll(`{NAME} get [flags]
  
   Get the secret with key 'user' from the 'google' namespace:
     {NAME} get -n Google -k user

   Get all secrets from the 'google' namespace:
     {NAME} get -n google`, "{NAME}", appLowerName)
}

func (c *cmdGet) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.namespace, "n", "Namespace.")
	fs.Var(&c.storeRef, "s", "Store name.")
	fs.Var(&c.keys, "k", "Secret key.")
	fs.Var(&c.output, "o", fmt.Sprintf("Output format, one of: %s", strings.Join(c.output.Choices, ",")))
}

func (c *cmdGet) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	sto, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer sto.Close()

	keys := c.keys.Values()
	if len(keys) == 1 {
		return c.extractOne(sto, fs)
	}

	return c.extractAll(sto, fs)
}

func (c *cmdGet) extractAll(sto kv.Store, fs *flag.FlagSet) error {
	all, err := sto.GetAll(c.namespace.String(), c.keys.Values()...)
	if err != nil {
		return err
	}

	of := c.output.String()
	i := 0
	for k, v := range all {
		c.exportFuncMap[of](fs.Output(), k, v)
		if i < len(all)-1 {
			fmt.Fprintln(fs.Output())
		}
	}

	return nil
}

func (c *cmdGet) extractOne(sto kv.Store, fs *flag.FlagSet) error {
	key := c.keys.Values()[0]
	val, err := sto.GetOne(c.namespace.String(), key)
	if err != nil {
		return err
	}

	if c.output.Value == fmtTxt {
		fmt.Fprintf(fs.Output(), "%s", val)
	} else {
		c.exportFuncMap[c.output.Value](fs.Output(), key, val)
	}

	clipboard.Write([]byte(val))

	return nil
}

func (c *cmdGet) complete(fs *flag.FlagSet) error {
	if len(c.namespace.Bytes()) == 0 {
		return fmt.Errorf("missing namespace")
	}

	if c.output.Value == "" {
		c.output.Set(fmtTxt)
	}

	pwd, err := getMasterSecret()
	if err != nil {
		return err
	}
	c.storeRef.MasterSecret = pwd

	return nil
}

type exportFunc func(w io.Writer, key, val string)
