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
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	keyfunc "github.com/MicahParks/keyfunc/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	SHOPIFY_JWT                = "shopify_jwt"
	SHOPIFY_JWT_HEADER_ENV     = "SHOPIFY_JWT_HEADER"
	SHOPIFY_JWKS_URL_ENV       = "SHOPIFY_JWKS_URL"
	SHOPIFY_USER_ID_HEADER_ENV = "SHOPIFY_USER_ID_HEADER"
	SHOPIFY_AUTHZ_URL_ENV      = "SHOPIFY_AUTHZ_URL"
	SHOPIFY_AUTHZ_GROUPS_ENV   = "SHOPIFY_AUTHZ_GROUPS"
	SHOPIFY_AUTHZ_USERNAME_ENV = "SHOPIFY_AUTHZ_USERNAME"
	SHOPIFY_AUTHZ_PASSWORD_ENV = "SHOPIFY_AUTHZ_PASSWORD"
)

var errDenyShopifyJwt = errors.New("not allowed: shopify_jwt security_policy enforced")

type membership struct {
	Group  string
	Member bool
}

type shopifyAuthzData struct {
	Memberships []membership
}

type shopifyAuthzResponse struct {
	Data shopifyAuthzData
}

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

func buildAuth(username string, password string) string {
	return fmt.Sprintf("Basic %s", b64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))
}

func authorizeUser(userId int) (bool, error) {
	url := os.Getenv(SHOPIFY_AUTHZ_URL_ENV)

	var groups string

	for _, group := range strings.Split(os.Getenv(SHOPIFY_AUTHZ_GROUPS_ENV), ",") {
		groups += fmt.Sprintf("\"%s\",", group)
	}

	query := map[string]string{
		"query": fmt.Sprintf(`
			{
				memberships(
					userEmployeeId: %d
					groups: [%s]
				) {
					member
					group
				}
			}
		`, userId, groups),
	}

	queryJson, err := json.Marshal(query)

	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(queryJson))

	if err != nil {
		return false, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", buildAuth(os.Getenv(SHOPIFY_AUTHZ_USERNAME_ENV), os.Getenv(SHOPIFY_AUTHZ_PASSWORD_ENV)))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to authorize user. status: %s", res.Status)
	}

	var jsonRes shopifyAuthzResponse
	err = json.NewDecoder(res.Body).Decode(&jsonRes)
	if err != nil {
		return false, err
	}

	for _, membership := range jsonRes.Data.Memberships {
		if membership.Member {
			return true, nil
		}
	}

	return false, fmt.Errorf("user is not a member of any authorized groups")
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
		jwtToken := req.Header.Get(os.Getenv(SHOPIFY_JWT_HEADER_ENV))

		if len(jwtToken) < 1 {
			log.Println("failed to get jwt token from header")
			return errDenyShopifyJwt
		}

		_, err := validateJWT(jwtToken, os.Getenv(SHOPIFY_JWKS_URL_ENV))

		if err != nil {
			log.Printf("invalid JWT token provided: %s", err)
			return err
		}

		userId := req.Header.Get(os.Getenv(SHOPIFY_USER_ID_HEADER_ENV))

		if len(userId) < 1 {
			log.Println("failed to get user id from header")
			return errDenyShopifyJwt
		}

		userIdInt, err := strconv.Atoi(userId)

		if err != nil {
			log.Printf("failed to convert user id to int: %s", err)
			return err
		}

		authorized, err := authorizeUser(userIdInt)

		if err != nil {
			log.Printf("failed to authorize user ID: %s, %v", userId, err)
			return err
		}

		if !authorized {
			log.Printf("user ID %s is not authorized", userId)
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
