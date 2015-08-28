// Package onesky tests
// Copyright (c) 2015 Sebastian Czoch <sebastian@czoch.eu>. All rights reserved.
// Use of this source code is governed by a GNU v2 license found in the LICENSE file.
package onesky

import (
	"net/url"
	"regexp"
	"testing"
)

// TestGetURLForEndpoint is testing getURLForEndpoint method
func TestGetURLForEndpoint(t *testing.T) {
	client := Client{}
	client.Secret = "test_secret"
	client.ProjectID = 1

	url, err := client.getURLForEndpoint("not_exits_endpoint")
	if err == nil {
		t.Errorf("getURLForEndpoint() = %+v, %+v, want %+v", url, err, "error")
	}

	want := APIAddress + "/" + APIVersion + "/" + "projects/1/translations"
	url, err = client.getURLForEndpoint("getFile")
	if url != want {
		t.Errorf("getURLForEndpoint() = %+v, %+v, want %+v", url, err, want)
	}
}

// TestGetFinalEndpointURL is testing getFinalEndpointURL method
func TestGetFinalEndpointURL(t *testing.T) {
	client := Client{}
	client.Secret = "test_secret"
	client.APIKey = "test_apikey"
	client.ProjectID = 1

	v := url.Values{}
	v.Set("test_key", "test_val")

	address, err := client.getFinalEndpointURL("getFile", v)
	found, _ := regexp.MatchString("https://platform\\.api\\.onesky\\.io/1/projects/1/translations\\?api_key=test_apikey&dev_hash=[a-z0-9]+&test_key=test_val&timestamp=[0-9]+", address)
	if !found {
		t.Errorf("getFinalEndpointURL() = %+v, %+v, want %+v,nil", address, err, "regexp(http://example\\.com/1/\\?api_key=test_apikey&dev_hash=[a-z0-9]+&test_key=test_val&timestamp=[0-9]+)")
	}
}

// TestGetURL is testing getURL method
func TestGetURL(t *testing.T) {
	client := Client{}
	client.Secret = "test_secret"
	client.APIKey = "test_apikey"
	client.ProjectID = 1

	want := "https://platform.api.onesky.io/1/projects/1/translations"
	address, err := client.getURL("getFile")
	if address != want {
		t.Errorf("getURL(%+v) = %+v, %+v, want %+v,nil", "getURL", address, err, want)
	}
}
