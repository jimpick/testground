package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/ipfs/testground/pkg/daemon/client"
	"github.com/ipfs/testground/pkg/server"
	"github.com/urfave/cli"
)

// DescribeCommand is the specification of the `list` command.
var DescribeCommand = cli.Command{
	Name:      "describe",
	Usage:     "describes a test plan or test case",
	ArgsUsage: "[term], where " + server.TermExplanation,
	Action:    describeCommand,
}

func describeCommand(c *cli.Context) error {
	if c.NArg() == 0 {
		_ = cli.ShowSubcommandHelp(c)
		return errors.New("missing term to describe; " + server.TermExplanation)
	}

	term := c.Args().First()

	api, cancel, err := setupClient()
	if err != nil {
		return err
	}
	defer cancel()

	req := &client.DescribeRequest{
		Term: term,
	}
	resp, err := api.Describe(context.Background(), req)
	if err != nil {
		return fmt.Errorf("fatal error from daemon: %s", err)
	}
	defer resp.Close()

	return client.ParseDescribeResponse(resp)
}
