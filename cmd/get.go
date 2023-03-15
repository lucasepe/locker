package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/locker/internal/kv"
)

const (
	outputTxt = "txt"
	outputEnv = "env"
)

func newCmdGet() *cmdGet {
	return &cmdGet{
		namespace: flags.Namespace{},
		key:       flags.Key{},
		storeRef: flags.Store{
			BaseDir: AppDir(),
		},
		output: flags.Enum{Choices: []string{outputEnv, outputTxt}},
	}
}

type cmdGet struct {
	namespace flags.Namespace
	key       flags.Key
	storeRef  flags.Store
	output    flags.Enum
}

func (*cmdGet) Name() string { return "get" }
func (*cmdGet) Synopsis() string {
	return "Get one or all secrets from a namespace."
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
	fs.Var(&c.key, "k", "Secret key.")
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

	if len(c.key.Bytes()) > 0 {
		return c.extractItem(sto, fs)
	}

	return c.extractAll(sto, fs)
}

func (c *cmdGet) extractAll(sto kv.Store, fs *flag.FlagSet) error {
	all, err := sto.GetAll(c.namespace.String())
	if err != nil {
		return err
	}

	for k, v := range all {
		if c.output.String() == "env" {
			fmt.Fprintf(fs.Output(), "%s=%s\n", strings.ToUpper(k), v)
		} else {
			fmt.Fprintf(fs.Output(), "%s: %s\n", k, v)
		}
	}

	return nil
}

func (c *cmdGet) extractItem(sto kv.Store, fs *flag.FlagSet) error {
	val, err := sto.GetOne(c.namespace.String(), c.key.String())
	if err != nil {
		return err
	}

	if c.output.Value == outputTxt {
		fmt.Fprintf(fs.Output(), "%s", val)
	} else {
		fmt.Fprintf(fs.Output(), "%s=%s", strings.ToUpper(c.key.String()), val)
	}

	return nil
}

func (c *cmdGet) complete(fs *flag.FlagSet) error {
	if len(c.namespace.Bytes()) == 0 {
		return fmt.Errorf("missing namespace")
	}

	if c.output.Value == "" {
		c.output.Set(outputTxt)
	}

	pwd, err := getMasterSecret()
	if err != nil {
		return err
	}
	c.storeRef.MasterSecret = pwd

	return nil
}
