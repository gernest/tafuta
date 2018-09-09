package tafuta

import (
	"io"
	"io/ioutil"
	"net/url"
	"syscall/js"
)

var global = js.Global()

// FetchValue returns global javascript fetch function.
func FetchValue() js.Value {
	return global.Get("fetch")
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

// RequestCache defines modes for cache. This defines how the request will
// interact with browser HTTP cache.
//
// Example
// Download a resource with cache busting, to bypass the cache
// completely.
//
// 	client := tafuta.NewClient()
// 	res, err := client.Do(&tafuta.Request{
// 		URL:   "some.json",
// 		Cache: tafuta.NoStore,
// 	})
// 	if err != nil {
// 		// handle error
// 	}
//
// Download a resource with cache busting, but update the HTTP
// cache with the downloaded resource.
//
// 	client := tafuta.NewClient()
// 	res, err := client.Do(&tafuta.Request{
// 		URL:   "some.json",
// 		Cache: tafuta.Reload,
// 	})
// 	if err != nil {
// 		// handle error
// 	}
//
// Download a resource with cache busting when dealing with a
// properly configured server that will send the correct ETag
// and Date headers and properly handle If-Modified-Since and
// If-None-Match request headers, therefore we can rely on the
// validation to guarantee a fresh response.
//
// 	client := tafuta.NewClient()
// 	res, err := client.Do(&tafuta.Request{
// 		URL:   "some.json",
// 		Cache: tafuta.NoCache,
// 	})
// 	if err != nil {
// 		// handle error
// 	}
//
// Download a resource with economics in mind!  Prefer a cached
// albeit stale response to conserve as much bandwidth as possible.
//
// 	client := tafuta.NewClient()
// 	res, err := client.Do(&tafuta.Request{
// 		URL:   "some.json",
// 		Cache: tafuta.ForceCache,
// 	})
// 	if err != nil {
// 		// handle error
// 	}
//
type RequestCache uint

const (
	// DefaultCache in this mode the browser looks for a matching request in its
	// HTTP cache.
	//
	// 	If there is a match and it is fresh, it will be returned from the cache.
	//
	// 	If there is a match but it is stale, the browser will make a conditional
	// 	request to the remote server. If the server indicates that the resource has
	// 	not changed, it will be returned from the cache. Otherwise the resource will
	// 	be downloaded from the server and the cache will be updated.
	//
	// 	If there is no match, the browser will make a normal request, and will
	// 	update the cache with the downloaded resource.
	DefaultCache RequestCache = 1 << iota

	// NoStore The browser fetches the resource from the remote server without
	// first looking in the cache, and will not update the cache with the
	// downloaded resource.
	NoStore

	// Reload The browser fetches the resource from the remote server without first
	// looking in the cache, but then will update the cache with the downloaded
	// resource.
	Reload

	// NoCache The browser looks for a matching request in its HTTP cache.
	//
	// 	If there is a match, fresh or stale, the browser will make a conditional
	// 	request to the remote server. If the server indicates that the resource has
	// 	not changed, it will be returned from the cache. Otherwise the resource will
	// 	be downloaded from the server and the cache will be updated.
	//
	// 	If there is no match, the browser will make a normal request, and will
	// 	update the cache with the downloaded resource.
	NoCache

	// ForceCache The browser looks for a matching request in its HTTP cache.
	//
	// 	If there is a match, fresh or stale, it will be returned from the cache.
	//
	// 	If there is no match, the browser will make a normal request, and will
	// 	update the cache with the downloaded resource.
	ForceCache

	// OnlyIfCached The browser looks for a matching request in its HTTP cache.
	//
	// 	If there is a match, fresh or stale, if will be returned from the cache.
	//
	// 	If there is no match, the browser will respond with a 504 Gateway timeout
	// 	status.
	//
	// The "only-if-cached" mode can only be used if the request's mode is
	// "same-origin". Cached redirects will be followed if the request's redirect
	// property is "follow" and the redirects do not violate the "same-origin"
	// mode.
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
	URL           string
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
type ResponseType uint

const (
	BasicResponse ResponseType = 1 << iota
	CorsResponse
	ErrorResponse
	OpaqueResponse
	OpaqueRedirectResponse
)

var respTypMap = map[ResponseType]string{
	BasicResponse:          "basic",
	CorsResponse:           "cors",
	ErrorResponse:          "error",
	OpaqueResponse:         "opaque",
	OpaqueRedirectResponse: "opaqueredirect",
}

func (r ResponseType) String() string {
	return respTypMap[r]
}

type Response struct {
	Headers    *Header
	Ok         bool
	Redirected bool
	Status     int
	StatusText string
	Type       ResponseType
	URL        *url.URL
	Body       io.ReadCloser
}

type Client struct {
	value js.Value
}

func NewClient() *Client {
	return &Client{value: FetchValue()}
}

func (c *Client) Do(req *Request) (res *Response, err error) {
	var resources resourceList
	defer func() {
		if resources != nil {
			resources.free()
		}
	}()
	args := []interface{}{req.URL}
	opts := make(map[string]interface{})
	if req.Method != "" {
		opts["method"] = req.Method
	}
	if req.Header != nil {
		opts["headers"] = req.Header.Value()
	}
	if mode := req.Mode.String(); mode != "" {
		opts["mode"] = mode
	}
	if creds := req.Credentials.String(); creds != "" {
		opts["credentials"] = creds
	}
	if cache := req.Cache.String(); cache != "" {
		opts["cache"] = cache
	}
	done := make(chan struct{})
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		a := js.TypedArrayOf(b)
		resources = append(resources, a)
		opts["body"] = a
	}
	if len(opts) > 0 {
		args = append(args, opts)
	}
	request := js.Global().Get("Request").New(args...)
	responseCallback := js.NewCallback(func(v []js.Value) {
		console(v[0])
		done <- struct{}{}
	})
	r := c.value.Invoke(request)
	r.Call("then", responseCallback)
	<-done
	return
}

type resource interface {
	Release()
}

// Simple helper for releasing multiple resounces.
type resourceList []resource

func (r resourceList) free() {
	for _, v := range r {
		v.Release()
	}
}

func console(v js.Value) {
	global.Get("console").Call("log", v)
}
