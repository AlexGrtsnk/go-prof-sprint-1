package main

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()

	client.
		SetRetryCount(1).
		SetRetryWaitTime(30 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(`aaabbb`).
		Post("http://localhost:8080/")

	if err != nil {
		panic(err)
	}
	qwe := string(resp.String())
	fmt.Println(qwe)
	resp_1, err_1 := client.R().
		SetHeader("Content-Type", "text/plain").
		Get("http://localhost:8080/" + qwe[22:])
	fmt.Println(resp_1)
	if err_1 != nil {
		panic(err)
	}
}
