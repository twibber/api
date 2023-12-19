package img

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	cfg "github.com/twibber/api/config"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// IMGConfig is a configuration for the imgproxy URL and image settings
type IMGConfig struct {
	Width   int
	Height  int
	Quality int // 0-100
}

// SignImageURL signs an image URL with the given key and salt for imgproxy
func SignImageURL(imgURL string, config ...IMGConfig) string {
	escapedImgURL := url.QueryEscape(imgURL + "?t=" + strconv.Itoa(time.Now().Nanosecond()))

	// Default configuration
	var imgConfig IMGConfig
	if len(config) > 0 {
		imgConfig = config[0]
	}

	// Building the path with only provided parameters
	var path strings.Builder

	// Adding the image settings
	if imgConfig.Width != 0 && imgConfig.Height != 0 {
		fmt.Fprintf(&path, "/rs:auto:%d:%d:0", imgConfig.Width, imgConfig.Height)
	}
	if imgConfig.Quality != 0 {
		fmt.Fprintf(&path, "/q:%d", imgConfig.Quality)
	}
	// Always setting the format to WebP
	fmt.Fprintf(&path, "/f:webp")

	// Adding the image URL
	fmt.Fprintf(&path, "/plain/%s", escapedImgURL)

	keyBin, _ := hex.DecodeString(cfg.Config.ImgproxyKey)
	saltBin, _ := hex.DecodeString(cfg.Config.ImgproxySalt)

	// Calculate the HMAC digest
	mac := hmac.New(sha256.New, keyBin)
	mac.Write(saltBin)               // Writing salt first
	mac.Write([]byte(path.String())) // Writing the path
	signature := mac.Sum(nil)

	// Base64 URL-Safe Encoding of the signature
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)
	encodedSignature = strings.TrimRight(encodedSignature, "=")

	// Construct the final signed URL
	return fmt.Sprintf("%s/%s%s", cfg.Config.ImgproxyURL, encodedSignature, path.String())
}
