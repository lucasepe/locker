package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/totp"
)

func newCmdTotp() *cmdTotp {
	return &cmdTotp{
		namespace: flags.Namespace{},
		storeRef: flags.Store{
			BaseDir: AppDir(),
		},
	}
}

type cmdTotp struct {
	namespace flags.Namespace
	storeRef  flags.Store
}

func (*cmdTotp) Name() string { return "totp" }
func (*cmdTotp) Synopsis() string {
	return "Generate a time-based OTP from a 'totp' key into a namespace."
}

func (*cmdTotp) Usage() string {
	return strings.ReplaceAll(`{NAME} totp [flags]
  
   Generate a TOTP from a totp url stored into the 'google' namespace:
     {NAME} totp -n google`, "{NAME}", appLowerName)
}

func (c *cmdTotp) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.namespace, "n", "Namespace.")
	fs.Var(&c.storeRef, "s", "Store name.")
}

func (c *cmdTotp) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	sto, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer sto.Close()

	uri, err := sto.GetOne(c.namespace.String(), "totp")
	if err != nil {
		return err
	}
	if len(uri) == 0 {
		return fmt.Errorf("totp url not found in namespace: %s", c.namespace.String())
	}

	opts, err := totp.ParseURI(uri)
	if err != nil {
		return err
	}

	code, err := totp.New(opts)
	if err != nil {
		return err
	}

	fmt.Fprint(fs.Output(), code)

	return err
}

func (c *cmdTotp) complete(fs *flag.FlagSet) error {
	if len(c.namespace.Bytes()) == 0 {
		return fmt.Errorf("missing namespace")
	}

	pwd, err := getMasterSecret()
	if err != nil {
		return err
	}
	c.storeRef.MasterSecret = pwd

	return nil
}
