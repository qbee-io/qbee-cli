package cmd

import (
	"context"

	"github.com/qbee-io/qbee-cli/client"
)

const (
	loginUserEmail    = "email"
	loginUserPassword = "password"
	loginBaseURL      = "base-url"
)

// Login is the login command.

var loginCommand = Command{
	Description: "Login to Qbee.io",
	Options: []Option{
		{
			Name:     loginUserEmail,
			Short:    "u",
			Help:     "User email.",
			Required: true,
		},
		{
			Name:     loginUserPassword,
			Short:    "p",
			Help:     "User password.",
			Required: true,
		},
		{
			Name:    loginBaseURL,
			Short:   "b",
			Help:    "Qbee.io base URL.",
			Default: "https://www.app.qbee.io",
			Hidden:  true,
		},
	},
	Target: func(opts Options) error {
		email := opts[loginUserEmail]
		password := opts[loginUserPassword]
		baseURL := opts[loginBaseURL]

		ctx := context.Background()
		cli := client.New().WithBaseURL(baseURL)

		if err := cli.Authenticate(ctx, email, password); err != nil {
			return err
		}

		config := Config{
			AuthToken: cli.GetAuthToken(),
			BaseURL:   cli.GetBaseURL(),
		}

		if err := WriteConfig(config); err != nil {
			return err
		}

		return nil
	},
}
