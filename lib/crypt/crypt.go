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
package crypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"html"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Base64 Functions
func base64Encode() (string, lang.Function) {
	name := "base64Encode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.StringValue(base64.StdEncoding.EncodeToString([]byte(string(str)))), nil
	}
	return name, fn
}

func base64Decode() (string, lang.Function) {
	name := "base64Decode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		decoded, err := base64.StdEncoding.DecodeString(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid base64 string: %w", name, err)
		}
		return lang.StringValue(string(decoded)), nil
	}
	return name, fn
}

func base64UrlEncode() (string, lang.Function) {
	name := "base64UrlEncode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.StringValue(base64.URLEncoding.EncodeToString([]byte(string(str)))), nil
	}
	return name, fn
}

func base64UrlDecode() (string, lang.Function) {
	name := "base64UrlDecode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		decoded, err := base64.URLEncoding.DecodeString(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid base64url string: %w", name, err)
		}
		return lang.StringValue(string(decoded)), nil
	}
	return name, fn
}

// Hex Functions
func hexEncode() (string, lang.Function) {
	name := "hexEncode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.StringValue(hex.EncodeToString([]byte(string(str)))), nil
	}
	return name, fn
}

func hexDecode() (string, lang.Function) {
	name := "hexDecode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		decoded, err := hex.DecodeString(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid hex string: %w", name, err)
		}
		return lang.StringValue(string(decoded)), nil
	}
	return name, fn
}

// Hash Functions
func hashMD5() (string, lang.Function) {
	name := "hashMd5"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hash := md5.Sum([]byte(string(str)))
		return lang.StringValue(hex.EncodeToString(hash[:])), nil
	}
	return name, fn
}

func hashSHA1() (string, lang.Function) {
	name := "hashSha1"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hash := sha1.Sum([]byte(string(str)))
		return lang.StringValue(hex.EncodeToString(hash[:])), nil
	}
	return name, fn
}

func hashSHA224() (string, lang.Function) {
	name := "hashSha224"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hash := sha256.Sum224([]byte(string(str)))
		return lang.StringValue(hex.EncodeToString(hash[:])), nil
	}
	return name, fn
}

func hashSHA256() (string, lang.Function) {
	name := "hashSha256"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hash := sha256.Sum256([]byte(string(str)))
		return lang.StringValue(hex.EncodeToString(hash[:])), nil
	}
	return name, fn
}

func hashSHA384() (string, lang.Function) {
	name := "hashSha384"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hash := sha512.Sum384([]byte(string(str)))
		return lang.StringValue(hex.EncodeToString(hash[:])), nil
	}
	return name, fn
}

func hashSHA512() (string, lang.Function) {
	name := "hashSha512"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hash := sha512.Sum512([]byte(string(str)))
		return lang.StringValue(hex.EncodeToString(hash[:])), nil
	}
	return name, fn
}

func hashCRC32() (string, lang.Function) {
	name := "hashCrc32"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		checksum := crc32.ChecksumIEEE([]byte(string(str)))
		return lang.NumberValue(float64(checksum)), nil
	}
	return name, fn
}

// HMAC Functions
func hmacMD5() (string, lang.Function) {
	name := "hmacMd5"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		key, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
		}
		message, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: message %w", name, err)
		}

		h := hmac.New(md5.New, []byte(string(key)))
		h.Write([]byte(string(message)))
		return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
	}
	return name, fn
}

func hmacSHA1() (string, lang.Function) {
	name := "hmacSha1"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		key, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
		}
		message, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: message %w", name, err)
		}

		h := hmac.New(sha1.New, []byte(string(key)))
		h.Write([]byte(string(message)))
		return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
	}
	return name, fn
}

func hmacSHA256() (string, lang.Function) {
	name := "hmacSha256"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		key, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
		}
		message, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: message %w", name, err)
		}

		h := hmac.New(sha256.New, []byte(string(key)))
		h.Write([]byte(string(message)))
		return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
	}
	return name, fn
}

func hmacSHA512() (string, lang.Function) {
	name := "hmacSha512"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		key, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
		}
		message, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: message %w", name, err)
		}

		h := hmac.New(sha512.New, []byte(string(key)))
		h.Write([]byte(string(message)))
		return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
	}
	return name, fn
}

// Binary/ASCII Functions
func toBinary() (string, lang.Function) {
	name := "toBinary"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		switch v := args[0].(type) {
		case lang.NumberValue:
			return lang.StringValue(strconv.FormatInt(int64(v), 2)), nil
		case lang.StringValue:
			result := ""
			for _, char := range []byte(string(v)) {
				if result != "" {
					result += " "
				}
				result += fmt.Sprintf("%08b", char)
			}
			return lang.StringValue(result), nil
		default:
			return nil, fmt.Errorf("%s: unsupported type %T", name, args[0])
		}
	}
	return name, fn
}

func fromBinary() (string, lang.Function) {
	name := "fromBinary"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		binaryStr := string(str)

		if strings.Contains(binaryStr, " ") {
			parts := strings.Split(binaryStr, " ")
			result := ""
			for _, part := range parts {
				if val, err := strconv.ParseInt(part, 2, 64); err == nil {
					result += string(byte(val))
				} else {
					return nil, fmt.Errorf("%s: invalid binary part '%s': %w", name, part, err)
				}
			}
			return lang.StringValue(result), nil
		}

		if val, err := strconv.ParseInt(binaryStr, 2, 64); err == nil {
			return lang.NumberValue(float64(val)), nil
		} else {
			return nil, fmt.Errorf("%s: invalid binary string '%s': %w", name, binaryStr, err)
		}
	}
	return name, fn
}

