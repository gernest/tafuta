package main

import (
	"fmt"

	"github.com/gernest/tafuta"
)

func main() {
	v := tafuta.FetchValue()
	fmt.Println(v.Type())
}
