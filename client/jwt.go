package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Token contains typical values returned from REST API.
type Token struct {
	Raw       string                 `json:"raw"`       // The raw token.  Populated when you Parse a token
	Header    map[string]interface{} `json:"headers"`   // The first segment of the token
	Claims    StandardClaims         `json:"claims"`    // The second segment of the token
	Signature string                 `json:"signature"` // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   `json:"valid"`     // Is the token valid?  Populated when you Parse/Verify a token
}

// qbee claims
// {
// 	"iat": 1701870535,
// 	"exp": 1702070535,
// 	"roles": [
// 	  "1925c605-87bf-46aa-9fcb-254cc5162722",
// 	  "b48152fe-74af-4810-bc67-482fa93ebfa3",
// 	  "97e95739-2c59-4ebf-ae04-ca1e5ada4983"
// 	],
// 	"username": "jonhenrik@qbee.io",
// 	"ip": "10.0.13.113"
// }

type StandardClaims struct {
	Audience  string   `json:"aud,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
	Id        string   `json:"jti,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	NotBefore int64    `json:"nbf,omitempty"`
	Subject   string   `json:"sub,omitempty"`
	Name      string   `json:"name,omitempty"`
	UserName  string   `json:"username,omitempty"`
	Ip        string   `json:"ip,omitempty"`
	Roles     []string `json:"roles,omitempty"`
}

func decodeTokenSegment(segment string) ([]byte, error) {
	if l := len(segment) % 4; l > 0 {
		segment += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(segment)
}

func DecodeAccessToken(accessToken string, claims StandardClaims) (*Token, error) {
	var err error
	token := &Token{Raw: accessToken}
	parts := strings.Split(accessToken, ".")
	var headerBytes []byte
	headerBytes, err = decodeTokenSegment(parts[0])
	if err = json.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, fmt.Errorf("An error occurred unmarshalling token")
	}
	var claimBytes []byte
	if claimBytes, err = decodeTokenSegment(parts[1]); err != nil {
		return token, fmt.Errorf("An error occurred unmarshalling token")
	}
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	err = dec.Decode(&claims)
	token.Claims = claims
	return token, err
}