func toOctal() (string, lang.Function) {
	name := "toOctal"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		switch v := args[0].(type) {
		case lang.NumberValue:
			return lang.StringValue(strconv.FormatInt(int64(v), 8)), nil
		default:
			return nil, fmt.Errorf("%s: unsupported type %T", name, args[0])
		}
	}
	return name, fn
}

func fromOctal() (string, lang.Function) {
	name := "fromOctal"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		octalStr := string(str)
		if val, err := strconv.ParseInt(octalStr, 8, 64); err == nil {
			return lang.NumberValue(float64(val)), nil
		} else {
			return nil, fmt.Errorf("%s: invalid octal string '%s': %w", name, octalStr, err)
		}
	}
	return name, fn
}

// ASCII Functions
func toASCII() (string, lang.Function) {
	name := "toAscii"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		result := make(lang.ListValue, len(string(str)))

		for i, char := range []byte(string(str)) {
			result[i] = lang.NumberValue(float64(char))
		}

		return result, nil
	}
	return name, fn
}

func fromASCII() (string, lang.Function) {
	name := "fromAscii"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, fmt.Errorf("%s: expected list, got %T", name, args[0])
		}

		result := make([]byte, len(list))
		for i, val := range list {
			num, err := lib.ToNumber(val)
			if err != nil {
				return nil, fmt.Errorf("%s: element %d %w", name, i, err)
			}
			ascii := int(num)
			if ascii < 0 || ascii > 255 {
				return nil, fmt.Errorf("%s: element %d value %d out of ASCII range (0-255)", name, i, ascii)
			}
			result[i] = byte(ascii)
		}

		return lang.StringValue(string(result)), nil
	}
	return name, fn
}

// URL-safe encoding functions
func base32Encode() (string, lang.Function) {
	name := "base32Encode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		input := []byte(string(str))
		return lang.StringValue(base32EncodeString(input)), nil
	}
	return name, fn
}

func base32Decode() (string, lang.Function) {
	name := "base32Decode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		decoded := base32DecodeString(string(str))
		if decoded == nil {
			return nil, fmt.Errorf("%s: invalid base32 string", name)
		}
		return lang.StringValue(string(decoded)), nil
	}
	return name, fn
}

// Escape/Unescape Functions
func htmlEscape() (string, lang.Function) {
	name := "htmlEscape"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.StringValue(html.EscapeString(string(str))), nil
	}
	return name, fn
}

func htmlUnescape() (string, lang.Function) {
	name := "htmlUnescape"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.StringValue(html.UnescapeString(string(str))), nil
	}
	return name, fn
}

// Hash verification functions
func hashVerify() (string, lang.Function) {
	name := "hashVerify"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}

		input, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: input %w", name, err)
		}
		expectedHash, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: expected hash %w", name, err)
		}
		algorithm, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: algorithm %w", name, err)
		}

		inputStr := string(input)
		expectedHashStr := strings.ToLower(string(expectedHash))
		algorithmStr := strings.ToLower(string(algorithm))

		var actualHash string

		switch algorithmStr {
		case "md5":
			hash := md5.Sum([]byte(inputStr))
			actualHash = hex.EncodeToString(hash[:])
		case "sha1":
			hash := sha1.Sum([]byte(inputStr))
			actualHash = hex.EncodeToString(hash[:])
		case "sha256":
			hash := sha256.Sum256([]byte(inputStr))
			actualHash = hex.EncodeToString(hash[:])
		case "sha512":
			hash := sha512.Sum512([]byte(inputStr))
			actualHash = hex.EncodeToString(hash[:])
		default:
			return nil, fmt.Errorf("%s: unsupported algorithm '%s'", name, algorithmStr)
		}

		return lang.BoolValue(actualHash == expectedHashStr), nil
	}
	return name, fn
}

// Simple Base32 implementation
func base32EncodeString(data []byte) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	if len(data) == 0 {
		return ""
	}

	result := ""
	buffer := 0
	bitsLeft := 0

	for _, b := range data {
		buffer = (buffer << 8) | int(b)
		bitsLeft += 8

		for bitsLeft >= 5 {
			result += string(alphabet[(buffer>>(bitsLeft-5))&0x1F])
			bitsLeft -= 5
		}
	}

	if bitsLeft > 0 {
		result += string(alphabet[(buffer<<(5-bitsLeft))&0x1F])
	}

	for len(result)%8 != 0 {
		result += "="
	}

	return result
}

func base32DecodeString(s string) []byte {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

	s = strings.TrimRight(s, "=")

	if len(s) == 0 {
		return []byte{}
	}

	var result []byte
	buffer := 0
	bitsLeft := 0

	for _, c := range s {
		idx := strings.IndexRune(alphabet, c)
		if idx == -1 {
			return nil
		}

		buffer = (buffer << 5) | idx
		bitsLeft += 5

		if bitsLeft >= 8 {
			result = append(result, byte(buffer>>(bitsLeft-8)))
			bitsLeft -= 8
		}
	}

	return result
}

var encodingFunctions = []func() (string, lang.Function){
	base64Encode,
	base64Decode,
	base64UrlEncode,
	base64UrlDecode,
	hexEncode,
	hexDecode,
	base32Encode,
	base32Decode,
	hashMD5,
	hashSHA1,
	hashSHA224,
	hashSHA256,
	hashSHA384,
	hashSHA512,
	hashCRC32,
	hmacMD5,
	hmacSHA1,
	hmacSHA256,
	hmacSHA512,
	toBinary,
	fromBinary,
	toOctal,
	fromOctal,
	toASCII,
	fromASCII,
	htmlEscape,
	htmlUnescape,
	hashVerify,
}

func Export() map[string]lang.Function {
	out := make(map[string]lang.Function)
	for _, value := range encodingFunctions {
		name, fn := value()
		out[name] = fn
	}
	return out
}
