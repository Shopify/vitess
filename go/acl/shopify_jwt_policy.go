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
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	keyfunc "github.com/MicahParks/keyfunc/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	SHOPIFY_JWT             = "shopify_jwt"
	SHOPIFY_COOKIE_NAME_ENV = "JWT_COOKIE_NAME"
	SHOPIFY_JWKS_URL_ENV    = "JWKS_URL"
)

var errDenyShopifyJwt = errors.New("not allowed: shopify_jwt security_policy enforced")

func jwksRequestFactory(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewReader(nil))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}

type shopifyJwt struct{}

func validateJWT(tokenString string, jwksURL string) (bool, error) {
	// Fetch the JWKS from the provided URL
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RequestFactory: jwksRequestFactory,
	})

	defer jwks.EndBackground()

	if err != nil {
		return false, err
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, jwks.Keyfunc)
	if err != nil {
		return false, fmt.Errorf("failed to parse JWT token: %v", err)
	}

	if token.Valid {
		return true, nil
	}

	return false, nil
}

// CheckAccessActor disallows actor access not verified by shopifyJwt
func (shopifyJwt) CheckAccessActor(actor, role string) error {
	switch role {
	case SHOPIFY_JWT:
		return nil
	default:
		return errDenyShopifyJwt
	}
}

// CheckAccessHTTP disallows HTTP access not verified by shopifyJwt
func (shopifyJwt) CheckAccessHTTP(req *http.Request, role string) error {
	switch role {
	case SHOPIFY_JWT:
		jwtCookie, err := req.Cookie(os.Getenv("SHOPIFY_COOKIE_NAME_ENV"))

		if err != nil {
			log.Printf("failed to get jwt token from cookie: %s", err)
			return err
		}

		_, err = validateJWT(jwtCookie.Value, os.Getenv(SHOPIFY_JWKS_URL_ENV))

		if err != nil {
			log.Printf("invalid JWT token provided: %s", err)
			return err
		}

		return nil
	default:
		return errDenyShopifyJwt
	}
}

func init() {
	RegisterPolicy(SHOPIFY_JWT, shopifyJwt{})
}
