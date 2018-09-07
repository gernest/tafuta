package tafuta

import (
	"io"
	"syscall/js"
)

func FetchValue() js.Value {
	return js.Global().Get("fetch")
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

type RequestCredentials uint

const (
	OmitCredential RequestCredentials = 1 << iota
	SameOrigin
	Include
)

func (c RequestCredentials) String() string {
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

type RequestMode uint

const (
	SameOriginMode RequestMode = 1 << iota
	NoCORSMode
	CORSMode
	NavigateMode
)

func (m RequestMode) String() string {
	switch m {
	case SameOriginMode:
		return "same-origin"
	case NoCORSMode:
		return "no-cors"
	case CORSMode:
		return "cors"
	case NavigateMode:
		return "navigate"
	default:
		return ""
	}
}

type RequestRedirect uint

const (
	FollowRedirect RequestRedirect = 1 << iota
	ErrorRedirect
	ManualRedirect
)

func (r RequestRedirect) String() string {
	switch r {
	case FollowRedirect:
		return "follow"
	case ErrorRedirect:
		return "error"
	case ManualRedirect:
		return "manual"
	default:
		return ""
	}
}

type RequestDestination uint

const (
	Audio RequestDestination = 1 << iota
	AudioWorklet
	Document
	Embed
	Font
	Image
	Manifest
	Object
	PaintWorklet
	Report
	Script
	ServiceWorker
	SharedWorker
	Style
	Track
	Video
	Worker
	XSLT
)

func (d RequestDestination) String() string {
	switch d {
	case Audio:
		return "audio"
	case AudioWorklet:
		return "audioworklet"
	case Document:
		return "document"
	case Embed:
		return "embed"
	case Font:
		return "font"
	case Image:
		return "image"
	case Manifest:
		return "manifest"
	case Object:
		return "object"
	case PaintWorklet:
		return "paintworklet"
	case Report:
		return "report"
	case Script:
		return "script"
	case ServiceWorker:
		return "serviceworker"
	case SharedWorker:
		return "sharedworker"
	case Style:
		return "style"
	case Track:
		return "track"
	case Video:
		return "video"
	case Worker:
		return "worker"
	case XSLT:
		return "xslt"
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
	Cache         RequestCache
	Credentials   RequestCredentials
	Destination   RequestDestination
	Header        *Header
	Integrity     string
	Method        string
	Mode          RequestMode
	Redirect      RequestRedirect
	Referer       string
	RefererPolicy string
	Body          io.Reader
}
