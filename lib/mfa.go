package lib

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	codeLen = 6
)

// StepDurationType defines the type of TOTP duration.
type StepDurationType int

// These constants represent the two types of durations.
const (
	MFACode           StepDurationType = iota // For Multi-Factor Authentication.
	EmailVerification                         // For email code verification.
)

// stepDurations map duration types to their respective durations.
var stepDurations = [...]int64{
	30,  // MFA code lasts 30 seconds.
	600, // Email verification code lasts 10 minutes.
}

var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateString produces a random string of the specified length.
func GenerateString(length int) string {
	result := make([]byte, length)
	for i := range result {
		for {
			b := make([]byte, 1)
			if _, err := rand.Read(b); err != nil {
				panic(err)
			}
			if int(b[0]) < 256-(256%len(chars)) {
				result[i] = chars[int(b[0])%len(chars)]
				break
			}
		}
	}
	return string(result)
}

// ComputeTOTP computes the TOTP value for a given secret and time.
func ComputeTOTP(secret string, timestamp int64) string {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		fmt.Println("Error decoding secret:", err)
		return ""
	}

	t := make([]byte, 8)
	binary.BigEndian.PutUint64(t, uint64(timestamp))

	hm := hmac.New(sha1.New, key)
	hm.Write(t)
	hash := hm.Sum(nil)

	offset := hash[19] & 0xf
	code := hash[offset : offset+4]

	fullCode := binary.BigEndian.Uint32(code) & 0x7fffffff
	strCode := fmt.Sprintf("%06d", fullCode%1000000)

	return strCode
}

// GenerateTOTP provides a TOTP code for the current time and desired type (MFA or Email Verification).
func GenerateTOTP(secret string, stepType StepDurationType) (string, error) {
	if stepType < 0 || int(stepType) >= len(stepDurations) {
		return "", errors.New("invalid stepType provided")
	}
	return ComputeTOTP(secret, time.Now().Unix()/stepDurations[stepType]), nil
}

// ValidateTOTP verifies if the provided code matches the expected TOTP value for the given secret and duration type.
func ValidateTOTP(secret, code string, stepType StepDurationType) bool {
	expectedCode, err := GenerateTOTP(secret, stepType)
	if err != nil {
		return false
	}
	return subtleCompare(code, expectedCode)
}

// subtleCompare does a constant-time comparison of two strings.
func subtleCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	diff := 0
	for i := 0; i < len(a); i++ {
		diff |= int(a[i] ^ b[i])
	}
	return diff == 0
}
