package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.qbee.io/client"
	"golang.org/x/term"
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
			Required: false,
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

		if password == "" {
			fmt.Printf("Enter password for %s: ", email)
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			// Print a newline to simulate the enter key press
			fmt.Println()
			password = strings.TrimSpace(string(bytePassword))
		}

		if err := cli.Authenticate(ctx, email, password); err != nil {
			return err
		}

		loginConfig := client.LoginConfig{
			AuthToken: cli.GetAuthToken(),
			BaseURL:   cli.GetBaseURL(),
		}

		if err := client.LoginWriteConfig(loginConfig); err != nil {
			return err
		}

		fmt.Printf("Successfully logged in as %s\n", email)
		return nil
	},
}
