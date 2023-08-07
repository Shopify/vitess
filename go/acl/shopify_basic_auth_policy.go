/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package acl

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	SHOPIFY_BASIC_AUTH          = "shopify_basic_auth"
	shopifyBasicAuthUsernameEnv = "SHOPIFY_BASIC_AUTH_USERNAME"
	shopifyBasicAuthPasswordEnv = "SHOPIFY_BASIC_AUTH_PASSWORD"
)

var errDenyShopifyBasicAuth = errors.New("not allowed: shopify_basic_auth security_policy enforced")

// The http API does not support https so using basic auth should not be considered secure as it's possible
// for the credentials to be sent in plain text.
type shopifyBasicAuth struct{}

func validateBasicAuth(username string, password string) bool {
	usernameEnv := strings.TrimSpace(os.Getenv(shopifyBasicAuthUsernameEnv))
	passwordEnv := strings.TrimSpace(os.Getenv(shopifyBasicAuthPasswordEnv))

	return (username == usernameEnv && password == passwordEnv)
}

// CheckAccessActor disallows actor access not verified by shopifyBasicAuth
func (shopifyBasicAuth) CheckAccessActor(actor, role string) error {
	switch role {
	case SHOPIFY_BASIC_AUTH:
		return nil
	default:
		return errDenyShopifyBasicAuth
	}
}

// CheckAccessHTTP disallows HTTP access not verified by shopifyBasicAuth
func (shopifyBasicAuth) CheckAccessHTTP(req *http.Request, role string) error {
	switch role {
	case SHOPIFY_BASIC_AUTH:
		username, password, ok := req.BasicAuth()

		if !ok {
			log.Printf("failed to parse basic auth headers")
			return errDenyShopifyBasicAuth
		}

		if !validateBasicAuth(strings.TrimSpace(username), strings.TrimSpace(password)) {
			log.Printf("username and password is invalid: %s:%s", username, password)
			return errDenyShopifyBasicAuth
		}

		return nil
	default:
		return errDenyShopifyBasicAuth
	}
}

func init() {
	RegisterPolicy(SHOPIFY_BASIC_AUTH, shopifyBasicAuth{})
}
