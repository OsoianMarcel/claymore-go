# claymore-go
Go library used to get claymore stats in human readable model

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/OsoianMarcel/claymore-go/blob/master/LICENSE)

## Example of simple web server using this library

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

## Server output
```json
{
	"data": {
		"stats": {
			"miner_version": "10.5 - ETH",
			"running_minutes": 1409,
			"eth_report": {
				"total_mhs": 80941,
				"shares": 1640,
				"rejected_shares": 1,
				"invalid_shares": 37,
				"pool_switches": 1,
				"mhs_per_gpu": [
					{
						"mhs": 13344,
						"gpu": 0
					},
					{
						"mhs": 13352,
						"gpu": 1
					},
					{
						"mhs": 13398,
						"gpu": 2
					},
					{
						"mhs": 14151,
						"gpu": 3
					},
					{
						"mhs": 13352,
						"gpu": 4
					},
					{
						"mhs": 13342,
						"gpu": 5
					}
				]
			},
			"alt_report": {
				"total_mhs": 0,
				"shares": 0,
				"rejected_shares": 0,
				"invalid_shares": 0,
				"pool_switches": 0,
				"mhs_per_gpu": [
					{
						"mhs": 0,
						"gpu": 0
					},
					{
						"mhs": 0,
						"gpu": 1
					},
					{
						"mhs": 0,
						"gpu": 2
					},
					{
						"mhs": 0,
						"gpu": 3
					},
					{
						"mhs": 0,
						"gpu": 4
					},
					{
						"mhs": 0,
						"gpu": 5
					}
				]
			},
			"temp_and_fan_reports": [
				{
					"temp": 66,
					"fan": 40,
					"gpu": 0
				},
				{
					"temp": 66,
					"fan": 40,
					"gpu": 1
				},
				{
					"temp": 67,
					"fan": 41,
					"gpu": 2
				},
				{
					"temp": 67,
					"fan": 40,
					"gpu": 3
				},
				{
					"temp": 67,
					"fan": 40,
					"gpu": 4
				},
				{
					"temp": 63,
					"fan": 80,
					"gpu": 5
				}
			],
			"pools": [
				"eu1.ethermine.org:4444"
			]
		},
		"extra": {
			"highest_temp": {
				"temp": 67,
				"fan": 41,
				"gpu": 2
			}
		}
	}
}
```

## Contribute

Contributions to the package are always welcome!

* Report any bugs or issues you find on the [issue tracker].
* You can grab the source code at the package's [Git repository].

## Donation

```
ETH: 0x58aaa089338901fcf5fb59342c97c17fa3dd1229
BTC: 13s6V1jzs84qijdcDTSXsEto7KZEE8cwZz
```

## License

All contents of this package are licensed under the [MIT license].

[issue tracker]: https://github.com/OsoianMarcel/claymore-go/issues
[Git repository]: https://github.com/OsoianMarcel/claymore-go
[MIT license]: LICENSE
