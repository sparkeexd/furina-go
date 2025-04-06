package network

import (
	"fmt"
	"maps"
)

const (
	// Content encodings.
	encodingGzip    = "gzip"
	encodingDeflate = "deflate"
	encodingBr      = "br"
)

// HTTP request wrapper for holding parameters required for HoYoLab endpoints.
// Used by Handler to be parsed into a HTTP request.
type Request struct {
	endpoint string
	method   string
	body     map[string]any
	params   map[string]string
	headers  map[string]string
}

// Request builder that defines a request structure using the builder pattern.
type RequestBuilder struct {
	request Request
}

// Constructor.
func NewRequest(endpoint string, method string) *RequestBuilder {
	builder := &RequestBuilder{
		request: Request{
			endpoint: endpoint,
			method:   method,
		},
	}

	builder.addDefaultHeaders()
	return builder
}

// Return built request instance.
func (builder *RequestBuilder) Build() Request {
	return Request{
		endpoint: builder.request.endpoint,
		method:   builder.request.method,
		body:     maps.Clone(builder.request.body),
		params:   maps.Clone(builder.request.params),
		headers:  maps.Clone(builder.request.headers),
	}
}

// Add body.
func (builder *RequestBuilder) AddBody(key string, value any) *RequestBuilder {
	if builder.request.body == nil {
		builder.request.body = make(map[string]any)
	}

	builder.request.body[key] = value
	return builder
}

// Add query parameter.
func (builder *RequestBuilder) AddParam(key string, value string) *RequestBuilder {
	if builder.request.params == nil {
		builder.request.params = make(map[string]string)
	}

	builder.request.params[key] = value
	return builder
}

// Add header.
func (builder *RequestBuilder) AddHeader(key string, value string) *RequestBuilder {
	if builder.request.headers == nil {
		builder.request.headers = make(map[string]string)
	}

	builder.request.headers[key] = value
	return builder
}

// Add cookie to request header.
func (builder *RequestBuilder) AddCookie(cookie Cookie) *RequestBuilder {
	return builder.AddHeader("Cookie", cookie.ParseCookie())
}

// Add referer to request header.
func (builder *RequestBuilder) AddReferer(referer string) *RequestBuilder {
	return builder.AddHeader("Referer", referer)
}

// Add language to request header.
func (builder *RequestBuilder) AddLanguage(language string) *RequestBuilder {
	return builder.AddHeader("X-Rpc-Language", language)
}

// Add default headers that are required by HoYoLab endpoints.
func (builder *RequestBuilder) addDefaultHeaders() *RequestBuilder {
	builder.AddHeader("Accept", "application/json, text/plain, */*")
	builder.AddHeader("Content-Type", "application/json")
	builder.AddHeader("Accept-Encoding", fmt.Sprintf("%s, %s, %s", encodingGzip, encodingDeflate, encodingBr))
	builder.AddHeader("Sec-Ch-Ua", `"Chromium";v="112", "Microsoft Edge";v="112", "Not:A-Brand";v="99"`)
	builder.AddHeader("Sec-Ch-Ua-Mobile", "?0")
	builder.AddHeader("Sec-Ch-Ua-Platform", `"Windows"`)
	builder.AddHeader("Sec-Fetch-Dest", "empty")
	builder.AddHeader("Sec-Fetch-Mode", "cors")
	builder.AddHeader("Sec-Fetch-Site", "same-site")
	builder.AddHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.46")
	builder.AddHeader("x-rpc-app_version", "1.5.0")
	builder.AddHeader("x-rpc-client_type", "5")
	builder.AddHeader("x-rpc-language", "en-us")
	return builder
}
