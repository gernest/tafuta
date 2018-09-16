// Package tafuta provide wrapper around the fetch API which allows easy and
// seamless of making http calls in a go wasm project.
//
// This abstracts away the low lever details necessary for interacting with
// browser's dom. In the process offering familiar API.
//
// Note that all calls in this package are blocking. So use goroutines and
// channels for synchronization. The choice of making blocking API is to make it
// easily integrated with existing go libraries/API's.
//
// This is how you send a GET request with header
// 	client := tafuta.NewClient()
// 	h := tafuta.NewHeader()
// 	h.Set("Content-Type", "image/jpeg")
// 	res, err := client.Do(&tafuta.Request{
// 		Method: "GET",
// 		URL:    "flowers.jpg",
// 		Header: h,
// 	})
// 	if err != nil {
// 		// handle error
// 	}
// 	println(res.Status)
package tafuta

import (
	"fmt"
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
// Headers object with a given key.
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
	URL         string
	Cache       RequestCache
	Credentials RequestCredentials
	Destination RequestDestination
	Header      *Header
	Method      string
	Mode        RequestMode
	Redirect    RequestRedirect

	// Can either be no-referrer, client, or a URL. The default is client.
	Referer string

	// Contains the subresource integrity value of the request (e.g.,
	// sha256-BpfBw7ivV8q2jLiT13fxDYAe2tJllusRSZ273h2nFSE=).
	Integrity string
	Body      io.Reader
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

var reverseRespTypMap map[string]ResponseType

func init() {
	reverseRespTypMap = make(map[string]ResponseType)
	for k, v := range respTypMap {
		reverseRespTypMap[v] = k
	}
}

func (r ResponseType) String() string {
	return respTypMap[r]
}

// Response represents a response to a fetch API Request.
type Response struct {
	Headers    *Header
	Ok         bool
	Redirected bool
	Status     int
	StatusText string
	Type       ResponseType
	URL        *url.URL
	Body       io.ReadCloser
	value      js.Value
}

// NewResponse creates *Response struct from Response js object.
func NewResponse(v js.Value) (*Response, error) {
	res := &Response{}
	res.Headers = &Header{value: v.Get("headers")}
	res.Ok = v.Get("ok").Bool()
	res.Status = v.Get("status").Int()
	res.StatusText = v.Get("statusText").String()
	resType := v.Get("type").String()
	res.Type = reverseRespTypMap[resType]
	u, err := url.Parse(v.Get("url").String())
	if err != nil {
		return nil, err
	}
	res.URL = u
	res.value = v
	return res, nil
}

// Text returns Response body contents as a string. This is a blocking call,
// please use this in a separate goroutines to avoid blocking execution of other
// code.
func (r *Response) Text() (res string) {
	done := make(chan struct{})
	responseCallback := js.NewCallback(func(v []js.Value) {
		res = v[0].String()
		done <- struct{}{}
	})
	defer responseCallback.Release()
	r.value.Call("text").Call("then", responseCallback)
	<-done
	return
}

type Client struct {
	value js.Value
}

func NewClient() *Client {
	return &Client{value: FetchValue()}
}

// Do sends request using fetch AP. This method is blocking, to avoid
// deadlocking your app please call this inside a goroutine.
func (c *Client) Do(req *Request) (res *Response, err error) {
	var resources resourceList
	defer func() {
		if resources != nil {
			resources.free()
		}
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
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
	if redirect := req.Redirect.String(); redirect != "" {
		opts["redirect"] = redirect
	}
	if req.Referer != "" {
		opts["referrer"] = req.Referer
	}
	if req.Integrity != "" {
		opts["integrity"] = req.Integrity
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
		res, err = NewResponse(v[0])
		done <- struct{}{}
	})
	r := c.value.Invoke(request)
	resources = append(resources, responseCallback)
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
