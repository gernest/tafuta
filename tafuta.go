package tafuta

import (
	"net/http"
	"syscall/js"
)

func FetchValue() js.Value {
	return js.Global().Get("fetch")
}

type Transport struct{}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, nil
}

// Header is a wrapper of javascript Headers
type Header struct {
	value js.Value
}

// NewHeader returns an instance of Header.
func NewHeader() *Header {
	v := js.Global().Get("Headers").New()
	return &Header{value: v}
}

// Add appends a new value onto an existing header inside a Headers object, or
// adds the header if it does not already exist.
func (h *Header) Add(key, value string) {
	h.value.Call("append", key, value)
}

// Del deletes a header from a Headers objectDel
func (h *Header) Del(key string) {
	h.value.Call("delete", key)
}

// Get Returns a string of all the values of a header within a
// Headers object with a given name.
func (h *Header) Get(key string) string {
	v := h.value.Get(key)
	if v.Type() == js.TypeString {
		return v.String()
	}
	return ""
}

// Set sets a new value for an existing header inside a Headers object, or adds
// the header if it does not already exist.
func (h *Header) Set(key, value string) {
	h.value.Call("set", key, value)
}

// Value returns js value of the header.
func (h *Header) Value() js.Value {
	return h.value
}
