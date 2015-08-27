# onesky-go [![Build Status](https://travis-ci.org/SebastianCzoch/onesky-go.svg?branch=master)](https://travis-ci.org/SebastianCzoch/onesky-go) [![GoDoc](https://godoc.org/github.com/SebastianCzoch/onesky-go?status.svg)](https://godoc.org/github.com/SebastianCzoch/onesky-go)
Go utils for working with [OneSky](http://www.oneskyapp.com/) translation service.

## Example

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

## Install

```
$ go get github.com/SebastianCzoch/onesky-go
````

## API

### (c *Client) DownloadFile(fileName, locale string) (string, error)
Downloads translation file from OneSky.

Returns file content via string.

## Tests

```
$ go test
````

## License

[GNU v2](./LICENSE)
