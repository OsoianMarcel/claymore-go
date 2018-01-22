# claymore-go
Go library used to get claymore stats

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/OsoianMarcel/claymore-go)

## Example

```go
package main

import (
	"github.com/OsoianMarcel/claymore-go"
	"net/http"
	"encoding/json"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
}

type ExtraResponse struct {
	HighestTemp claymore.TempAndFanReport `json:"highest_temp"`
}

type StatsResponse struct {
	Stats claymore.StatsModel `json:"stats"`
	Extra ExtraResponse       `json:"extra"`
}

func main() {
	cc := claymore.NewClient("localhost:3333")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		stats, err := cc.GetStats()

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}

		extraResp := ExtraResponse{}

		if ht, err := stats.GetHighestTemp(); err == nil {
			extraResp.HighestTemp = ht
		}

		statsResp := StatsResponse{stats, extraResp}

		json.NewEncoder(w).Encode(DataResponse{statsResp})
	})

	http.ListenAndServe(":8080", nil)
}
```



## Contribute

Contributions to the package are always welcome!

* Report any bugs or issues you find on the [issue tracker].
* You can grab the source code at the package's [Git repository].

## License

All contents of this package are licensed under the [MIT license].

[issue tracker]: https://github.com/OsoianMarcel/bnm-rates/issues
[Git repository]: https://github.com/OsoianMarcel/bnm-rates
[MIT license]: LICENSE
