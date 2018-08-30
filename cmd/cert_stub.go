// +build !cert

package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var Cert = cli.Command{
	Name:        "cert",
	Usage:       "Generate self-signed certificate",
	Description: `Please use build tags "cert" to rebuild Gitote in order to have this ability`,
	Action:      runCert,
}

func runCert(ctx *cli.Context) error {
	fmt.Println("Command cert not available, please use build tags 'cert' to rebuild.")
	os.Exit(1)

	return nil
}
