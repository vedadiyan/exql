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

func headerFn() (string, lang.Function) {
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

func headersFn() (string, lang.Function) {
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

func trailerFn() (string, lang.Function) {
	name := "trailer"
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

		return lang.StringValue(protocol.Trailers().Get(string(headerName))), nil
	}
	return name, fn
}

func trailersFn() (string, lang.Function) {
	name := "trailers"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		value := make(lang.MapValue)

		for key, val := range protocol.Trailers() {
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

func routeValuesFn() (string, lang.Function) {
	name := "routeValues"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		value := make(lang.MapValue)

		pattern := strings.Split(protocol.Pattern(), "/")
		path := strings.Split(protocol.Url().Path, "/")
		for i := 0; i < len(pattern); i++ {
			value[strings.TrimLeft(pattern[i], ":")] = path[i]
		}
		return value, nil
	}
	return name, fn
}

func patternFn() (string, lang.Function) {
	name := "pattern"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		return lang.StringValue(protocol.Pattern()), nil
	}
	return name, fn
}

func protoFn() (string, lang.Function) {
	name := "proto"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		return lang.StringValue(protocol.Proto()), nil
	}
	return name, fn
}

func protoMajorFn() (string, lang.Function) {
	name := "protoMajor"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		return lang.NumberValue(protocol.ProtoMajor()), nil
	}
	return name, fn
}

func protoMinorFn() (string, lang.Function) {
	name := "protoMinor"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		return lang.NumberValue(protocol.ProtoMinor()), nil
	}
	return name, fn
}

func transferEncodingFn() (string, lang.Function) {
	name := "transferEncoding"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		transferEncoding := protocol.TransferEncoding()

		value := make(lang.ListValue, len(transferEncoding))

		for i := 0; i < len(transferEncoding); i++ {
			value = append(value, lang.StringValue(transferEncoding[i]))
		}

		return value, nil
	}
	return name, fn
}

func methodFn() (string, lang.Function) {
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

func pathFn() (string, lang.Function) {
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

func queryFn() (string, lang.Function) {
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

func queryParamFn() (string, lang.Function) {
	name := "queryParam"
	_, Query := queryFn()
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

func bodyFn() (string, lang.Function) {
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

func statusFn() (string, lang.Function) {
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

func ipFn() (string, lang.Function) {
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

func userAgentFn() (string, lang.Function) {
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

func contentTypeFn() (string, lang.Function) {
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

func contentLengthFn() (string, lang.Function) {
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

func hostFn() (string, lang.Function) {
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

func schemeFn() (string, lang.Function) {
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

func portFn() (string, lang.Function) {
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

func cookiesFn() (string, lang.Function) {
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

func cookieFn() (string, lang.Function) {
	name := "cookie"
	_, Cookies := cookiesFn()
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

func urlFn() (string, lang.Function) {
	name := "url"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		protocol, ok := args[0].(HttpProtocol)
		if !ok {
			return nil, lib.ArgumenErrorType(name, 0, "HttpProtocol", args[0])
		}

		url := protocol.Url()

		value := make(lang.MapValue)
		value["forceQuery"] = url.ForceQuery
		value["fragment"] = url.Fragment
		value["host"] = url.Host
		value["omitHost"] = url.OmitHost
		value["opaque"] = url.Opaque
		value["path"] = url.Path
		value["rawFragment"] = url.RawFragment
		value["rawPath"] = url.RawPath
		value["rawQuery"] = url.RawQuery
		value["scheme"] = url.Scheme
		user := make(lang.MapValue)
		user["username"] = url.User.Username()
		user["password"], _ = url.User.Password()
		value["user"] = user

		return value, nil
	}
	return name, fn
}

func refererFn() (string, lang.Function) {
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

func authorizationFn() (string, lang.Function) {
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

func acceptFn() (string, lang.Function) {
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
	headerFn,
	headersFn,
	methodFn,
	pathFn,
	queryFn,
	queryParamFn,
	bodyFn,
	statusFn,
	ipFn,
	userAgentFn,
	contentTypeFn,
	contentLengthFn,
	hostFn,
	schemeFn,
	portFn,
	cookiesFn,
	cookieFn,
	refererFn,
	authorizationFn,
	acceptFn,
	trailerFn,
	trailersFn,
	routeValuesFn,
	patternFn,
	protoFn,
	protoMajorFn,
	protoMinorFn,
	transferEncodingFn,
	urlFn,
}

func Export() map[string]lang.Function {
	out := make(map[string]lang.Function)
	for _, value := range httpFunctions {
		name, fn := value()
		out[name] = fn
	}
	return out
}
