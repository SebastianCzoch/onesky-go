# onesky-go
[![Build Status](https://travis-ci.org/SebastianCzoch/onesky-go.svg?branch=master)](https://travis-ci.org/SebastianCzoch/onesky-go) [![Code Climate](https://codeclimate.com/github/SebastianCzoch/onesky-go/badges/gpa.svg)](https://codeclimate.com/github/SebastianCzoch/onesky-go) [![Coverage Status](https://coveralls.io/repos/SebastianCzoch/onesky-go/badge.svg?branch=feature%2Fcoverage&service=github)](https://coveralls.io/github/SebastianCzoch/onesky-go?branch=feature%2Fcoverage)  [![GoDoc](https://godoc.org/github.com/SebastianCzoch/onesky-go?status.svg)](https://godoc.org/github.com/SebastianCzoch/onesky-go)  [![License](https://img.shields.io/badge/licence-GNU%20v2-green.svg)](./LICENSE)



Go utils for working with [OneSky](http://www.oneskyapp.com/) translation service.

## Install

```
$ go get github.com/SebastianCzoch/onesky-go
````

or via [Godep](https://github.com/tools/godep)
```
$ godep get github.com/SebastianCzoch/onesky-go
```

## Examples
### Example 1 - Download file

```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	fmt.Println(onesky.DownloadFile("filename", "locale"))
}
```

### Example 2 - Upload file

```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	_, err := onesky.UploadFile("messages.yml", "YAML", "en-US", true)
	if err != nil {
		fmt.Println("Can not upload file")
	}
}
```

### Example 3 - Delete file

```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	err := onesky.DeleteFile("messages.yml")
	if err != nil {
		fmt.Println("Can not delete file")
	}
}
```

### Example 4 - Get informations about uploaded files
```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	list, err := onesky.ListFiles(1, 100)
	if err != nil {
		fmt.Println("Can not download list of uploaded files")
	}
	fmt.Println(list)
}
```

### Example 5 - List import tasks
```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	list, err := onesky.ImportTasks(map[string]interface{}{
		"per_page": 50,
		"status": "completed", // all, completed, in-progress, failed
	})
	if err != nil {
		fmt.Println("Can not download list of import tasks")
	}
	fmt.Println(list)
}
```

### Example 6 - Show an import task
```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	task, err := onesky.ImportTask(773572) // import id
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(task)
}
```

### Example 7 - Get a project's translations status
```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	status, err := onesky.GetTranslationsStatus("string.po", "ja-JP")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(status)
}
```

### Example 8 - Get languages
```
package main

import (
	"fmt"
	"github.com/SebastianCzoch/onesky-go"
)

func main() {
	onesky := onesky.Client{APIKey: "abcdef", Secret: "abcdef", ProjectID: 1}
	languages, err := onesky.GetLanguages()
	if err != nil {
		fmt.Println("Can not get languages")
	}
	fmt.Println(languages)
}
```


## API

### (c *Client) DownloadFile(fileName, locale string) (string, error)
Downloads translation file from OneSky.

Returns file content via string.

### (c *Client) UploadFile(file, fileFormat, locale string, keepStrings bool) (UploadData, error)
Upload translation file to OneSky.
* `file` should be a full path to file

### (c *Client) DeleteFile(fileName string) error
Permanently remove file from OneSky service (with translations)!

### (c *Client) ListFiles(page, perPage int) ([]FileData, error)
Get informations about files uploaded to OneSky

### (c *Client) ImportTasks(params) ([]TaskData, error)
List import tasks. (Default params: `{"page": 1, "per_page": 50, "status": "all"}`)

### (c *Client) ImportTask(importID) (TaskData, error)
Show an import task.

### (c *Client) GetTranslationsStatus(fileName, locale string) (TranslationsStatus, error)
Shows a project's translations status.

### (c *Client) GetLanguages() ([]Language, error)
Get informations about available languages in project

## Tests

```
$ go test ./...
````

## License

[GNU v2](./LICENSE)

## Support

Issues for this project should be reported on GitHub issues
