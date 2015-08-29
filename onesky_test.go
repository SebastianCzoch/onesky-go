// Package onesky tests
// Copyright (c) 2015 Sebastian Czoch <sebastian@czoch.eu>. All rights reserved.
// Use of this source code is governed by a GNU v2 license found in the LICENSE file.
package onesky

import (
	"net/url"
	"regexp"
	"testing"
)

var testEndpoints = map[string]apiEndpoint{
	"getFile":    apiEndpoint{"projects/%d/translations", "GET"},
	"postFile":   apiEndpoint{"projects/%d/files", "POST"},
	"deleteFile": apiEndpoint{"projects/%d/files", "DELETE"},
}

// TestFull is testing full method on apiEndpoint struct
func TestFull(t *testing.T) {
	client := Client{
		Secret:    "test_secret",
		APIKey:    "test_apikey",
		ProjectID: 1,
	}

	v := url.Values{}
	v.Set("test_key", "test_val")
	endpoint := testEndpoints["getFile"]
	want := "https://platform\\.api\\.onesky\\.io/1/projects/1/translations\\?api_key=test_apikey&dev_hash=[a-z0-9]+&test_key=test_val&timestamp=[0-9]+"
	address, err := endpoint.full(&client, v)
	if err != nil {
		t.Errorf("full() = %+v, %+v, want %+v,nil", address, err, want)
	}

	found, _ := regexp.MatchString(want, address)
	if !found {
		t.Errorf("full() = %+v, %+v, want %+v,nil", address, err, want)
	}
}

// TestGetEndpoint is testing getEndpoint function
func TestGetEndpoint(t *testing.T) {
	endpointName := "not_exist_endpoint"
	endpoint, err := getEndpoint(endpointName)
	if err == nil {
		t.Errorf("getEndpoint(%s) = %+v, %s, want %s,%s", endpointName, endpoint, err, nil, "endpoint not_exist_endpoint not found")
	}

	endpointName = "getFile"
	endpoint, err = getEndpoint(endpointName)
	if err != nil {
		t.Errorf("getEndpoint(%s) = %+v, %s, want %s,%s", endpointName, endpoint, err, nil, "endpoint not_exist_endpoint not found")
	}
}
