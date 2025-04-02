package main

import (
	"context"
	"fmt"
	"os"

	"go.qbee.io/client"
	"go.qbee.io/client/broker"
)

const (
	brokerPasswordOption       = "password"
	brokerUsernameOption       = "username"
	brokerAuthTokenOption      = "auth-token"
	brokerListenPortOption     = "listen-port"
	brokerRemoteHostOption     = "remote-host"
	brokerRemotePortOption     = "remote-port"
	brokerRemoteProtocolOption = "remote-protocol"
	brokerBaseUrlOption        = "base-url"
)

var brokerCommand = Command{
	Description: "Start the http broker",
	Options: []Option{
		{
			Name:     brokerUsernameOption,
			Short:    "u",
			Help:     "Username for authentication",
			Required: false,
		},
		{
			Name:     brokerPasswordOption,
			Short:    "p",
			Help:     "Password for authentication",
			Required: false,
		},
		{
			Name:   brokerBaseUrlOption,
			Short:  "b",
			Help:   "Base URL for the broker",
			Hidden: true,
		},
		{
			Name:     brokerAuthTokenOption,
			Help:     "Token for authentication",
			Required: false,
		},
		{
			Name:    brokerListenPortOption,
			Help:    "Port to listen on",
			Default: broker.DefaultListenPort,
		},
		{
			Name:    brokerRemoteHostOption,
			Help:    "Remote host to connect to",
			Default: broker.DefaultRemoteHost,
		},
		{
			Name:    brokerRemotePortOption,
			Help:    "Default port to connect to",
			Default: broker.DefaultRemotePort,
		},
		{
			Name:    brokerRemoteProtocolOption,
			Help:    "Default protocol to use",
			Default: broker.DefaultRemoteProtocol,
			Hidden:  true,
		},
	},
	OptionsHandler: func(opts Options) error {

		if opts["username"] != "" {
			err := os.Setenv("QBEE_USERNAME", opts[brokerUsernameOption])
			if err != nil {
				return err
			}
		}

		if opts["password"] != "" {
			err := os.Setenv("QBEE_PASSWORD", opts[brokerPasswordOption])
			if err != nil {
				return err
			}
		}

		if opts["base-url"] != "" {
			os.Setenv("QBEE_BASEURL", opts[brokerBaseUrlOption])
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

		s := broker.NewService().WithClient(cli).WithListenPort(opts[brokerListenPortOption])

		if opts["auth-token"] != "" {
			s = s.WithAuthToken(opts[brokerAuthTokenOption])
		}

		if opts["remote-port"] != "" {
			s = s.WithRemotePort(opts[brokerRemotePortOption])
		}

		if opts["remote-protocol"] != "" {
			s = s.WithRemoteProtocol(opts[brokerRemoteProtocolOption])
		}

		return s.Start(ctx)
	},
}
