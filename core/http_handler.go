package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func setDns(dialer Dialer) *net.Resolver {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, "1.1.1.1:53")
		},
	}

	return resolver
}

func (c *Client) setProxy() error {
	// Parse the proxy URL
	parsedURL, err := url.Parse(c.proxy)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %v", err)
	}

	// Extract username and password if available
	var proxyUser, proxyPass string
	if parsedURL.User != nil {
		proxyUser = parsedURL.User.Username()
		proxyPass, _ = parsedURL.User.Password()
	}

	// Handle based on scheme
	switch parsedURL.Scheme {
	case "http", "https":
		// HTTP/HTTPS Proxy
		transport := &http.Transport{
			Proxy: http.ProxyURL(parsedURL), // Handles HTTP Proxy with auth
		}

		c.httpClient = &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
		}

		// Test the proxy connection by making a request
		req, err := http.NewRequest("GET", "https://ipinfo.io/json", nil) // Test with a simple request
		if err != nil {
			// Fallback to no proxy
			c.httpClient = &http.Client{
				Transport: nil,
				Timeout:   15 * time.Second,
			}

			return fmt.Errorf("error creating test request, fallback to no proxy: %v", err)
		}

		res, err := c.httpClient.Do(req)
		if err != nil || res.StatusCode >= 400 {
			// If the proxy request fails or returns error, fallback to no proxy
			c.httpClient = &http.Client{
				Transport: nil,
				Timeout:   15 * time.Second,
			}

			return fmt.Errorf("proxy failed (status: %v)", res.StatusCode)
		}
		defer res.Body.Close()

	case "socks5":
		// SOCKS5 Proxy
		var auth *proxy.Auth
		if proxyUser != "" && proxyPass != "" {
			auth = &proxy.Auth{
				User:     proxyUser,
				Password: proxyPass,
			}
		}

		// Create SOCKS5 dialer
		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
		if err != nil {
			// Fallback to no proxy
			c.httpClient = &http.Client{
				Transport: nil,
				Timeout:   15 * time.Second,
			}
			return fmt.Errorf("error creating SOCKS5 dialer!")
		}

		// Set the transport to use the SOCKS5 dialer
		transport := &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := dialer.Dial(network, addr)
				if err != nil {
					// Log the error and fallback to no proxy
					c.httpClient = &http.Client{
						Transport: nil,
						Timeout:   15 * time.Second,
					}
					return nil, err
				}
				return conn, nil
			},
		}
		c.httpClient = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}

	default:
		return fmt.Errorf("Unsupported proxy scheme: %s", parsedURL.Scheme)
	}

	return nil
}

func (c *Client) checkIp() (map[string]interface{}, error) {
	result, err := c.makeRequest("GET", "https://ipinfo.io/json", nil)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	return result, nil
}

func (c *Client) makeRequest(method, url string, payload interface{}) (map[string]interface{}, error) {
	res, err := c.processRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("Request error: %v", err)
	}

	result, err := handleResponse(res)
	if err != nil {
		return nil, fmt.Errorf("Handle response error: %v", err)
	}

	return result, nil
}

func (c *Client) processRequest(method, url string, payload interface{}) ([]byte, error) {
	var reqBody []byte
	var err error

	if payload != nil {
		reqBody, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if !strings.Contains(url, "raw.githubusercontent.com") {
		c.setHeader(req)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		status := http.StatusText(res.StatusCode)
		if status == "" {
			status = "Unknown Error"
		}

		return nil, fmt.Errorf("error status: %v, error message: %s", res.StatusCode, status)
	}

	return io.ReadAll(res.Body)
}

func handleResponse(resBody []byte) (map[string]interface{}, error) {
	var responseObject interface{}
	err := json.Unmarshal(resBody, &responseObject)
	if err != nil {
		return nil, fmt.Errorf("Parse JSON response error: %v", err)
	}

	var result map[string]interface{}

	switch value := responseObject.(type) {
	case []interface{}:
		result = make(map[string]interface{})
		for i, v := range value {
			result[fmt.Sprintf("%d", i)] = v
		}
	case map[string]interface{}:
		result = value
	default:
		return nil, fmt.Errorf("Unsupported response type: %T", value)
	}

	return result, nil
}
