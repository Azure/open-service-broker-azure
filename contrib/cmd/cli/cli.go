package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "broker-cli"
	app.Usage = "demo the azure service broker with ease"
	app.UsageText = "broker-cli [global options] <command> [command options] " +
		"[arguments...]"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagsHost,
			Usage: "specify the broker's host",
			Value: "localhost",
		},
		cli.IntFlag{
			Name:  flagsPort,
			Usage: "specify the broker's port",
			Value: 8080,
		},
		cli.StringFlag{
			Name:  flagsUsername,
			Usage: "specify a username for authenticating to the broker",
		},
		cli.StringFlag{
			Name:  flagsPassword,
			Usage: "specify a password for authenticating to the broker",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "catalog",
			Usage:  "list available services and service plans",
			Action: catalog,
		},
		{
			Name:  "provision",
			Usage: "provision a new service instance",
			UsageText: "broker-cli [global options] provision --service-id " +
				"<service id> --plan-id <plan id> [other command options]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagsServiceID,
					Usage: "specify the `<service id>`; required",
				},
				cli.StringFlag{
					Name:  flagsPlanID,
					Usage: "specify the `<plan id>`; required",
				},
				cli.StringSliceFlag{
					Name:  flagTag,
					Usage: "specify a tag as a k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsParameter,
					Usage: "specify a service-specific string provisioning parameter " +
						"as a k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsIntParameter,
					Usage: "specify a service-specific int provisioning parameter as a " +
						"k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsFloatParameter,
					Usage: "specify a service-specific floating point provisioning " +
						"parameter as a k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsBoolParameter,
					Usage: "specify a service-specific bool provisioning parameter as " +
						"a k=v pair",
				},
				cli.BoolFlag{
					Name: flagPoll,
					Usage: "poll the instance for status until provisioning succeeds " +
						"or fails",
				},
			},
			Action: provision,
		},
		{
			Name:  "update",
			Usage: "update a existing service instance",
			UsageText: "broker-cli [global options] update --instance-id " +
				"<instance id> --service-id <service id> " +
				"[--plan-id <plan id>] [other command options]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagsInstanceID,
					Usage: "specify the `<instance id>`; required",
				},
				cli.StringFlag{
					Name:  flagsServiceID,
					Usage: "specify the `<service id>`; required",
				},
				cli.StringFlag{
					Name:  flagsPlanID,
					Usage: "specify the `<plan id>`; optional",
				},
				cli.StringSliceFlag{
					Name: flagsParameter,
					Usage: "specify a service-specific string update parameter as a " +
						"k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsIntParameter,
					Usage: "specify a service-specific int update parameter as a k=v " +
						"pair",
				},
				cli.StringSliceFlag{
					Name: flagsFloatParameter,
					Usage: "specify a service-specific floating point update parameter " +
						" as a k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsBoolParameter,
					Usage: "specify a service-specific bool update parameter as a k=v " +
						"pair",
				},
				cli.BoolFlag{
					Name: flagPoll,
					Usage: "poll the instance for status until updating succeeds " +
						"or fails",
				},
			},
			Action: update,
		},
		{
			Name:  "poll",
			Usage: "poll instance status",
			UsageText: "broker-cli [global options] poll --instance-id " +
				"<instance id> --operation <provisioning|deprovisioning> " +
				"[other command options]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagsInstanceID,
					Usage: "specify the `<instance id>`; required",
				},
				cli.StringFlag{
					Name:  flagsOperation,
					Usage: "specify the `<operation>`; required",
				},
			},
			Action: poll,
		},
		{
			Name:  "bind",
			Usage: "bind to a service instance",
			UsageText: "broker-cli [global options] bind --instance-id " +
				"<instance id> [other command options]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagsInstanceID,
					Usage: "specify the `<instance id>`; required",
				},
				cli.StringSliceFlag{
					Name: flagsParameter,
					Usage: "specify a service-specific string binding parameter as a " +
						"k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsIntParameter,
					Usage: "specify a service-specific int binding parameter as a k=v " +
						"pair",
				},
				cli.StringSliceFlag{
					Name: flagsFloatParameter,
					Usage: "specify a service-specific floating point binding " +
						"parameter as a k=v pair",
				},
				cli.StringSliceFlag{
					Name: flagsBoolParameter,
					Usage: "specify a service-specific bool binding parameter as a k=v " +
						"pair",
				},
			},
			Action: bind,
		},
		{
			Name:  "unbind",
			Usage: "unbind from a service instance",
			UsageText: "broker-cli [global options] unbind --instance-id " +
				"<instance id> --binding-id <binding id> [other command options]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagsInstanceID,
					Usage: "specify the `<instance id>`; required",
				},
				cli.StringFlag{
					Name:  flagsBindingID,
					Usage: "specify the `<binding id>`; required",
				},
			},
			Action: unbind,
		},
		{
			Name:  "deprovision",
			Usage: "deprovision a service instance",
			UsageText: "broker-cli [global options] deprovision --instance-id " +
				"<instance id> [other command options]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagsInstanceID,
					Usage: "specify the `<instance id>`; required",
				},
				cli.BoolFlag{
					Name: flagPoll,
					Usage: "poll the instance for status until deprovisioning succeeds " +
						"or fails",
				},
			},
			Action: deprovision,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\n%s\n\n", err)
		os.Exit(1)
	}
}
