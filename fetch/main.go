package main

import (
	"bytes"
	"encoding/json"

	"github.com/gernest/tafuta"
)

func main() {
	v := tafuta.NewClient()
	b, _ := json.Marshal([]int{12345})
	v.Do(&tafuta.Request{
		URL:    "/nothing",
		Method: "POST",
		Body:   bytes.NewReader(b),
	})
}
