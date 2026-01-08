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
	loginPrintToken   = "print-token"
)

// Login is the login command.

var loginCommand = Command{
	Description: "Login to Qbee.io",
	Options: []Option{
		{
			Name:     loginUserEmail,
			Short:    "u",
			Help:     "User email. Can also be set via QBEE_EMAIL environment variable.",
			Required: false,
		},
		{
			Name:     loginUserPassword,
			Short:    "p",
			Help:     "User password. Can also be set via QBEE_PASSWORD environment variable.",
			Required: false,
		},
		{
			Name:    loginBaseURL,
			Short:   "b",
			Help:    "Qbee.io base URL.",
			Default: "https://www.app.qbee.io",
			Hidden:  true,
		},
		{
			Name: loginPrintToken,
			Help: "Print the authentication token to stdout instead of writing configuration file.",
			Flag: "true",
		},
	},
	Target: func(opts Options) error {
		// Check for environment variables first
		email := os.Getenv("QBEE_EMAIL")
		password := os.Getenv("QBEE_PASSWORD")
		baseURL := os.Getenv("QBEE_BASEURL")

		// Use command line options if environment variables are not set
		if email == "" {
			email = opts[loginUserEmail]
		}
		if password == "" {
			password = opts[loginUserPassword]
		}
		if baseURL == "" {
			baseURL = opts[loginBaseURL]
		}

		ctx := context.Background()
		cli := client.New().WithBaseURL(baseURL)

		// Validate that we have required information
		if email == "" {
			if err := cli.InteractiveOAuth2DeviceAuthorizationFlow(ctx); err != nil {
				return err
			}
		} else {
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
		}

		// Check if user wants to print token instead of writing config
		if opts[loginPrintToken] == "true" {
			fmt.Println(cli.GetAuthToken())
			return nil
		}

		loginConfig := client.LoginConfig{
			AuthToken:    cli.GetAuthToken(),
			RefreshToken: cli.GetRefreshToken(),
			BaseURL:      cli.GetBaseURL(),
		}

		if err := client.LoginWriteConfig(loginConfig); err != nil {
			return err
		}

		fmt.Printf("Successfully authenticated!\n")
		return nil
	},
}
