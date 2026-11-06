// Copyright 2026 qbee.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"regexp"
	"testing"
)

// secretRefPattern matches a secret reference, e.g. secret:ddb285d5-4f9c-4815-8fcf-7e6c4ba1d5ce (UUID4).
var secretRefPattern = regexp.MustCompile(
	`^secret:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func Test_CreateSecret(t *testing.T) {
	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	secretRef, err := cli.CreateSecret(ctx, "qbee-cli-integration-test-secret")
	if err != nil {
		t.Fatal(err)
	}

	if !secretRefPattern.MatchString(secretRef) {
		t.Fatalf("expected secret_ref in format secret:<uuid4>, got %q", secretRef)
	}
}
