//go:build !windows

/***************************************************************
 *
 * Copyright (C) 2024, University of Nebraska-Lincoln
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
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/pelicanplatform/pelican/config"
	"github.com/pelicanplatform/pelican/namespaces"
	"github.com/pelicanplatform/pelican/test_utils"
)

func TestMain(m *testing.M) {
	if err := config.InitClient(); err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// TestIsPort calls main.hasPort with a hostname, checking
// for a valid return value.
func TestIsPort(t *testing.T) {

	if hasPort("blah.not.port:") {
		t.Fatal("Failed to parse port when : at end")
	}

	if !hasPort("host:1") {
		t.Fatal("Failed to parse with port = 1")
	}

	if hasPort("https://example.com") {
		t.Fatal("Failed when scheme is specified")
	}
}

// TestNewTransferDetails checks the creation of transfer details
func TestNewTransferDetails(t *testing.T) {
	os.Setenv("http_proxy", "http://proxy.edu:3128")

	// Case 1: cache with http
	testCache := namespaces.Cache{
		AuthEndpoint: "cache.edu:8443",
		Endpoint:     "cache.edu:8000",
		Resource:     "Cache",
	}
	transfers := newTransferDetails(testCache, transferDetailsOptions{false, ""})
	assert.Equal(t, 2, len(transfers))
	assert.Equal(t, "cache.edu:8000", transfers[0].Url.Host)
	assert.Equal(t, "http", transfers[0].Url.Scheme)
	assert.Equal(t, true, transfers[0].Proxy)
	assert.Equal(t, "cache.edu:8000", transfers[1].Url.Host)
	assert.Equal(t, "http", transfers[1].Url.Scheme)
	assert.Equal(t, false, transfers[1].Proxy)

	// Case 2: cache with https
	transfers = newTransferDetails(testCache, transferDetailsOptions{true, ""})
	assert.Equal(t, 1, len(transfers))
	assert.Equal(t, "cache.edu:8443", transfers[0].Url.Host)
	assert.Equal(t, "https", transfers[0].Url.Scheme)
	assert.Equal(t, false, transfers[0].Proxy)

	testCache.Endpoint = "cache.edu"
	// Case 3: cache without port with http
	transfers = newTransferDetails(testCache, transferDetailsOptions{false, ""})
	assert.Equal(t, 2, len(transfers))
	assert.Equal(t, "cache.edu:8000", transfers[0].Url.Host)
	assert.Equal(t, "http", transfers[0].Url.Scheme)
	assert.Equal(t, true, transfers[0].Proxy)
	assert.Equal(t, "cache.edu:8000", transfers[1].Url.Host)
	assert.Equal(t, "http", transfers[1].Url.Scheme)
	assert.Equal(t, false, transfers[1].Proxy)

	// Case 4. cache without port with https
	testCache.AuthEndpoint = "cache.edu"
	transfers = newTransferDetails(testCache, transferDetailsOptions{true, ""})
	assert.Equal(t, 2, len(transfers))
	assert.Equal(t, "cache.edu:8444", transfers[0].Url.Host)
	assert.Equal(t, "https", transfers[0].Url.Scheme)
	assert.Equal(t, false, transfers[0].Proxy)
	assert.Equal(t, "cache.edu:8443", transfers[1].Url.Host)
	assert.Equal(t, "https", transfers[1].Url.Scheme)
	assert.Equal(t, false, transfers[1].Proxy)
}

func TestNewTransferDetailsEnv(t *testing.T) {

	testCache := namespaces.Cache{
		AuthEndpoint: "cache.edu:8443",
		Endpoint:     "cache.edu:8000",
		Resource:     "Cache",
	}

	os.Setenv("OSG_DISABLE_PROXY_FALLBACK", "")
	err := config.InitClient()
	assert.Nil(t, err)
	transfers := newTransferDetails(testCache, transferDetailsOptions{false, ""})
	assert.Equal(t, 1, len(transfers))
	assert.Equal(t, true, transfers[0].Proxy)

	transfers = newTransferDetails(testCache, transferDetailsOptions{true, ""})
	assert.Equal(t, 1, len(transfers))
	assert.Equal(t, "https", transfers[0].Url.Scheme)
	assert.Equal(t, false, transfers[0].Proxy)
	os.Unsetenv("OSG_DISABLE_PROXY_FALLBACK")
	viper.Reset()
	err = config.InitClient()
	assert.Nil(t, err)
}

func TestSlowTransfers(t *testing.T) {
	ctx, _, _ := test_utils.TestContext(context.Background(), t)

	// Adjust down some timeouts to speed up the test
	viper.Set("Client.SlowTransferWindow", 5)
	viper.Set("Client.SlowTransferRampupTime", 10)

	channel := make(chan bool)
	slowDownload := 1024 * 10 // 10 KiB/s < 100 KiB/s
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Add("Content-Length", "1024000")
			w.WriteHeader(http.StatusOK)
			return
		}
		buffer := make([]byte, slowDownload)
		for {
			select {
			case <-channel:
				return
			default:
				_, err := w.Write(buffer)
				if err != nil {
					return
				}
				w.(http.Flusher).Flush()
				time.Sleep(1 * time.Second)
			}
		}
	}))

	defer svr.CloseClientConnections()
	defer svr.Close()

	testCache := namespaces.Cache{
		AuthEndpoint: svr.URL,
		Endpoint:     svr.URL,
		Resource:     "Cache",
	}
	transfers := newTransferDetails(testCache, transferDetailsOptions{false, ""})
	assert.Equal(t, 2, len(transfers))
	assert.Equal(t, svr.URL, transfers[0].Url.String())

	finishedChannel := make(chan bool)
	var err error
	// Do a quick timeout
	go func() {
		_, _, _, err = downloadHTTP(ctx, nil, nil, transfers[0], filepath.Join(t.TempDir(), "test.txt"), "", nil)
		finishedChannel <- true
	}()

	select {
	case <-finishedChannel:
		if err == nil {
			t.Fatal("Error is nil, download should have failed")
		}
	case <-time.After(time.Second * 160):
		// 120 seconds for warmup, 30 seconds for download
		t.Fatal("Maximum downloading time reach, download should have failed")
	}

	// Close the channel to allow the download to complete
	channel <- true

	// Make sure the errors are correct
	assert.NotNil(t, err)
	assert.IsType(t, &SlowTransferError{}, err)
}

// Test stopped transfer
func TestStoppedTransfer(t *testing.T) {
	os.Setenv("http_proxy", "http://proxy.edu:3128")
	ctx, _, _ := test_utils.TestContext(context.Background(), t)

	// Adjust down the timeouts
	viper.Set("Client.StoppedTransferTimeout", 3)
	viper.Set("Client.SlowTransferRampupTime", 100)

	channel := make(chan bool)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Add("Content-Length", "102400")
			w.WriteHeader(http.StatusOK)
			return
		}
		buffer := make([]byte, 1024*100)
		for {
			select {
			case <-channel:
				return
			default:
				_, err := w.Write(buffer)
				if err != nil {
					return
				}
				w.(http.Flusher).Flush()
				time.Sleep(1 * time.Second)
				buffer = make([]byte, 0)
			}
		}
	}))

	defer svr.CloseClientConnections()
	defer svr.Close()

	testCache := namespaces.Cache{
		AuthEndpoint: svr.URL,
		Endpoint:     svr.URL,
		Resource:     "Cache",
	}
	transfers := newTransferDetails(testCache, transferDetailsOptions{false, ""})
	assert.Equal(t, 2, len(transfers))
	assert.Equal(t, svr.URL, transfers[0].Url.String())

	finishedChannel := make(chan bool)
	var err error

	go func() {
		_, _, _, err = downloadHTTP(ctx, nil, nil, transfers[0], filepath.Join(t.TempDir(), "test.txt"), "", nil)
		finishedChannel <- true
	}()

	select {
	case <-finishedChannel:
		if err == nil {
			t.Fatal("Download should have failed")
		}
	case <-time.After(time.Second * 150):
		t.Fatal("Download should have failed")
	}

	// Close the channel to allow the download to complete
	channel <- true

	// Make sure the errors are correct
	assert.NotNil(t, err)
	assert.IsType(t, &StoppedTransferError{}, err, err.Error())
}

// Test connection error
func TestConnectionError(t *testing.T) {
	ctx, _, _ := test_utils.TestContext(context.Background(), t)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("dialClosedPort: Listen failed: %v", err)
	}
	addr := l.Addr().String()
	l.Close()

	_, _, _, err = downloadHTTP(ctx, nil, nil, transferAttemptDetails{Url: &url.URL{Host: addr, Scheme: "http"}, Proxy: false}, filepath.Join(t.TempDir(), "test.txt"), "", nil)

	assert.IsType(t, &ConnectionSetupError{}, err)

}

func TestTrailerError(t *testing.T) {
	ctx, _, _ := test_utils.TestContext(context.Background(), t)

	// Set up an HTTP server that returns an error trailer
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Trailer", "X-Transfer-Status")
		w.Header().Set("X-Transfer-Status", "500: Unable to read test.txt; input/output error")

		chunkedWriter := httputil.NewChunkedWriter(w)
		defer chunkedWriter.Close()

		_, err := chunkedWriter.Write([]byte("Test data"))
		if err != nil {
			t.Fatalf("Error writing to chunked writer: %v", err)
		}
	}))

	defer svr.Close()

	testCache := namespaces.Cache{
		AuthEndpoint: svr.URL,
		Endpoint:     svr.URL,
		Resource:     "Cache",
	}
	transfers := newTransferDetails(testCache, transferDetailsOptions{false, ""})
	assert.Equal(t, 2, len(transfers))
	assert.Equal(t, svr.URL, transfers[0].Url.String())

	// Call DownloadHTTP and check if the error is returned correctly
	_, _, _, err := downloadHTTP(ctx, nil, nil, transfers[0], filepath.Join(t.TempDir(), "test.txt"), "", nil)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "transfer error: Unable to read test.txt; input/output error")
}

func TestUploadZeroLengthFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//t.Logf("%s", dump)
		assert.Equal(t, "PUT", r.Method, "Not PUT Method")
		assert.Equal(t, int64(0), r.ContentLength, "ContentLength should be 0")
	}))
	defer ts.Close()
	reader := bytes.NewReader([]byte{})
	request, err := http.NewRequest("PUT", ts.URL, reader)
	if err != nil {
		assert.NoError(t, err)
	}

	request.Header.Set("Authorization", "Bearer test")
	errorChan := make(chan error, 1)
	responseChan := make(chan *http.Response)
	go runPut(request, responseChan, errorChan)
	select {
	case err := <-errorChan:
		assert.NoError(t, err)
	case response := <-responseChan:
		assert.Equal(t, http.StatusOK, response.StatusCode)
	case <-time.After(time.Second * 2):
		assert.Fail(t, "Timeout while waiting for response")
	}
}

func TestFailedUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//t.Logf("%s", dump)
		assert.Equal(t, "PUT", r.Method, "Not PUT Method")
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error"))
		assert.NoError(t, err)

	}))
	defer ts.Close()
	reader := strings.NewReader("test")
	request, err := http.NewRequest("PUT", ts.URL, reader)
	if err != nil {
		assert.NoError(t, err)
	}
	request.Header.Set("Authorization", "Bearer test")
	errorChan := make(chan error, 1)
	responseChan := make(chan *http.Response)
	go runPut(request, responseChan, errorChan)
	select {
	case err := <-errorChan:
		assert.Error(t, err)
	case response := <-responseChan:
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	case <-time.After(time.Second * 2):
		assert.Fail(t, "Timeout while waiting for response")
	}
}
