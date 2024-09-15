package core

import (
	"MajorBot/helper"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	username   string
	authToken  string
	httpClient *http.Client
}

func (c *Client) makeRequest(method string, endpoint string, jsonBody interface{}) ([]byte, error) {
	var fullURL string
	if endpoint == "https://raw.githubusercontent.com/dancayairdrop/blum/main/durov.json" {
		fullURL = endpoint
	} else {
		fullURL = "https://major.bot/api" + endpoint
	}

	// Convert body to JSON
	var reqBody []byte
	var err error
	if jsonBody != nil {
		reqBody, err = json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
	}

	// Create new request
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	if endpoint == "https://raw.githubusercontent.com/dancayairdrop/blum/main/durov.json" {
		var header = map[string]string{
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"accept-language":           "en-US,en;q=0.9,id;q=0.8",
			"cache-control":             "max-age=0",
			"if-none-match":             "W/\"71f09c6fd3ed7a37cbd730c88157a3038913475c5af18078f9ecf5ad565cc800\"",
			"priority":                  "u=0, i",
			"sec-ch-ua":                 "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"",
			"sec-ch-ua-mobile":          "?0",
			"sec-ch-ua-platform":        "\"Windows\"",
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-site":            "none",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
		}

		for key, value := range header {
			req.Header.Set(key, value)
		}
	} else {
		setHeader(req, c.authToken)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 status code
	if resp.StatusCode >= 400 {
		// Read the response body to include in the error message
		bodyBytes, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, fmt.Errorf("error status: %v, and failed to read body: %v", resp.StatusCode, bodyErr)
		}
		return nil, fmt.Errorf("error status: %v, error message: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// Login
func (c *Client) getToken(account *Account) map[string]interface{} {
	payload := map[string]string{
		"init_data": account.QueryData,
	}

	req, err := c.makeRequest("POST", "/auth/tg/", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to login: %v", c.username, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Visit
func (c *Client) visit() map[string]interface{} {
	res, err := c.makeRequest("POST", "/user-visits/visit/?", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to visit: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Streak
// func (c *Client) streak() map[string]interface{} {
// 	res, err := c.makeRequest("GET", "/user-visits/streak/", nil)
// 	if err != nil {
// 		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to daily streak: %v", c.username, err))
// 		return nil
// 	}

// 	result, err := handleResponseMap(res)
// 	if err != nil {
// 		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
// 		return nil
// 	}

// 	return result
// }

// Get User Info
func (c *Client) getUserInfo(account *Account) map[string]interface{} {
	req, err := c.makeRequest("GET", fmt.Sprintf("/users/%v/", account.UserId), nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to get user info: %v", c.username, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Join Squad
func (c *Client) joinSquad() map[string]interface{} {
	res, err := c.makeRequest("POST", "/squads/2414599412/join/?", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to join squad: %v", c.username, err))
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Leave Squad
func (c *Client) leaveSquad() map[string]interface{} {
	res, err := c.makeRequest("POST", "/squads/leave/", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to leave squad: %v", c.username, err))
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Squad
func (c *Client) getSquad() map[string]interface{} {
	res, err := c.makeRequest("GET", "/squads/2414599412", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to get squad: %v", c.username, err))
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Daily Task List
func (c *Client) getDailyTask() []map[string]interface{} {
	req, err := c.makeRequest("GET", "/tasks/?is_daily=true", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to get daily task: %v", c.username, err))
		return nil
	}

	res, err := handleResponseArray(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Get Another Task List
func (c *Client) getAnotherTask() []map[string]interface{} {
	req, err := c.makeRequest("GET", "/tasks/?is_daily=false", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to get another task: %v", c.username, err))
		return nil
	}

	res, err := handleResponseArray(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Completing Task
func (c *Client) completingTask(taskId int, taskName string) map[string]interface{} {
	payload := map[string]int{
		"task_id": taskId,
	}

	req, err := c.makeRequest("POST", "/tasks/", payload)

	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to completing task %v: %v", c.username, taskName, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Check Swipe Coins
func (c *Client) checkSwipeCoins() map[string]interface{} {
	res, err := c.makeRequest("GET", "/swipe_coin/", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to check swipe coins: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Play Swipe Coins
func (c *Client) playSwipeCoins(swipeCoins int) map[string]interface{} {
	payload := map[string]int{
		"coins": swipeCoins,
	}

	res, err := c.makeRequest("POST", "/swipe_coin/", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to play swipe coins: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Check Swipe Coins
func (c *Client) checkHoldCoins() map[string]interface{} {
	res, err := c.makeRequest("GET", "/bonuses/coins/", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to check hold coins: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Play Hold Coins
func (c *Client) playHoldCoins(holdCoins int) map[string]interface{} {
	payload := map[string]int{
		"coins": holdCoins,
	}

	res, err := c.makeRequest("POST", "/bonuses/coins/", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to play hold coins: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Bind Wallet
func (c *Client) bindWallet(walletAddress string) map[string]interface{} {
	payload := map[string]string{
		"address": walletAddress,
	}

	res, err := c.makeRequest("POST", "/users/address/", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to bind wallet: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Solve Durov Puzzle
func (c *Client) getSolvePuzzle() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://raw.githubusercontent.com/dancayairdrop/blum/main/durov.json", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to get solve puzzle: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Check Durov Puzzle
func (c *Client) checkDurovPuzzle() map[string]interface{} {
	res, err := c.makeRequest("GET", "/durov/", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to check durov puzzle: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Play Durov Puzzle
func (c *Client) playDurovPuzzle(answer map[string]interface{}) map[string]interface{} {
	res, err := c.makeRequest("POST", "/durov/", answer)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed to check durov puzzle: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("%s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}
