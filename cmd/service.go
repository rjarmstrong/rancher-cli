package cmd

import (
	"time"

	"github.com/nowait-tools/rancher-cli/rancher"
	"github.com/nowait-tools/rancher-cli/rancher/config"
	"github.com/urfave/cli"
)

func ServiceCommand() cli.Command {
	return cli.Command{
		Name:  "service",
		Usage: "Operations on services",
		Subcommands: []cli.Command{
			{
				Name:  "upgrade",
				Usage: "Upgrade a service",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "service",
					},
					cli.StringFlag{
						Name: "service-like",
					},
					cli.StringFlag{
						Name:  "env-file",
						Usage: "File containing environment variables that will be used for validating that the Rancher service has all variables defined",
					},
					cli.StringSliceFlag{
						Name:  "env",
						Usage: "Environment variables to add when upgrading the service",
					},
					cli.StringFlag{
						Name: "runtime-tag",
					},
					cli.StringFlag{
						Name: "code-tag",
					},
					cli.Int64Flag{
						Name:  "interval",
						Usage: "Interval between starting new containers and stopping old ones",
					},
					cli.BoolFlag{
						Name:  "wait",
						Usage: "Wait for the upgrade to fully complete",
					},
				},
				Action: UpgradeAction,
			},
			{
				Name:  "upgrade-finish",
				Usage: "",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "service",
					},
				},
				Action: func(c *cli.Context) error {
					client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret, "")

					if err != nil {
						return err
					}
					_, err = client.FinishServiceUpgrade(c.String("service"))
					return err
				},
			},
		},
	}
}

func UpgradeAction(c *cli.Context) error {

	envFile := c.String("env-file")
	env := c.StringSlice("env")

	if err := config.ValidateEnvFlag(env); err != nil {
		return err
	}

	client, err := rancher.NewClient(cattleUrl, cattleAccessKey, cattleSecret, envFile)
	if err != nil {
		return err
	}
	interval := time.Duration(0)
	if interval = time.Duration(c.Int64("interval") * int64(time.Second)); interval == 0 {
		interval = defaultUpgradeInterval
	}

	opts := config.UpgradeOpts{
		Envs:        env,
		Interval:    interval,
		ServiceLike: c.String("service-like"),
		Service:     c.String("service"),
		CodeTag:     c.String("code-tag"),
		RuntimeTag:  c.String("runtime-tag"),
		Wait:        c.Bool("wait"),
	}
	if name := opts.ServiceLike; name != "" {
		return client.UpgradeServiceWithNameLike(opts)
	}

	_, err = client.UpgradeService(opts)

	return err
}
