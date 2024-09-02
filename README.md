# qbee-cli

qbee-cli is a client library and a command line tool used to interact with [qbee.io](https://qbee.io) IoT/Linux device management platform.


# Use as command line tool

## Build the binary

```shell
go build -o qbee-cli ./cmd
```

## Providing credentials

Currently, the only way to provide credentials is through environmental variables: `QBEE_EMAIL` & `QBEE_PASSWORD`. 

If your account is configured with two-factor authentication, you will either be prompted for which of your configured
2FA providers you want to use, or you can set the `QBEE_2FA_CODE` environment variable to provide a code for the 
Google provider directly.

Please remember to rotate your credentials regularly.

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
