package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/strcase"

	"gopkg.in/yaml.v3"
)

func newCmdImport() *cmdImport {
	return &cmdImport{
		file:     flags.FileFlag{},
		storeRef: flags.StoreFlag{},
	}
}

type cmdImport struct {
	file     flags.FileFlag
	storeRef flags.StoreFlag
}

func (*cmdImport) Name() string { return "import" }
func (*cmdImport) Synopsis() string {
	return "Import secrets."
}

func (*cmdImport) Usage() string {
	return strings.ReplaceAll("{NAME} import [flags]", "{NAME}", app.Name)
}

func (c *cmdImport) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.storeRef, "s", "Store name.")
	fs.Var(&c.file, "f", "File to import.")
}

func (c *cmdImport) Execute(fs *flag.FlagSet) error {
	if err := c.complete(fs); err != nil {
		return err
	}

	rdr, err := os.OpenFile(c.file.String(), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer rdr.Close()

	db, err := c.storeRef.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	count := 0
	decoder := yaml.NewDecoder(rdr)
	for {
		var d SecretList
		if err := decoder.Decode(&d); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("document decode failed: %w", err)
		}

		for _, el := range d.Secrets {
			err := db.PutOne(strcase.Kebab(d.Namespace), strcase.Snake(el.Key), el.Value)
			if err != nil {
				return err
			}
		}
		count = count + 1
	}

	if count > 0 {
		fmt.Fprintf(fs.Output(), "successfully imported %d documents\n", count)
	}

	return nil
}

func (c *cmdImport) complete(fs *flag.FlagSet) error {
	if len(c.file.String()) == 0 {
		return fmt.Errorf("file to import not specified")
	}

	return nil
}

// An Secret holds a label/value pair.
type Secret struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type SecretList struct {
	Namespace string   `yaml:"namespace"`
	Secrets   []Secret `yaml:"secrets"`
}
