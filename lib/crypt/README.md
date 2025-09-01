# Crypt Package

The crypt package provides cryptographic and encoding functions for data transformation, hashing, and validation.

## Base64 Functions

### `base64Encode(data)`
Encodes a string to Base64 format.
- **Parameters:** `data` (string) - The string to encode
- **Returns:** Base64 encoded string
- **Example:** `base64Encode("hello")` → `"aGVsbG8="`

### `base64Decode(data)`
Decodes a Base64 string to its original form.
- **Parameters:** `data` (string) - The Base64 string to decode
- **Returns:** Decoded string
- **Example:** `base64Decode("aGVsbG8=")` → `"hello"`

### `base64UrlEncode(data)`
Encodes a string to URL-safe Base64 format.
- **Parameters:** `data` (string) - The string to encode
- **Returns:** URL-safe Base64 encoded string
- **Example:** `base64UrlEncode("hello?world")` → `"aGVsbG8_d29ybGQ="`

### `base64UrlDecode(data)`
Decodes a URL-safe Base64 string to its original form.
- **Parameters:** `data` (string) - The URL-safe Base64 string to decode
- **Returns:** Decoded string
- **Example:** `base64UrlDecode("aGVsbG8_d29ybGQ=")` → `"hello?world"`

## Base32 Functions

### `base32Encode(data)`
Encodes a string to Base32 format.
- **Parameters:** `data` (string) - The string to encode
- **Returns:** Base32 encoded string
- **Example:** `base32Encode("hello")` → `"NBSWY3DP"`

### `base32Decode(data)`
Decodes a Base32 string to its original form.
- **Parameters:** `data` (string) - The Base32 string to decode
- **Returns:** Decoded string
- **Example:** `base32Decode("NBSWY3DP")` → `"hello"`

## Hex Functions

### `hexEncode(data)`
Encodes a string to hexadecimal format.
- **Parameters:** `data` (string) - The string to encode
- **Returns:** Hexadecimal encoded string
- **Example:** `hexEncode("hello")` → `"68656c6c6f"`

### `hexDecode(data)`
Decodes a hexadecimal string to its original form.
- **Parameters:** `data` (string) - The hexadecimal string to decode
- **Returns:** Decoded string
- **Example:** `hexDecode("68656c6c6f")` → `"hello"`

## Hash Functions

### `hashMd5(data)`
Generates an MD5 hash of the input string.
- **Parameters:** `data` (string) - The string to hash
- **Returns:** MD5 hash as hexadecimal string
- **Example:** `hashMd5("hello")` → `"5d41402abc4b2a76b9719d911017c592"`

### `hashSha1(data)`
Generates a SHA-1 hash of the input string.
- **Parameters:** `data` (string) - The string to hash
- **Returns:** SHA-1 hash as hexadecimal string
- **Example:** `hashSha1("hello")` → `"aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"`

### `hashSha224(data)`
Generates a SHA-224 hash of the input string.
- **Parameters:** `data` (string) - The string to hash
- **Returns:** SHA-224 hash as hexadecimal string

### `hashSha256(data)`
Generates a SHA-256 hash of the input string.
- **Parameters:** `data` (string) - The string to hash
- **Returns:** SHA-256 hash as hexadecimal string
- **Example:** `hashSha256("hello")` → `"2cf24dba4f21d4288094e4966731b9347d8e88c0"`

### `hashSha384(data)`
Generates a SHA-384 hash of the input string.
- **Parameters:** `data` (string) - The string to hash
- **Returns:** SHA-384 hash as hexadecimal string

### `hashSha512(data)`
Generates a SHA-512 hash of the input string.
- **Parameters:** `data` (string) - The string to hash
- **Returns:** SHA-512 hash as hexadecimal string

### `hashCrc32(data)`
Generates a CRC32 checksum of the input string.
- **Parameters:** `data` (string) - The string to checksum
- **Returns:** CRC32 checksum as number
- **Example:** `hashCrc32("hello")` → `907060870`

