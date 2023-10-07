package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func CheckCaptcha(captcha string) error {
	// don't need to verify actions under a dev env
	if Config.Debug {
		return nil
	}

	if captcha == "" {
		return ErrInvalidCaptcha
	}

	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", Config.CaptchaSecret)
	q.Add("response", captcha)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// idek
	defer func(Body io.ReadCloser) {
		if e := Body.Close(); e != nil {
			panic(e)
		}
	}(resp.Body)

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	// Check recaptcha verification success.
	if !body.Success || body.Score < 0.7 {
		return ErrInvalidCaptcha
	}

	return nil
}
