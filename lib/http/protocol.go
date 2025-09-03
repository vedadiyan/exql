package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type (
	HttpProtocol interface {
		ContentLength() int64
		Cookies() []*http.Cookie
		Form() url.Values
		GetBody() (io.ReadCloser, error)
		Headers() http.Header
		Host() string
		Method() string
		Pattern() string
		Proto() string
		ProtoMajor() int
		ProtoMinor() int
		RemoteAddress() string
		StatusCode() int
		Trailers() http.Header
		TransferEncoding() []string
		Url() *url.URL

		Type() string
	}
	httpProtocol[T http.Request | http.Response] struct {
		v any
	}
)

func New[T http.Request | http.Response](v *T) *httpProtocol[T] {
	out := new(httpProtocol[T])
	out.v = v
	return out
}

func (hp httpProtocol[T]) ContentLength() int64 {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.ContentLength
		}
	case *http.Response:
		{
			return v.ContentLength
		}
	default:
		{
			return 0
		}
	}
}

func (hp httpProtocol[T]) Cookies() []*http.Cookie {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Cookies()
		}
	case *http.Response:
		{
			return v.Cookies()
		}
	default:
		{
			return nil
		}
	}
}

func (hp httpProtocol[T]) Form() url.Values {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Form
		}
	case *http.Response:
		{
			if v.Request == nil {
				return nil
			}
			return v.Request.Form
		}
	default:
		{
			return nil
		}
	}
}

func (hp httpProtocol[T]) GetBody() (io.ReadCloser, error) {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.GetBody()
		}
	case *http.Response:
		{
			copy, err := io.ReadAll(v.Body)
			if err != nil {
				return nil, err
			}
			v.Body = io.NopCloser(bytes.NewBuffer(copy))
			return io.NopCloser(bytes.NewBuffer(copy)), nil
		}
	default:
		{
			return nil, nil
		}
	}
}

func (hp httpProtocol[T]) Headers() http.Header {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Header
		}
	case *http.Response:
		{
			return v.Header
		}
	default:
		{
			return nil
		}
	}
}

func (hp httpProtocol[T]) Host() string {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Host
		}
	case *http.Response:
		{
			if v.Request == nil {
				return ""
			}
			return v.Request.Host
		}
	default:
		{
			return ""
		}
	}
}

func (hp httpProtocol[T]) Method() string {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Method
		}
	case *http.Response:
		{
			if v.Request == nil {
				return ""
			}
			return v.Request.Method
		}
	default:
		{
			return ""
		}
	}
}

func (hp httpProtocol[T]) Pattern() string {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Pattern
		}
	case *http.Response:
		{
			if v.Request == nil {
				return ""
			}
			return v.Request.Pattern
		}
	default:
		{
			return ""
		}
	}
}

func (hp httpProtocol[T]) Proto() string {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Proto
		}
	case *http.Response:
		{
			return v.Proto
		}
	default:
		{
			return ""
		}
	}
}

func (hp httpProtocol[T]) ProtoMajor() int {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.ProtoMajor
		}
	case *http.Response:
		{
			return v.ProtoMajor
		}
	default:
		{
			return 0
		}
	}
}

func (hp httpProtocol[T]) ProtoMinor() int {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.ProtoMinor
		}
	case *http.Response:
		{
			return v.ProtoMinor
		}
	default:
		{
			return 0
		}
	}
}

func (hp httpProtocol[T]) RemoteAddress() string {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.RemoteAddr
		}
	case *http.Response:
		{
			if v.Request == nil {
				return ""
			}
			return v.Request.RemoteAddr
		}
	default:
		{
			return ""
		}
	}
}

func (hp httpProtocol[T]) StatusCode() int {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			if v.Response == nil {
				return 0
			}
			return v.Response.StatusCode
		}
	case *http.Response:
		{
			return v.StatusCode
		}
	default:
		{
			return 0
		}
	}
}

func (hp httpProtocol[T]) Trailers() http.Header {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.Trailer
		}
	case *http.Response:
		{
			return v.Trailer
		}
	default:
		{
			return nil
		}
	}
}

func (hp httpProtocol[T]) TransferEncoding() []string {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.TransferEncoding
		}
	case *http.Response:
		{
			return v.TransferEncoding
		}
	default:
		{
			return nil
		}
	}
}

func (hp httpProtocol[T]) Url() *url.URL {
	switch v := hp.v.(type) {
	case *http.Request:
		{
			return v.URL
		}
	case *http.Response:
		{
			if v.Request == nil {
				return nil
			}
			return v.Request.URL
		}
	default:
		{
			return nil
		}
	}
}

func (hp httpProtocol[T]) Type() string {
	switch hp.v.(type) {
	case *http.Request:
		{
			return "request"
		}
	case *http.Response:
		{
			return "response"
		}
	default:
		{
			return ""
		}
	}
}
