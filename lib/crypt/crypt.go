package crypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Encoding/Hashing Functions
// These functions provide encoding, decoding, and hashing capabilities

// Base64 Functions
func base64Encode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("base64_encode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("base64_encode: %w", err)
	}
	return lang.StringValue(base64.StdEncoding.EncodeToString([]byte(string(str)))), nil
}

func base64Decode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("base64_decode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("base64_decode: %w", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(str))
	if err != nil {
		return nil, fmt.Errorf("base64_decode: invalid base64 string: %w", err)
	}
	return lang.StringValue(string(decoded)), nil
}

func base64UrlEncode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("base64_url_encode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("base64_url_encode: %w", err)
	}
	return lang.StringValue(base64.URLEncoding.EncodeToString([]byte(string(str)))), nil
}

func base64UrlDecode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("base64_url_decode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("base64_url_decode: %w", err)
	}
	decoded, err := base64.URLEncoding.DecodeString(string(str))
	if err != nil {
		return nil, fmt.Errorf("base64_url_decode: invalid base64url string: %w", err)
	}
	return lang.StringValue(string(decoded)), nil
}

// Hex Functions
func hexEncode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hex_encode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hex_encode: %w", err)
	}
	return lang.StringValue(hex.EncodeToString([]byte(string(str)))), nil
}

func hexDecode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hex_decode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hex_decode: %w", err)
	}
	decoded, err := hex.DecodeString(string(str))
	if err != nil {
		return nil, fmt.Errorf("hex_decode: invalid hex string: %w", err)
	}
	return lang.StringValue(string(decoded)), nil
}

// Hash Functions
func hashMD5(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hash_md5: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_md5: %w", err)
	}
	hash := md5.Sum([]byte(string(str)))
	return lang.StringValue(hex.EncodeToString(hash[:])), nil
}

func hashSHA1(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hash_sha1: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_sha1: %w", err)
	}
	hash := sha1.Sum([]byte(string(str)))
	return lang.StringValue(hex.EncodeToString(hash[:])), nil
}

func hashSHA224(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hash_sha224: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_sha224: %w", err)
	}
	hash := sha256.Sum224([]byte(string(str)))
	return lang.StringValue(hex.EncodeToString(hash[:])), nil
}

func hashSHA256(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hash_sha256: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_sha256: %w", err)
	}
	hash := sha256.Sum256([]byte(string(str)))
	return lang.StringValue(hex.EncodeToString(hash[:])), nil
}

func hashSHA384(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hash_sha384: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_sha384: %w", err)
	}
	hash := sha512.Sum384([]byte(string(str)))
	return lang.StringValue(hex.EncodeToString(hash[:])), nil
}

func hashSHA512(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hash_sha512: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_sha512: %w", err)
	}
	hash := sha512.Sum512([]byte(string(str)))
	return lang.StringValue(hex.EncodeToString(hash[:])), nil
}

func hashCRC32(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("hash_crc32: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_crc32: %w", err)
	}
	checksum := crc32.ChecksumIEEE([]byte(string(str)))
	return lang.NumberValue(float64(checksum)), nil
}

// HMAC Functions
func hmacMD5(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("hmac_md5: expected 2 arguments")
	}
	key, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hmac_md5: key %w", err)
	}
	message, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("hmac_md5: message %w", err)
	}

	h := hmac.New(md5.New, []byte(string(key)))
	h.Write([]byte(string(message)))
	return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
}

func hmacSHA1(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("hmac_sha1: expected 2 arguments")
	}
	key, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hmac_sha1: key %w", err)
	}
	message, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("hmac_sha1: message %w", err)
	}

	h := hmac.New(sha1.New, []byte(string(key)))
	h.Write([]byte(string(message)))
	return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
}

func hmacSHA256(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("hmac_sha256: expected 2 arguments")
	}
	key, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hmac_sha256: key %w", err)
	}
	message, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("hmac_sha256: message %w", err)
	}

	h := hmac.New(sha256.New, []byte(string(key)))
	h.Write([]byte(string(message)))
	return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
}

func hmacSHA512(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("hmac_sha512: expected 2 arguments")
	}
	key, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hmac_sha512: key %w", err)
	}
	message, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("hmac_sha512: message %w", err)
	}

	h := hmac.New(sha512.New, []byte(string(key)))
	h.Write([]byte(string(message)))
	return lang.StringValue(hex.EncodeToString(h.Sum(nil))), nil
}

// Binary/ASCII Functions
func toBinary(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("to_binary: expected 1 argument")
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
		return nil, fmt.Errorf("to_binary: unsupported type %T", args[0])
	}
}

func fromBinary(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("from_binary: expected 1 argument")
	}

	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("from_binary: %w", err)
	}
	binaryStr := string(str)

	// Handle space-separated binary (for strings)
	if strings.Contains(binaryStr, " ") {
		parts := strings.Split(binaryStr, " ")
		result := ""
		for _, part := range parts {
			if val, err := strconv.ParseInt(part, 2, 64); err == nil {
				result += string(byte(val))
			} else {
				return nil, fmt.Errorf("from_binary: invalid binary part '%s': %w", part, err)
			}
		}
		return lang.StringValue(result), nil
	}

	// Handle single binary number
	if val, err := strconv.ParseInt(binaryStr, 2, 64); err == nil {
		return lang.NumberValue(float64(val)), nil
	} else {
		return nil, fmt.Errorf("from_binary: invalid binary string '%s': %w", binaryStr, err)
	}
}

