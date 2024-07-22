package main

import (
	"context"
	"encoding/json"

	"go.qbee.io/client"
)

const (
	execDeviceOption  = "device"
	execCommandOption = "command"
)

var execCommand = Command{
	Description: "Execute a command on a device",
	Options: []Option{
		{
			Name:     execDeviceOption,
			Short:    "d",
			Help:     "Device ID",
			Required: true,
		},
		{
			Name:     execCommandOption,
			Short:    "c",
			Help:     "Command to execute as JSON string",
			Required: true,
		},
	},
	Target: func(opts Options) error {
		ctx := context.Background()
		cli, err := client.LoginGetAuthenticatedClient(ctx)
		if err != nil {
			return err
		}

		deviceID := opts[execDeviceOption]
		command := opts[execCommandOption]

		var cmd []string

		if err := json.Unmarshal([]byte(command), &cmd); err != nil {
			return err
		}

		if err := cli.ExecuteCommandStream(ctx, deviceID, cmd); err != nil {
			return err
		}

		return nil
	},
}
