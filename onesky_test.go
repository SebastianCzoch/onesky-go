package onesky

import (
	"testing"
	"net/url"
	"regexp"
)

func TestGetUrlForEndpoint(t *testing.T) {
	options := Options{}
	options.Secret = "test_secret"
	options.ProjectID = 1
	
	
	url, err := options.getUrlForEndpoint("not_exits_endpoint");
	if err == nil {
		t.Errorf("getUrlForEndpoint() = %+v, %+v, want %+v", url, err, "error")
	}
	
	want := API_ADDRESS + "/" + API_VERSION + "/" + "projects/1/translations"
	url, err = options.getUrlForEndpoint("getFile");
	if url != want {
		t.Errorf("getUrlForEndpoint() = %+v, %+v, want %+v", url, err, want)
	}
}

func TestGetFinalEndpointUrl(t *testing.T) {
	options := Options{}
	options.Secret = "test_secret"
	options.ProjectID = 1
	options.ApiKey = "test_apikey"
	
	address, err := options.getFinalEndpointUrl("not_exits_endpoint", url.Values{});
	if err == nil {
		t.Errorf("getFinalEndpointUrl() = %+v, %+v, want %+v", address, err, "error")
	}

	v := url.Values{}
	v.Set("test_key", "test_val")
	
	address, err = options.getFinalEndpointUrl("http://example.com/%d/", v);
	found, _ := regexp.MatchString("http://example\\.com/1/\\?api_key=test_apikey&dev_hash=[a-z0-9]+&test_key=test_val&timestamp=[0-9]+", address)
	if !found {
		t.Errorf("getFinalEndpointUrl() = %+v, %+v, want %+v,nil", address, err, "regexp(http://example\\.com/1/\\?api_key=test_apikey&dev_hash=[a-z0-9]+&test_key=test_val&timestamp=[0-9]+)")
	}
}
