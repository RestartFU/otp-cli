package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"math"
	"strings"
	"time"
)

func counter(tm time.Time, period float64) uint64 {
	return uint64(math.Floor(float64(tm.Unix()) / period))
}

// Generate generates a 6 digit OTP code using the given secret and counter
func Generate(secret string) (passcode int, err error) {
	// As noted in issue #10 and #17 this adds support for TOTP secrets that are
	// missing their padding.
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	// As noted in issue #24 Google has started producing base32 in lower case,
	// but the StdEncoding (and the RFC), expect a dictionary of only upper case letters.
	secret = strings.ToUpper(secret)

	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return 0, errors.New("decoding of secret as base32 failed")
	}

	buf := make([]byte, 8)
	mac := hmac.New(sha1.New, secretBytes)
	binary.BigEndian.PutUint64(buf, counter(time.Now(), 30))

	mac.Write(buf)
	sum := mac.Sum(nil)

	// "Dynamic truncation" in RFC 4226
	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	mod := int32(value % int64(math.Pow10(6)))
	return int(mod), nil
}
