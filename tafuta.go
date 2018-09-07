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

type RequestCache uint

const (
	DefaultCache RequestCache = 1 << iota
	NoStore
	Reload
	NoCache
	ForceCache
	OnlyIfCached
)

func (c RequestCache) String() string {
	switch c {
	case DefaultCache:
		return "default"
	case NoStore:
		return "no-store"
	case Reload:
		return "reload"
	case NoCache:
		return "no-cache"
	case ForceCache:
		return "force-cache"
	case OnlyIfCached:
		return "only-if-cached"
	default:
		return ""
	}
}

type Credentials uint

const (
	OmitCredential Credentials = 1 << iota
	SameOrigin
	Include
)

type RequestMode string

const (
	SameOriginMode RequestMode = "same-origin"
	NoCORSMode     RequestMode = "no-cors"
	CORSMode       RequestMode = "cors"
	NavigateMode   RequestMode = "navigate"
)

type RequestRedirect string

const (
	FollowRedirect RequestRedirect = "follow"
	ErrorRedirect  RequestRedirect = "error"
	ManualRedirect RequestRedirect = "manual"
)

func (c Credentials) String() string {
	switch c {
	case OmitCredential:
		return "omit"
	case SameOrigin:
		return "same-origin"
	case Include:
		return "include"
	default:
		return ""
	}
}

type Iterator struct {
	js.Value
}

func (i *Iterator) Next() (done bool, value js.Value) {
	v := i.Get("next")
	done = v.Get("done").Bool()
	value = v.Get("value")
	return
}

// Range iterates over items and calling fn for every value. This will stop when
// fn returns false.
func (i *Iterator) Range(fn func(js.Value) bool) {
	for {
		done, value := i.Next()
		if done {
			return
		}
		if !fn(value) {
			return
		}
	}
}

type Request struct {
	Cache       RequestCache
	Credentials Credentials
	Method      string
	Mode        RequestMode
	Redirect    RequestRedirect
}
