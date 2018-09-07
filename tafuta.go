package tafuta

import (
	"syscall/js"
)

func FetchValue() js.Value {
	return js.Global().Get("fetch")
}
