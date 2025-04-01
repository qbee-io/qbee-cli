package main

import (
	"context"
	"fmt"
	"os"

	"go.qbee.io/client"
	"go.qbee.io/client/broker"
)

var brokerCommand = Command{
	Description: "Start the http broker",
	Options: []Option{
		{
			Name:    "port",
			Short:   "p",
			Help:    "Port to listen on",
			Default: "8080",
		},
		{
			Name:     "username",
			Help:     "Username for authentication",
			Required: false,
		},
		{
			Name:     "password",
			Help:     "Password for authentication",
			Required: false,
		},
		{
			Name:     "auth-token",
			Help:     "Token for authentication",
			Required: false,
		},
		{
			Name:   "base-url",
			Help:   "Base URL for the broker",
			Hidden: true,
		},
		{
			Name:    "remote-host",
			Help:    "Remote host to connect to",
			Default: "localhost",
		},
		{
			Name:    "remote-port",
			Help:    "Default port to connect to",
			Default: "80",
		},
		{
			Name:    "remote-protocol",
			Help:    "Default protocol to use",
			Default: "http",
		},
	},
	OptionsHandler: func(opts Options) error {

		if opts["username"] != "" {
			err := os.Setenv("QBEE_USERNAME", opts["username"])
			if err != nil {
				return err
			}
		}

		if opts["password"] != "" {
			err := os.Setenv("QBEE_PASSWORD", opts["password"])
			if err != nil {
				return err
			}
		}

		if opts["base-url"] != "" {
			os.Setenv("QBEE_BASEURL", opts["base-url"])
		}

		if os.Getenv("QBEE_PASSWORD") == "" {
			return fmt.Errorf("no password provided")
		}

		if os.Getenv("QBEE_USERNAME") == "" {
			return fmt.Errorf("no username provided")
		}

		return nil
	},
	Target: func(opts Options) error {
		ctx := context.Background()

		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		s := broker.NewService().WithClient(cli).WithPort(opts["port"])

		if opts["auth-token"] != "" {
			s = s.WithAuthToken(opts["auth-token"])
		}

		if opts["remote-port"] != "" {
			s = s.WithRemotePort(opts["remote-port"])
		}

		if opts["remote-protocol"] != "" {
			s = s.WithRemoteProtocol(opts["remote-protocol"])
		}

		return s.Start(ctx)
	},
}
