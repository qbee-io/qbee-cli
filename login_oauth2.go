package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuth2DeviceAuthorizationResponse represents the response from the OAuth2 Device Authorization endpoint.
type OAuth2DeviceAuthorizationResponse struct {
	// DeviceCode is a long string used to verify the session between the client and the authorization server.
	// The client uses this parameter to request the access token from the authorization server.
	DeviceCode string `json:"device_code"`

	// UserCode is a short string shown to the user used to identify the session on a secondary device.
	UserCode string `json:"user_code"`

	// VerificationURI is the URL where the user can enter the UserCode to authorize the device.
	VerificationURI string `json:"verification_uri"`

	// ExpiresIn is the number of seconds before the device_code and user_code expire.
	ExpiresIn int `json:"expires_in"`

	// Interval is the minimum number of seconds that the client MUST wait between polling requests to the token endpoint.
	Interval int `json:"interval"`

	// ClientIP is the IP address of the client that initiated the device authorization request.
	ClientIP string `json:"client_ip"`

	// UserAgent is the user agent string of the client that initiated the device authorization request.
	UserAgent string `json:"user_agent"`
}

const oauth2DevicePath = "/oauth2/devicecode"

// OAuth2DeviceAuthorizationRequest initiates the OAuth2 Device Authorization Flow by requesting a device code.
func (c *Client) OAuth2DeviceAuthorizationRequest(ctx context.Context) (*OAuth2DeviceAuthorizationResponse, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+oauth2DevicePath, nil)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("User-Agent", UserAgent)
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected API response: %s", responseBody)
	}

	deviceAuth := new(OAuth2DeviceAuthorizationResponse)

	if err = json.Unmarshal(responseBody, deviceAuth); err != nil {
		return nil, fmt.Errorf("error decoding JSON response (%w): %s", err, responseBody)
	}

	return deviceAuth, nil
}

// OAuth2TokenResponse represents the OAuth2 token response.
type OAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var (
	// ErrOAuth2AuthorizationPending indicates that the authorization is still pending.
	// The client should continue polling the token endpoint after waiting the interval specified by the authorization server.
	ErrOAuth2AuthorizationPending = fmt.Errorf("authorization_pending")

	// ErrOAuth2AuthorizationDeclined indicates that the authorization was declined by the user.
	// No further attempts should be made to obtain a token.
	ErrOAuth2AuthorizationDeclined = fmt.Errorf("authorization_declined")

	// ErrOAuth2BadVerificationCode indicates that the provided verification code is invalid.
	// No further attempts should be made to obtain a token.
	ErrOAuth2BadVerificationCode = fmt.Errorf("bad_verification_code")

	// ErrOAuth2ExpiredToken indicates that the device code has expired.
	// No further attempts should be made to obtain a token.
	ErrOAuth2ExpiredToken = fmt.Errorf("expired_token")
)

const oauth2TokenPath = "/oauth2/token"

// OAuth2GetToken retrieves the OAuth2 token using the provided device code.
// It handles specific error responses related to the device authorization flow - retry or fail accordingly.
func (c *Client) OAuth2GetTokenForDeviceCode(ctx context.Context, deviceCode string) (*OAuth2TokenResponse, error) {
	payload := url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {deviceCode},
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+oauth2TokenPath, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("User-Agent", UserAgent)
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	defer func() { _ = httpResponse.Body.Close() }()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode == http.StatusBadRequest {
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err = json.Unmarshal(responseBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("error decoding JSON error response (%w): %s", err, responseBody)
		}

		switch errorResponse.Error {
		case "authorization_pending":
			return nil, ErrOAuth2AuthorizationPending
		case "authorization_declined":
			return nil, ErrOAuth2AuthorizationDeclined
		case "bad_verification_code":
			return nil, ErrOAuth2BadVerificationCode
		case "expired_token":
			return nil, ErrOAuth2ExpiredToken
		}
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected API response: %s", responseBody)
	}

	tokenResponse := new(OAuth2TokenResponse)

	if err = json.Unmarshal(responseBody, tokenResponse); err != nil {
		return nil, fmt.Errorf("error decoding JSON response (%w): %s", err, responseBody)
	}

	return tokenResponse, nil
}

// ApproveOAuth2DeviceAuthorization approves the device authorization using the provided user code.
func (c *Client) ApproveOAuth2DeviceAuthorization(ctx context.Context, userCode string) error {
	return c.Call(ctx, http.MethodPost, "/api/v2/oauth2/device-auth/"+userCode, nil, nil)
}

// DeclineOAuth2DeviceAuthorization declines the device authorization using the provided user code.
func (c *Client) DeclineOAuth2DeviceAuthorization(ctx context.Context, userCode string) error {
	return c.Call(ctx, http.MethodDelete, "/api/v2/oauth2/device-auth/"+userCode, nil, nil)
}

// InteractiveOAuth2DeviceAuthorizationFlow performs the interactive OAuth2 Device Authorization Flow.
func (c *Client) InteractiveOAuth2DeviceAuthorizationFlow(ctx context.Context) error {
	deviceAuth, err := c.OAuth2DeviceAuthorizationRequest(ctx)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("To authorize, visit: %s\n", deviceAuth.VerificationURI)
	fmt.Println("")
	fmt.Printf("And verify that the code  %s  is shown before authorization.\n", deviceAuth.UserCode)
	fmt.Println("")

	startTime := time.Now()
	expirationTime := startTime.Add(time.Duration(deviceAuth.ExpiresIn) * time.Second)
	pollingDuration := time.Duration(deviceAuth.Interval) * time.Second

	for time.Now().Before(expirationTime) {
		time.Sleep(pollingDuration)

		tokenResponse, err := c.OAuth2GetTokenForDeviceCode(ctx, deviceAuth.DeviceCode)
		switch err {
		case nil:
			c.WithAuthToken(tokenResponse.AccessToken)
			c.SetRefreshToken(tokenResponse.RefreshToken)
			return nil
		case ErrOAuth2AuthorizationPending:
			continue
		case ErrOAuth2AuthorizationDeclined:
			return fmt.Errorf("authorization was declined by the user")
		case ErrOAuth2ExpiredToken:
			return fmt.Errorf("authorization request has expired")
		default:
			return err
		}
	}

	return nil
}
