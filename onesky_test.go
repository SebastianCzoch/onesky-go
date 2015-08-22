package onesky

import (
	"testing"
	// "crypto/md5"
	// "time"
	// "fmt"
)

func TestGetUrlForEndpoint(t *testing.T) {
	options := Options{}
	options.Secret = "test_secret"
	
	url, err := options.getUrlForEndpoint("not_exits_endpoint");
	if err == nil {
		t.Errorf("getUrlForEndpoint() = %+v, %+v, want %+v", url, err, "error")
	}
	
	want := API_ADDRESS + "/" + API_VERSION + "/" + "projects/translations"
	url, err = options.getUrlForEndpoint("getFile");
	if url != want {
		t.Errorf("getUrlForEndpoint() = %+v, %+v, want %+v", url, err, want)
	}
}


