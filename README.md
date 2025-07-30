# qbee-cli

qbee-cli is a client library and a command line tool used to interact with [qbee.io](https://qbee.io) IoT/Linux device management platform.


# Use as command line tool

## Download the binary

Open the [releases](https://github.com/qbee-io/qbee-cli/releases) page, scroll down to **Assets** and download the latest version for your platform.

*NOTE:* The binary is not signed, so you might need to allow it to run in your system settings. For Windows, you need to add ".exe" to the binary name.

## Build the binary

```shell
go build -o qbee-cli ./cmd
```

## Providing credentials

qbee-cli supports multiple authentication methods to provide flexibility for different use cases:

### 1. Interactive Login (Default)

```shell
qbee-cli login
```

This will prompt for your email and password interactively. If two-factor authentication is enabled, you'll be prompted to select a 2FA provider and enter the code.

### 2. Environment Variables for Login

You can provide credentials via environment variables to avoid interactive prompts:

```shell
export QBEE_EMAIL=alice@example.com
export QBEE_PASSWORD=secret
qbee-cli login
```

For accounts with two-factor authentication enabled, you can also set:

```shell
export QBEE_EMAIL=alice@example.com
export QBEE_PASSWORD=secret
export QBEE_2FA_CODE=123456
qbee-cli login
```

If you need to authenticate against a different qbee.io instance, you can set the `QBEE_BASEURL` environment variable:

```shell
export QBEE_EMAIL=alice@example.com
export QBEE_PASSWORD=secret
export QBEE_BASEURL=https://www.app.qbee.example.com
qbee-cli login
```

### 3. Authentication Token (QBEE_TOKEN)

For automated workflows and multiple command executions, you can use authentication tokens. First, obtain a token using the `--print-token` flag:

```shell
# Get a token and save it for reuse
QBEE_TOKEN=$(QBEE_EMAIL=alice@example.com QBEE_PASSWORD=secret QBEE_2FA_CODE=123456 qbee-cli login --print-token)
export QBEE_TOKEN
```

Then use the token for subsequent commands without needing to re-authenticate:

```shell
# Use the token for any qbee-cli command
QBEE_TOKEN=your_token_here qbee-cli device list
QBEE_TOKEN=your_token_here qbee-cli files list /
```

Or export it once and use multiple commands:

```shell
export QBEE_TOKEN=your_token_here
qbee-cli device list
qbee-cli files list /
qbee-cli connect -d device123 -t 8080:localhost:80
```

**Note:** The `--print-token` flag prints the authentication token to stdout instead of writing the configuration file to disk, making it ideal for automation and CI/CD workflows.

Please remember to rotate your credentials regularly and keep tokens secure.

## Run latest using Go

```shell
go run go.qbee.io/client/cmd@latest
```

## Remote access using qbee-cli

```shell
export QBEE_EMAIL=<email>
export QBEE_PASSWORD=<password>

qbee-cli connect -d <deviceID> -t <target>[,<target> ...]
```

Where:
- `deviceID` identifies to which device we want to connect (public key digest)
- `target` defines a singe port forwarding target as `[<localHost>:]<localPort>:<remoteHost>:<remotePort>`
- `localHost` is optional and defaults to _localhost_ to only listen on the loopback interface
- `localPort` is the local port on which tunnel will listen for connections/packets
- `remoteHost` must be set to _localhost_
- `remotePort` is the remote port on which tunnel will connect to on the device

# Use as a Go module

```go
package demo

import (
	"context"
	"log"
	"os"

	"go.qbee.io/client"
)

func main() {
	cli := client.New()
	ctx := context.Background()

	email := os.Getenv("QBEE_EMAIL")
	password := os.Getenv("QBEE_PASSWORD")

	if err := cli.Authenticate(ctx, email, password); err != nil {
		log.Fatalf("authentication failed: %v", err)
	}
}
```

## Contributing

We welcome contributions to this project!
Please see [Contribution License Agreement](https://qbee.io/docs/contribution-license-agreement.html) for more information.
