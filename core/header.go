package core

import (
	"fmt"
	"net/http"
)

func (c *Client) setHeader(http *http.Request) {

	userAgent, os := generateRandomUserAgent()
	if userAgent == "" || os == "" {
		userAgent = "Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-A500FU Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.2 Chrome/38.0.2125.102 Mobile Safari/537.36"
		os = "Android"
	}

	header := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9,id;q=0.8",
		"content-type":       "application/json",
		"priority":           "u=1, i",
		"sec-ch-ua":          `"Android WebView";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`,
		"sec-ch-ua-platform": fmt.Sprintf("\"%s\"", os),
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"Origin":             "https://major.bot",
		"Referer":            "https://major.bot/",
		"Referrer-Policy":    "strict-origin-when-cross-origin",
		"X-Requested-With":   "org.telegram.messenger.web",
		"User-Agent":         userAgent,
	}

	if c.accessToken != "" {
		header["authorization"] = c.accessToken
	}

	for key, value := range header {
		http.Header.Set(key, value)
	}
}
