package tafuta

import (
	"syscall/js"
)

func FatchValue() js.Value {
	return js.Global().Get("fetch")
}
