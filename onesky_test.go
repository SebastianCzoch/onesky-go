// Package onesky tests
// Copyright (c) 2015 Sebastian Czoch <sebastian@czoch.eu>. All rights reserved.
// Use of this source code is governed by a GNU v2 license found in the LICENSE file.
package onesky

import (
	"fmt"
	"net/url"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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
		t.Errorf("getEndpoint(%s) = %+v, %s, want %v,%s", endpointName, endpoint, err, nil, "endpoint not_exist_endpoint not found")
	}

	endpointName = "getFile"
	endpoint, err = getEndpoint(endpointName)
	if err != nil {
		t.Errorf("getEndpoint(%s) = %+v, %s, want %v,%s", endpointName, endpoint, err, nil, "endpoint not_exist_endpoint not found")
	}
}

func TestDeleteFileWithSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	err := client.DeleteFile("test.yml")
	assert.Nil(t, err)
}

func TestDeleteFileWithFailure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	err := client.DeleteFile("test.yml")
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}

func TestListFilesWithFailure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	_, err := client.ListFiles(1, 1)
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}

func TestListFilesWithSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"meta":{"status":200,"record_count":16},"data":[{"name":"strings.po","string_count":236,"last_import":{"id":123,"status":"in-progress"},"uploaded_at":"2013-10-07T15:27:10+0000","uploaded_at_timestamp":1381159630},{"name":"en.yml","string_count":335,"last_import":{"id":109,"status":"completed"},"uploaded_at":"2013-10-05T12:36:52+0000","uploaded_at_timestamp":1380976612},{"name":"Manuallyinput","string_count":285}]}`))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	res, err := client.ListFiles(1, 1)
	assert.Nil(t, err)

	assert.Equal(t,
		[]FileData{
			FileData{
				Name:        "strings.po",
				StringCount: 236,
				LastImport: LastImport{
					ID:     123,
					Status: "in-progress",
				},
				UpoladedAt:          "2013-10-07T15:27:10+0000",
				UpoladedAtTimestamp: 1381159630,
			},
			FileData{
				Name:        "en.yml",
				StringCount: 335,
				LastImport: LastImport{
					ID:     109,
					Status: "completed",
				},
				UpoladedAt:          "2013-10-05T12:36:52+0000",
				UpoladedAtTimestamp: 1380976612,
			},
			FileData{
				Name:        "Manuallyinput",
				StringCount: 285,
				LastImport: LastImport{
					ID:     0,
					Status: "",
				},
				UpoladedAt:          "",
				UpoladedAtTimestamp: 0,
			},
		}, res)
}