## HMAC Functions

### `hmacMd5(key, message)`
Generates an HMAC-MD5 signature.
- **Parameters:** 
  - `key` (string) - The secret key
  - `message` (string) - The message to sign
- **Returns:** HMAC-MD5 signature as hexadecimal string

### `hmacSha1(key, message)`
Generates an HMAC-SHA1 signature.
- **Parameters:** 
  - `key` (string) - The secret key
  - `message` (string) - The message to sign
- **Returns:** HMAC-SHA1 signature as hexadecimal string

### `hmacSha256(key, message)`
Generates an HMAC-SHA256 signature.
- **Parameters:** 
  - `key` (string) - The secret key
  - `message` (string) - The message to sign
- **Returns:** HMAC-SHA256 signature as hexadecimal string

### `hmacSha512(key, message)`
Generates an HMAC-SHA512 signature.
- **Parameters:** 
  - `key` (string) - The secret key
  - `message` (string) - The message to sign
- **Returns:** HMAC-SHA512 signature as hexadecimal string

## Binary/ASCII Functions

### `toBinary(input)`
Converts a number or string to binary representation.
- **Parameters:** `input` (number|string) - The value to convert
- **Returns:** Binary representation as string
- **Example:** 
  - `toBinary(10)` → `"1010"`
  - `toBinary("hi")` → `"01101000 01101001"`

### `fromBinary(binary)`
Converts binary string back to number or string.
- **Parameters:** `binary` (string) - Binary string to convert
- **Returns:** Converted value (number for single binary, string for space-separated)
- **Example:** 
  - `fromBinary("1010")` → `10`
  - `fromBinary("01101000 01101001")` → `"hi"`

### `toOctal(number)`
Converts a number to octal representation.
- **Parameters:** `number` (number) - The number to convert
- **Returns:** Octal representation as string
- **Example:** `toOctal(8)` → `"10"`

### `fromOctal(octal)`
Converts octal string back to number.
- **Parameters:** `octal` (string) - Octal string to convert
- **Returns:** Number value
- **Example:** `fromOctal("10")` → `8`

### `toAscii(text)`
Converts a string to ASCII code array.
- **Parameters:** `text` (string) - The string to convert
- **Returns:** Array of ASCII codes
- **Example:** `toAscii("hi")` → `[104, 105]`

### `fromAscii(codes)`
Converts ASCII code array back to string.
- **Parameters:** `codes` (array) - Array of ASCII codes
- **Returns:** Converted string
- **Example:** `fromAscii([104, 105])` → `"hi"`

## HTML Functions

### `htmlEscape(text)`
Escapes HTML special characters in a string.
- **Parameters:** `text` (string) - The text to escape
- **Returns:** HTML-escaped string
- **Example:** `htmlEscape("<script>")` → `"&lt;script&gt;"`

### `htmlUnescape(text)`
Unescapes HTML entities in a string.
- **Parameters:** `text` (string) - The HTML-escaped text
- **Returns:** Unescaped string
- **Example:** `htmlUnescape("&lt;script&gt;")` → `"<script>"`

## Hash Verification

### `hashVerify(input, expectedHash, algorithm)`
Verifies if an input matches the expected hash using specified algorithm.
- **Parameters:** 
  - `input` (string) - The input to verify
  - `expectedHash` (string) - The expected hash value
  - `algorithm` (string) - Hash algorithm ("md5", "sha1", "sha256", "sha512")
- **Returns:** Boolean indicating if the hash matches
- **Example:** `hashVerify("hello", "5d41402abc4b2a76b9719d911017c592", "md5")` → `true`

## Usage Notes

- All hash functions return lowercase hexadecimal strings
- Binary functions handle both individual numbers and space-separated binary strings for text
- Base64 and Base32 functions handle padding automatically
- HMAC functions require both a key and message parameter
- HTML escape/unescape functions handle standard HTML entities