// Package onesky tests
// Copyright (c) 2015 Sebastian Czoch <sebastian@czoch.eu>. All rights reserved.
// Use of this source code is governed by a GNU v2 license found in the LICENSE file.
package onesky

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testEndpoints = map[string]apiEndpoint{
	"getFile":               apiEndpoint{"projects/%d/translations", "GET"},
	"postFile":              apiEndpoint{"projects/%d/files", "POST"},
	"deleteFile":            apiEndpoint{"projects/%d/files", "DELETE"},
	"listFiles":             apiEndpoint{"projects/%d/files", "GET"},
	"importTasks":           apiEndpoint{"projects/%d/import-tasks", "GET"},
	"importTask":            apiEndpoint{"projects/%d/import-tasks/%d", "GET"},
	"getTranslationsStatus": apiEndpoint{"projects/%d/translations/status", "GET"},
	"getLanguages":          apiEndpoint{"projects/%d/languages", "GET"},
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

func TestDownloadFileWithFailure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	_, err := client.DownloadFile("test.yml", "en_US")
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}

func TestDownloadFileWithSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `test: translatedTest`))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	res, err := client.DownloadFile("test.yml", "en_US")
	fmt.Println(res)
	assert.Nil(t, err)

	assert.Equal(t, `test: translatedTest`, res)
}

func TestUploadFileWithSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(201, `{"meta":{"status":201},"data":{"name":"string.po","format":"GNU_PO","language":{"code":"en-US","english_name":"English (United States)","local_name":"English (United States)","locale":"en","region":"US"},"import":{"id":154,"created_at":"2013-10-07T15:27:10+0000","created_at_timestamp":1381159630}}}`))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	tmpdir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpdir)

	filename := path.Join(tmpdir, "string.po")
	ioutil.WriteFile(filename, []byte("test"), 0666)

	res, err := client.UploadFile(filename, "GNU_PO", "en_US", true)
	assert.Nil(t, err)

	assert.Equal(t, UploadData{
		Name:   "string.po",
		Format: "GNU_PO",
		Language: Language{
			Code:        "en-US",
			EnglishName: "English (United States)",
			LocalName:   "English (United States)",
			Locale:      "en",
			Region:      "US",
		},
		Import: TaskData{
			ID:                  154,
			OriginalID:          154.0,
			CreateddAt:          "2013-10-07T15:27:10+0000",
			CreateddAtTimestamp: 1381159630,
		},
	}, res)
}

func TestUploadFileWithFailure(t *testing.T) {
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	tmpdir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpdir)

	filename := path.Join(tmpdir, "not_found")
	_, err = client.UploadFile(filename, "GNU_PO", "en_US", true)
	assert.NotNil(t, err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))

	ioutil.WriteFile(filename, []byte("test"), 0666)
	_, err = client.UploadFile(filename, "GNU_PO", "en_US", true)
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}

func TestImportTasksWithFailure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	_, err := client.ImportTasks(map[string]interface{}{"page": 1, "per_page": 50, "status": "all"})
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}

func TestImportTasksWithSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"meta":{"status":200},"data":[{"id":"177","file":{"name":"string2.po"},"status":"in-progress","created_at":"2013-10-07T15:25:00+0000","created_at_timestamp":1381159500},{"id":"176","file":{"name":"string.po"},"status":"in-progress","created_at":"2013-10-07T15:27:10+0000","created_at_timestamp":1381159630}]}`))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	res, err := client.ImportTasks(map[string]interface{}{"page": 1, "per_page": 50, "status": "in-progress"})
	assert.Nil(t, err)

	assert.Equal(t,
		[]TaskData{
			TaskData{
				ID:         177,
				OriginalID: "177",
				File: TaskFile{
					Name: "string2.po",
				},
				Status:              "in-progress",
				CreateddAt:          "2013-10-07T15:25:00+0000",
				CreateddAtTimestamp: 1381159500,
			},
			TaskData{
				ID:         176,
				OriginalID: "176",
				File: TaskFile{
					Name: "string.po",
				},
				Status:              "in-progress",
				CreateddAt:          "2013-10-07T15:27:10+0000",
				CreateddAtTimestamp: 1381159630,
			},
		}, res)
}

func TestImportTaskWithFailure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	_, err := client.ImportTask(1)
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}
func TestImportTaskWithSuccess(t *testing.T) {
	httpmock.Activate()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"meta":{"status":200},"data":{"id":176,"file":{"name":"string.po","format":"GNU_PO","locale":{"code":"en-US","english_name":"English (United States)","local_name":"English (United States)","locale":"en","region":"US"}},"string_count":236,"word_count":1260,"status":"in-progress","created_at":"2013-10-07T15:27:10+0000","created_at_timestamp":1381159630}}`))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	res, err := client.ImportTask(1)
	assert.Nil(t, err)

	assert.Equal(t,
		TaskData{
			ID:         176,
			OriginalID: float64(176),
			File: TaskFile{
				Name:   "string.po",
				Format: "GNU_PO",
				Locale: Language{
					Code:        "en-US",
					EnglishName: "English (United States)",
					LocalName:   "English (United States)",
					Locale:      "en",
					Region:      "US",
				},
			},
			StringCount:         236,
			WordCount:           1260,
			Status:              "in-progress",
			CreateddAt:          "2013-10-07T15:27:10+0000",
			CreateddAtTimestamp: 1381159630,
		}, res)
}

func TestGetTranslationsStatusWithFailure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	_, err := client.GetTranslationsStatus("string.po", "ja-JP")
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}
func TestGetTranslationsStatusWithSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"meta":{"status":200},"data":{"file_name":"string.po","locale":{"code":"ja-JP","english_name":"Japanese","local_name":"\u65e5\u672c\u8a9e","locale":"ja","region":"JP"},"progress":"92%","string_count":1359,"word_count":3956}}`))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	res, err := client.GetTranslationsStatus("string.po", "ja-JP")
	assert.Nil(t, err)

	assert.Equal(t,
		TranslationsStatus{
			FileName: "string.po",
			Locale: Language{
				Code:         "ja-JP",
				EnglishName:  "Japanese",
				LocalName:    "日本語",
				CustomLocale: "",
				Locale:       "ja",
				Region:       "JP",
			},
			Progress:    "92%",
			StringCount: 1359,
			WordCount:   3956,
		}, res)
}

func TestGetLanguagesWithFailure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(500, ""))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	_, err := client.GetLanguages()
	assert.Equal(t, err, fmt.Errorf("bad status: %d", 500))
}
func TestGetLanguagesWithSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"meta":{"status":200,"record_count":17},"data":[{"code":"it","english_name":"Italian","local_name":"Italiano\u0000","custom_locale":"","Locale":"it","region":"","translation_progress":"0.0"},{"code":"de","english_name":"German","local_name":"Deutsch\u0000","custom_locale":"","locale":"de","region":"","translation_progress":"0.0"},{"code":"fr","english_name":"French","local_name":"Français\u0000","custom_locale":"","locale":"fr","region":"","translation_progress":"0.0"}]}`))
	client := Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}

	res, err := client.GetLanguages()
	assert.Nil(t, err)

	assert.Equal(t,
		[]Language{
			Language{
				Code:                "it",
				EnglishName:         "Italian",
				LocalName:           "Italiano\u0000",
				CustomLocale:        "",
				Locale:              "it",
				Region:              "",
				TranslationProgress: "0.0",
			},
			Language{
				Code:                "de",
				EnglishName:         "German",
				LocalName:           "Deutsch\u0000",
				CustomLocale:        "",
				Locale:              "de",
				Region:              "",
				TranslationProgress: "0.0",
			},
			Language{
				Code:                "fr",
				EnglishName:         "French",
				LocalName:           "Français\u0000",
				CustomLocale:        "",
				Locale:              "fr",
				Region:              "",
				TranslationProgress: "0.0",
			},
		}, res)
}
