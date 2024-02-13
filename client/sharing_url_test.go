/***************************************************************
 *
 * Copyright (C) 2023, University of Nebraska-Lincoln
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you
 * may not use this file except in compliance with the License.  You may
 * obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 ***************************************************************/

package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/pelicanplatform/pelican/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectorGeneration(t *testing.T) {
	returnError := false
	returnErrorRef := &returnError

	handler := func(w http.ResponseWriter, r *http.Request) {
		discoveryConfig := `{"director_endpoint": "https://location.example.com", "namespace_registration_endpoint": "https://location.example.com/namespace", "jwks_uri": "https://location.example.com/jwks"}`
		if *returnErrorRef {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(discoveryConfig))
			assert.NoError(t, err)
		}
	}
	server := httptest.NewTLSServer(http.HandlerFunc(handler))
	defer server.Close()
	serverURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	objectUrl := url.URL{
		Scheme: "pelican",
		Host:   serverURL.Host,
		Path:   "/test/foo",
	}

	// Discovery works to get URL
	viper.Reset()
	viper.Set("TLSSkipVerify", true)
	err = config.InitClient()
	require.NoError(t, err)
	dUrl, err := getDirectorFromUrl(&objectUrl)
	require.NoError(t, err)
	assert.Equal(t, dUrl, "https://location.example.com")

	// Discovery URL overrides the federation config.
	viper.Reset()
	viper.Set("TLSSkipVerify", true)
	viper.Set("Federation.DirectorURL", "https://location2.example.com")
	dUrl, err = getDirectorFromUrl(&objectUrl)
	require.NoError(t, err)
	assert.Equal(t, dUrl, "https://location.example.com")

	// Fallback to configuration if no discovery present
	viper.Reset()
	viper.Set("Federation.DirectorURL", "https://location2.example.com")
	objectUrl.Host = ""
	dUrl, err = getDirectorFromUrl(&objectUrl)
	require.NoError(t, err)
	assert.Equal(t, dUrl, "https://location2.example.com")

	// Error if server has an error
	viper.Reset()
	returnError = true
	viper.Set("TLSSkipVerify", true)
	objectUrl.Host = serverURL.Host
	_, err = getDirectorFromUrl(&objectUrl)
	require.Error(t, err)

	// Error if neither config nor hostname provided.
	viper.Reset()
	objectUrl.Host = ""
	_, err = getDirectorFromUrl(&objectUrl)
	require.Error(t, err)

	// Error on unknown scheme
	viper.Reset()
	objectUrl.Scheme = "buzzard"
	_, err = getDirectorFromUrl(&objectUrl)
	require.Error(t, err)
}

func TestSharingUrl(t *testing.T) {
	// Construct a local server that we can poke with QueryDirector
	myUrl := "http://redirect.com"
	myUrlRef := &myUrl
	log.SetLevel(log.DebugLevel)
	handler := func(w http.ResponseWriter, r *http.Request) {
		issuerLoc := *myUrlRef + "/issuer"

		if strings.HasPrefix(r.URL.Path, "/test") {
			w.Header().Set("Location", *myUrlRef)
			w.Header().Set("X-Pelican-Namespace", "namespace=/test, require-token=true")
			w.Header().Set("X-Pelican-Authorization", fmt.Sprintf("issuer=%s", issuerLoc))
			w.Header().Set("X-Pelican-Token-Generation", fmt.Sprintf("issuer=%s, base-path=/test, strategy=OAuth2", issuerLoc))
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else if r.URL.Path == "/issuer/.well-known/openid-configuration" {
			w.WriteHeader(http.StatusOK)
			oidcConfig := fmt.Sprintf(`{"token_endpoint": "%s/token", "registration_endpoint": "%s/register", "grant_types_supported": ["urn:ietf:params:oauth:grant-type:device_code"], "device_authorization_endpoint": "%s/device_authz"}`, issuerLoc, issuerLoc, issuerLoc)
			_, err := w.Write([]byte(oidcConfig))
			assert.NoError(t, err)
		} else if r.URL.Path == "/issuer/register" {
			//requestBytes, err := io.ReadAll(r.Body)
			//assert.NoError(t, err)
			clientConfig := `{"client_id": "client1", "client_secret": "secret", "client_secret_expires_at": 0}`
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(clientConfig))
			assert.NoError(t, err)
		} else if r.URL.Path == "/issuer/device_authz" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"device_code": "1234", "user_code": "5678", "interval": 1, "verification_uri": "https://example.com", "expires_in": 20}`))
			assert.NoError(t, err)
		} else if r.URL.Path == "/issuer/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"access_token": "token1234", "token_type": "jwt"}`))
			assert.NoError(t, err)
		} else {
			fmt.Println(r)
			requestBytes, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			fmt.Println(string(requestBytes))
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	myUrl = server.URL

	os.Setenv("PELICAN_SKIP_TERMINAL_CHECK", "password")
	defer os.Unsetenv("PELICAN_SKIP_TERMINAL_CHECK")
	viper.Set("Federation.DirectorURL", myUrl)
	viper.Set("ConfigDir", t.TempDir())
	err := config.InitClient()
	assert.NoError(t, err)

	// Call QueryDirector with the test server URL and a source path
	testUrl, err := url.Parse("/test/foo/bar")
	require.NoError(t, err)
	os.Setenv(config.GetPreferredPrefix()+"_SKIP_TERMINAL_CHECK", "true")
	token, err := CreateSharingUrl(testUrl, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	fmt.Println(token)
	os.Unsetenv(config.GetPreferredPrefix() + "_SKIP_TERMINAL_CHECK")
}
