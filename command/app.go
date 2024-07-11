package command

import "github.com/urfave/cli/v3"

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "certsuck",
		Commands: []*cli.Command{
			ChainCommand(),
		},
	}
}

func ChainCommand() *cli.Command {
	return &cli.Command{
		Name: "chain",
		Description: "Show the CA chain that validates the host",
		Usage: "Show ca chain of the host",
		// Action: ,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "output",
			},
		},
	}
}