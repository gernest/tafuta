package main

import (
	"fmt"
	"syscall/js"

	"github.com/gernest/tafuta"
)

func main() {
	v := tafuta.FetchValue()
	fmt.Println(v.Type())
	h := tafuta.NewHeader()
	h.Set("Content-Type", "text/xml")
	js.Global().Set("someHead", h.Value())
	js.Global().Get("console").Call("log", h.Value())
}
