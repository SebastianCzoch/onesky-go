// Package onesky tests
package onesky

import (
	"net/url"
	"regexp"
	"testing"
)

// TestGetURLForEndpoint is testing GetURLForEndpoint method
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

// TestGetFinalEndpointURL is testing GetFinalEndpointURL method
func TestGetFinalEndpointURL(t *testing.T) {
	client := Client{}
	client.Secret = "test_secret"
	client.APIKey = "test_apikey"

	v := url.Values{}
	v.Set("test_key", "test_val")

	address, err := client.getFinalEndpointURL("http://example.com/1/", v)
	found, _ := regexp.MatchString("http://example\\.com/1/\\?api_key=test_apikey&dev_hash=[a-z0-9]+&test_key=test_val&timestamp=[0-9]+", address)
	if !found {
		t.Errorf("getFinalEndpointURL() = %+v, %+v, want %+v,nil", address, err, "regexp(http://example\\.com/1/\\?api_key=test_apikey&dev_hash=[a-z0-9]+&test_key=test_val&timestamp=[0-9]+)")
	}
}
