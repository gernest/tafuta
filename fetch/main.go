package main

import (
	"fmt"

	"github.com/gernest/tafuta"
)

func main() {
	client := tafuta.NewClient()
	h := tafuta.NewHeader()
	h.Set("Content-Type", "image/jpeg")
	res, err := client.Do(&tafuta.Request{
		Method: "GET",
		URL:    "flowers.jpg",
		Header: h,
	})
	if err != nil {
		// handle error
	}
	fmt.Println(res.Text())
}