func toOctal(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("to_octal: expected 1 argument")
	}

	switch v := args[0].(type) {
	case lang.NumberValue:
		return lang.StringValue(strconv.FormatInt(int64(v), 8)), nil
	default:
		return nil, fmt.Errorf("to_octal: unsupported type %T", args[0])
	}
}

func fromOctal(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("from_octal: expected 1 argument")
	}

	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("from_octal: %w", err)
	}
	octalStr := string(str)
	if val, err := strconv.ParseInt(octalStr, 8, 64); err == nil {
		return lang.NumberValue(float64(val)), nil
	} else {
		return nil, fmt.Errorf("from_octal: invalid octal string '%s': %w", octalStr, err)
	}
}

// ASCII Functions
func toASCII(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("to_ascii: expected 1 argument")
	}

	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("to_ascii: %w", err)
	}
	result := make(lang.ListValue, len(string(str)))

	for i, char := range []byte(string(str)) {
		result[i] = lang.NumberValue(float64(char))
	}

	return result, nil
}

func fromASCII(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("from_ascii: expected 1 argument")
	}

	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("from_ascii: expected list, got %T", args[0])
	}

	result := make([]byte, len(list))
	for i, val := range list {
		num, err := lib.ToNumber(val)
		if err != nil {
			return nil, fmt.Errorf("from_ascii: element %d %w", i, err)
		}
		ascii := int(num)
		if ascii < 0 || ascii > 255 {
			return nil, fmt.Errorf("from_ascii: element %d value %d out of ASCII range (0-255)", i, ascii)
		}
		result[i] = byte(ascii)
	}

	return lang.StringValue(string(result)), nil
}

// URL-safe encoding functions
func base32Encode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("base32_encode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("base32_encode: %w", err)
	}
	input := []byte(string(str))
	return lang.StringValue(base32EncodeString(input)), nil
}

func base32Decode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("base32_decode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("base32_decode: %w", err)
	}
	decoded := base32DecodeString(string(str))
	if decoded == nil {
		return nil, errors.New("base32_decode: invalid base32 string")
	}
	return lang.StringValue(string(decoded)), nil
}

// Escape/Unescape Functions
func htmlEscape(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("html_escape: expected 1 argument")
	}

	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("html_escape: %w", err)
	}
	s := string(str)
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")

	return lang.StringValue(s), nil
}

func htmlUnescape(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("html_unescape: expected 1 argument")
	}

	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("html_unescape: %w", err)
	}
	s := string(str)
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&#x27;", "'")

	return lang.StringValue(s), nil
}

// Hash verification functions
func hashVerify(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("hash_verify: expected 3 arguments")
	}

	input, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("hash_verify: input %w", err)
	}
	expectedHash, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("hash_verify: expected hash %w", err)
	}
	algorithm, err := lib.ToString(args[2])
	if err != nil {
		return nil, fmt.Errorf("hash_verify: algorithm %w", err)
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
		return nil, fmt.Errorf("hash_verify: unsupported algorithm '%s'", algorithmStr)
	}

	return lang.BoolValue(actualHash == expectedHashStr), nil
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

	// Add padding
	for len(result)%8 != 0 {
		result += "="
	}

	return result
}

func base32DecodeString(s string) []byte {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

	// Remove padding
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
			return nil // Invalid character
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

// Functions that would be in the BuiltinFunctions map:
var EncodingFunctions = map[string]func([]lang.Value) (lang.Value, error){
	// Base64
	"base64_encode":     base64Encode,
	"base64_decode":     base64Decode,
	"base64_url_encode": base64UrlEncode,
	"base64_url_decode": base64UrlDecode,

	// Hex
	"hex_encode": hexEncode,
	"hex_decode": hexDecode,

	// Base32
	"base32_encode": base32Encode,
	"base32_decode": base32Decode,

	// Hashing
	"hash_md5":    hashMD5,
	"hash_sha1":   hashSHA1,
	"hash_sha224": hashSHA224,
	"hash_sha256": hashSHA256,
	"hash_sha384": hashSHA384,
	"hash_sha512": hashSHA512,
	"hash_crc32":  hashCRC32,

	// HMAC
	"hmac_md5":    hmacMD5,
	"hmac_sha1":   hmacSHA1,
	"hmac_sha256": hmacSHA256,
	"hmac_sha512": hmacSHA512,

	// Binary/ASCII
	"to_binary":   toBinary,
	"from_binary": fromBinary,
	"to_octal":    toOctal,
	"from_octal":  fromOctal,
	"to_ascii":    toASCII,
	"from_ascii":  fromASCII,

	// HTML
	"html_escape":   htmlEscape,
	"html_unescape": htmlUnescape,

	// Verification
	"hash_verify": hashVerify,
}
