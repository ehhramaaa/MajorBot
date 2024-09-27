package core

import "fmt"

const apiUrl = "https://major.bot/api"

// Login
func (c *Client) getToken(query string) (string, error) {
	payload := map[string]string{
		"init_data": query,
	}

	result, err := c.makeRequest("POST", apiUrl+"/auth/tg/", payload)
	if err != nil {
		return "", err
	}

	if token, exits := result["access_token"].(string); exits {
		return token, nil
	} else {
		return "", fmt.Errorf("token not found")
	}
}

// Get User Info
func (c *Client) getUserInfo() (map[string]interface{}, error) {
	result, err := c.makeRequest("GET", fmt.Sprintf("%s/users/%v/", apiUrl, c.account.userId), nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Join Squad
func (c *Client) joinSquad() {
	c.makeRequest("POST", apiUrl+"/squads/2414599412/join/?", nil)
}

// Leave Squad
func (c *Client) leaveSquad() {
	c.makeRequest("POST", apiUrl+"/squads/leave/", nil)
}

// Get Squad
func (c *Client) getSquad() (map[string]interface{}, error) {
	result, err := c.makeRequest("GET", apiUrl+"/squads/2414599412", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Visit
func (c *Client) visit() (map[string]interface{}, error) {
	result, err := c.makeRequest("POST", apiUrl+"/user-visits/visit/?", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get Daily Task List
func (c *Client) getDailyTask() (map[string]interface{}, error) {
	result, err := c.makeRequest("GET", apiUrl+"/tasks/?is_daily=true", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get Another Task List
func (c *Client) getAnotherTask() (map[string]interface{}, error) {
	result, err := c.makeRequest("GET", apiUrl+"/tasks/?is_daily=false", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Completing Task
func (c *Client) completingTask(taskId int, taskName string) (map[string]interface{}, error) {
	payload := map[string]int{
		"task_id": taskId,
	}
	result, err := c.makeRequest("POST", apiUrl+"/tasks/", payload)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Check Swipe Coins
func (c *Client) checkSwipeCoins() (bool, error) {
	result, err := c.makeRequest("GET", apiUrl+"/swipe_coin/", nil)
	if err != nil {
		return false, err
	}

	if success, exits := result["success"].(bool); exits && result["success"].(bool) {
		return success, nil
	} else {
		return false, fmt.Errorf("Success Status Not Exits!")
	}
}

// Play Swipe Coins
func (c *Client) playSwipeCoins(swipeCoins int) (bool, error) {
	payload := map[string]int{
		"coins": swipeCoins,
	}

	result, err := c.makeRequest("POST", apiUrl+"/swipe_coin/", payload)
	if err != nil {
		return false, err
	}

	if success, exits := result["success"].(bool); exits && result["success"].(bool) {
		return success, nil
	} else {
		return false, fmt.Errorf("Success Status Not Exits!")
	}
}

// Check Hold Coins
func (c *Client) checkHoldCoins() (bool, error) {
	result, err := c.makeRequest("GET", apiUrl+"/bonuses/coins/", nil)
	if err != nil {
		return false, err
	}

	if success, exits := result["success"].(bool); exits && result["success"].(bool) {
		return success, nil
	} else {
		return false, fmt.Errorf("Success Status Not Exits!")
	}
}

// Play Hold Coins
func (c *Client) playHoldCoins(holdCoins int) (bool, error) {
	payload := map[string]int{
		"coins": holdCoins,
	}

	result, err := c.makeRequest("POST", apiUrl+"/bonuses/coins/", payload)
	if err != nil {
		return false, err
	}

	if success, exits := result["success"].(bool); exits && result["success"].(bool) {
		return success, nil
	} else {
		return false, fmt.Errorf("Success Status Not Exits!")
	}
}

// Get Solve Durov Puzzle
func (c *Client) getSolvePuzzle() (map[string]interface{}, error) {
	result, err := c.makeRequest("GET", "https://raw.githubusercontent.com/dancayairdrop/blum/main/durov.json", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Check Durov Puzzle
func (c *Client) checkDurovPuzzle() (bool, error) {
	result, err := c.makeRequest("GET", apiUrl+"/durov/", nil)
	if err != nil {
		return false, err
	}

	if success, exits := result["success"].(bool); exits && result["success"].(bool) {
		return success, nil
	} else {
		return false, fmt.Errorf("Success Status Not Exits!")
	}
}

// Play Durov Puzzle
func (c *Client) playDurovPuzzle(answer map[string]interface{}) (map[string]interface{}, error) {
	result, err := c.makeRequest("POST", apiUrl+"/durov/", answer)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Check Roulette
func (c *Client) checkRoulette() (bool, error) {
	result, err := c.makeRequest("GET", apiUrl+"/roulette/", nil)
	if err != nil {
		return false, err
	}

	if success, exits := result["success"].(bool); exits && result["success"].(bool) {
		return success, nil
	} else {
		return false, fmt.Errorf("Success Status Not Exits!")
	}
}

// Play Roulette
func (c *Client) playRoulette() (map[string]interface{}, error) {
	result, err := c.makeRequest("POST", apiUrl+"/roulette/", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Bind Wallet
func (c *Client) bindWallet() error {
	payload := map[string]string{
		"address": c.account.walletAddress,
	}

	_, err := c.makeRequest("POST", apiUrl+"/users/address/", payload)
	if err != nil {
		return err
	} else {
		return nil
	}
}
