package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lucasepe/locker/cmd/app"
	"github.com/lucasepe/locker/cmd/flags"
	"github.com/lucasepe/locker/internal/store"
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
	return "Import items."
}

func (*cmdImport) Usage() string {
	return strings.ReplaceAll("{NAME} import [flags]", "{NAME}", app.Name)
}

func (c *cmdImport) SetFlags(fs *flag.FlagSet) {
	fs.Var(&c.storeRef, "n", "Locker name.")
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
		var d Doc
		if err := decoder.Decode(&d); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("document decode failed: %w", err)
		}

		bkt, err := store.NewBucket(db, []byte(d.Box))
		if err != nil {
			return err
		}

		items := make([]struct{ Key, Value []byte }, len(d.Items))
		for i, el := range d.Items {
			val, err := app.Encrypt([]byte(el.Value))
			if err != nil {
				return err
			}
			items[i].Key = []byte(el.Label)
			items[i].Value = val
		}
		if err := bkt.Insert(items); err != nil {
			return err
		}

		count = count + 1
	}

	if count > 0 {
		fmt.Printf("successfully imported %d documents\n", count)
	}

	return nil
}

func (c *cmdImport) complete(fs *flag.FlagSet) error {
	if len(app.MasterPassword()) == 0 {
		return app.ErrUnsetMasterPassword
	}

	if len(c.file.String()) == 0 {
		return fmt.Errorf("file to import not specified")
	}

	return nil
}

// An Item holds a key/value pair.
type Item struct {
	Label string `yaml:"label"`
	Value string `yaml:"value"`
}

type Doc struct {
	Box   string `yaml:"box"`
	Items []Item `yaml:"items"`
}
