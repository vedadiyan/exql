/*
 * Copyright 2025 Pouya Vedadiyan
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package http

import (
	"fmt"
	"io"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func header() (string, lang.Function) {
	name := "header"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		headerName, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: header name %w", name, err)
		}

		return lang.StringValue(protocol.Headers().Get(string(headerName))), nil
	}
	return name, fn
}

func headers() (string, lang.Function) {
	name := "headers"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		value := make(lang.MapValue)

		for key, val := range protocol.Headers() {
			values := make(lang.ListValue, len(val))
			for i := 0; i < len(val); i++ {
				values[i] = val[i]
			}
			value[key] = values
		}
		return value, nil
	}
	return name, fn
}

func method() (string, lang.Function) {
	name := "method"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		return lang.StringValue(protocol.Method()), nil
	}
	return name, fn
}

func path() (string, lang.Function) {
	name := "path"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Url().Path), nil
	}
	return name, fn
}

func query() (string, lang.Function) {
	name := "query"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		value := make(lang.MapValue)

		for key, val := range protocol.Url().Query() {
			values := make(lang.ListValue, len(val))
			for i := 0; i < len(val); i++ {
				values[i] = val[i]
			}
			value[key] = values
		}
		return value, nil
	}
	return name, fn
}

func queryParam() (string, lang.Function) {
	name := "queryParam"
	_, Query := query()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		queryParams, err := Query(args[:1])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		paramName, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: param name %w", name, err)
		}

		if queryMap, ok := queryParams.(lang.MapValue); ok {
			return queryMap[string(paramName)], nil
		}
		return nil, nil
	}
	return name, fn
}

func body() (string, lang.Function) {
	name := "body"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		body, err := protocol.GetBody()
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(body)
		if err != nil {
			return nil, err
		}
		return lang.StringValue(string(data)), nil
	}
	return name, fn
}

func status() (string, lang.Function) {
	name := "status"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		return lang.NumberValue(protocol.StatusCode()), nil
	}
	return name, fn
}

func ip() (string, lang.Function) {
	name := "ip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		xForwardedFor := protocol.Headers().Get("X-Forwarded-For")
		if xForwardedFor != "" {
			return lang.StringValue(strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])), nil
		}
		xRealIP := protocol.Headers().Get("X-Real-IP")
		if xRealIP != "" {
			return lang.StringValue(xRealIP), nil
		}
		return lang.StringValue(protocol.RemoteAddress()), nil
	}
	return name, fn
}

func userAgent() (string, lang.Function) {
	name := "userAgent"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Headers().Get("User-Agent")), nil
	}
	return name, fn
}

func contentType() (string, lang.Function) {
	name := "contentType"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Headers().Get("Content-Type")), nil
	}
	return name, fn
}

func contentLength() (string, lang.Function) {
	name := "contentLength"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.NumberValue(protocol.ContentLength()), nil
	}
	return name, fn
}

func host() (string, lang.Function) {
	name := "host"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Host()), nil
	}
	return name, fn
}

func scheme() (string, lang.Function) {
	name := "scheme"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Url().Scheme), nil
	}
	return name, fn
}

func port() (string, lang.Function) {
	name := "port"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Url().Port()), nil
	}
	return name, fn
}

func cookies() (string, lang.Function) {
	name := "cookies"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		cookies := protocol.Cookies()

		values := make(lang.ListValue, len(cookies))

		for i := 0; i < len(cookies); i++ {
			value := make(lang.MapValue)
			ref := cookies[i]
			value["domain"] = lang.StringValue(ref.Domain)
			value["expires"] = lang.StringValue(ref.Expires.String())
			value["httpOnly"] = lang.BoolValue(ref.HttpOnly)
			value["maxAge"] = lang.NumberValue(ref.MaxAge)
			value["name"] = lang.StringValue(ref.Name)
			value["partitioned"] = lang.BoolValue(ref.Partitioned)
			value["path"] = lang.StringValue(ref.Path)
			value["quoted"] = lang.BoolValue(ref.Quoted)
			value["raw"] = lang.StringValue(ref.Raw)
			value["rawExpires"] = lang.StringValue(ref.RawExpires)
			value["sameSite"] = lang.NumberValue(ref.SameSite)
			value["secure"] = lang.BoolValue(ref.Secure)
			unparsed := make(lang.ListValue, len(ref.Unparsed))
			for x := 0; x < len(ref.Unparsed); x++ {
				unparsed[x] = lang.StringValue(ref.Unparsed[x])
			}
			value["unparsed"] = unparsed
			value["value"] = lang.StringValue(ref.Value)
			values[i] = value
		}

		return values, nil
	}
	return name, fn
}

func cookie() (string, lang.Function) {
	name := "cookie"
	_, Cookies := cookies()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		cookies, err := Cookies(args[:1])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		cookieName, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: cookie name %w", name, err)
		}

		if cookiesMap, ok := cookies.(lang.MapValue); ok {
			return cookiesMap[string(cookieName)], nil
		}
		return nil, nil
	}
	return name, fn
}

func referer() (string, lang.Function) {
	name := "referer"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Headers().Get("Referer")), nil
	}
	return name, fn
}

func authorization() (string, lang.Function) {
	name := "authorization"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Headers().Get("Authorization")), nil
	}
	return name, fn
}

func accept() (string, lang.Function) {
	name := "accept"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}
		return lang.StringValue(protocol.Headers().Get("Accept")), nil
	}
	return name, fn
}

var httpFunctions = []func() (string, lang.Function){
	header,
	headers,
	method,
	path,
	query,
	queryParam,
	body,
	status,
	ip,
	userAgent,
	contentType,
	contentLength,
	host,
	scheme,
	port,
	cookies,
	cookie,
	referer,
	authorization,
	accept,
}

func Export() map[string]lang.Function {
	out := make(map[string]lang.Function)
	for _, value := range httpFunctions {
		name, fn := value()
		out[name] = fn
	}
	return out
}
