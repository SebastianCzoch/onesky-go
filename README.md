# onesky-go
[![Build Status](https://travis-ci.org/SebastianCzoch/onesky-go.svg?branch=master)](https://travis-ci.org/SebastianCzoch/onesky-go) [![Code Climate](https://codeclimate.com/github/SebastianCzoch/onesky-go/badges/gpa.svg)](https://codeclimate.com/github/SebastianCzoch/onesky-go) [![Coverage Status](https://coveralls.io/repos/SebastianCzoch/onesky-go/badge.svg?branch=feature%2Fcoverage&service=github)](https://coveralls.io/github/SebastianCzoch/onesky-go?branch=feature%2Fcoverage)  [![GoDoc](https://godoc.org/github.com/SebastianCzoch/onesky-go?status.svg)](https://godoc.org/github.com/SebastianCzoch/onesky-go)  [![License](https://img.shields.io/badge/licence-GNU%20v2-green.svg)](./LICENSE)



Go utils for working with [OneSky](http://www.oneskyapp.com/) translation service.

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
	err := onesky.UploadFile("messages.yml", "YAML", "en-US")
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

## Install

```
$ go get github.com/SebastianCzoch/onesky-go
````

or via [Godep](https://github.com/tools/godep)
```
$ godep get github.com/SebastianCzoch/onesky-go
```


## API

### (c *Client) DownloadFile(fileName, locale string) (string, error)
Downloads translation file from OneSky.

Returns file content via string.

### (c *Client) UploadFile(file, fileFormat, locale string) error
Upload translation file to OneSky.
* `file` should be a full path to file

### (c *Client) DeleteFile(fileName string) error
Permanently remove file from OneSky service (with translations)!

### (c *Client) ListFiles(page, perPage int) ([]FileData, error)
Get informations about files uploaded to OneSky

## Tests

```
$ go test ./...
````

## License

[GNU v2](./LICENSE)

## Support

Issues for this project should be reported on GitHub issues

Staff responsible for project:

* [Sebastian Czoch <sebastian@czoch.eu>](sebastian@czoch.eu)
