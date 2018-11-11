package cmd

import (
	"os"

	"github.com/urfave/cli"
	"github.com/xenolf/lego/log"
)

const filePerm os.FileMode = 0600

func Before(c *cli.Context) error {
	if len(c.GlobalString("path")) == 0 {
		log.Fatal("Could not determine current working directory. Please pass --path.")
	}

	err := createNonExistingFolder(c.GlobalString("path"))
	if err != nil {
		log.Fatalf("Could not check/create path: %v", err)
	}

	if len(c.GlobalString("server")) == 0 {
		log.Fatal("Could not determine current working server. Please pass --server.")
	}

	return nil
}

func createNonExistingFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return nil
}
