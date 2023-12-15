package lib

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

// SiteVerifyResponse struct maps the JSON response from reCAPTCHA verification.
type SiteVerifyResponse struct {
	Success     bool      `json:"success"`      // Indicates if the captcha was successful
	Score       float64   `json:"score"`        // Score for the captcha action
	Action      string    `json:"action"`       // Action associated with the captcha
	ChallengeTS time.Time `json:"challenge_ts"` // Timestamp of the captcha challenge
	Hostname    string    `json:"hostname"`     // Hostname of the site where the captcha was solved
	ErrorCodes  []string  `json:"error-codes"`  // Any error codes returned by the verification
}

// CheckCaptcha verifies the reCAPTCHA response token.
func CheckCaptcha(captcha string) error {
	// Skip verification in development environment.
	if Config.Debug {
		return nil
	}

	// Return error if captcha token is empty.
	if captcha == "" {
		log.Warn("reCAPTCHA verification failed: captcha token is empty")
		return ErrInvalidCaptcha
	}

	// Prepare a POST request to the Google reCAPTCHA API.
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return err
	}

	// Set the necessary query parameters for the request.
	q := req.URL.Query()
	q.Add("secret", Config.CaptchaSecret)
	q.Add("response", captcha)
	req.URL.RawQuery = q.Encode()

	// Execute the HTTP request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Ensure response body is closed after function execution.
	defer func(Body io.ReadCloser) {
		if e := Body.Close(); e != nil {
			panic(e) // Panic if closing the response body fails, might replace with logging
		}
	}(resp.Body)

	// Parse the JSON response from the reCAPTCHA server.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	// Validate if reCAPTCHA verification was successful and score is above the threshold.
	if !body.Success || body.Score < 0.3 {
		log.WithFields(log.Fields{
			"success":      body.Success,
			"score":        body.Score,
			"action":       body.Action,
			"challenge_ts": body.ChallengeTS,
			"hostname":     body.Hostname,
			"error_codes":  body.ErrorCodes,
		}).Warn("reCAPTCHA verification failed")

		return ErrInvalidCaptcha
	}

	return nil
}
