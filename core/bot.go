package core

import (
	"MajorBot/tools"
	"fmt"
	"time"

	"github.com/gookit/config/v2"
)

func (c *Client) autoCompleteTask() int {
	swipeCoins := tools.RandomNumber(config.Int("SWIPE_COINS.MIN"), config.Int("SWIPE_COINS.MAX"))
	holdCoins := tools.RandomNumber(config.Int("HOLD_COINS.MIN"), config.Int("HOLD_COINS.MAX"))

	var points int

	token, err := c.getToken(c.account.queryData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get token: %v", err))
		return points
	}

	if token != "" {
		c.accessToken = fmt.Sprintf("Bearer %s", token)
	} else {
		tools.Logger("error", "Token not found")
		return points
	}

	userData, err := c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return points
	}

	points = int(userData["rating"].(float64))

	tools.Logger("success", fmt.Sprintf("| %s | Points: %v", c.account.username, points))

	if squadId, exits := userData["squad_id"].(float64); exits {
		if int64(squadId) != 2414599412 {
			c.leaveSquad()
			c.joinSquad()
		}
	} else {
		c.joinSquad()
	}

	squadInfo, err := c.getSquad()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get squad info: %v", c.account.username, err))
	}

	if squadInfo != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Squad: %s | Points: %v | Member: %v", c.account.username, squadInfo["name"].(string), int(squadInfo["rating"].(float64)), int(squadInfo["members_count"].(float64))))
	}

	dailyVisit, err := c.visit()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to visit: %v", c.account.username, err))
	}

	if dailyVisit != nil {
		if dailyVisit["is_increased"].(bool) {
			tools.Logger("success", fmt.Sprintf("| %s | Daily Streak: %v", c.account.username, int(dailyVisit["streak"].(float64))))
		}
	}

	var allTask []map[string]interface{}

	dailyTask, err := c.getDailyTask()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get daily task: %v", c.account.username, err))
	}

	anotherTask, err := c.getAnotherTask()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get another task: %v", c.account.username, err))
	}

	for _, value := range dailyTask {
		task := value.(map[string]interface{})
		allTask = append(allTask, task)
	}

	for _, value := range anotherTask {
		task := value.(map[string]interface{})
		allTask = append(allTask, task)
	}

	for _, task := range allTask {
		if !task["is_completed"].(bool) {
			completingTask, err := c.completingTask(int(task["id"].(float64)), task["title"].(string))
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Failed to completing task: %v", c.account.username, err))
			}

			if completingTask != nil {
				if completingTask["is_completed"].(bool) {
					tools.Logger("success", fmt.Sprintf("| %s | Claim Task: %s Completed | Award: %v | Sleep 15s Before Completing Next Task...", c.account.username, task["title"].(string), int(task["award"].(float64))))
				} else {
					tools.Logger("error", fmt.Sprintf("| %s | Claim Task: %v Failed | Sleep 15s Before Completing Next Task...", c.account.username, task["title"].(string)))
				}
			}

			time.Sleep(15 * time.Second)
		}
	}

	isSwipeCoins, err := c.checkSwipeCoins()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check swipe coins: %v", c.account.username, err))
	}

	if isSwipeCoins {
		tools.Logger("success", fmt.Sprintf("| %s | Start Playing Swipe Coins After 5s...", c.account.username))

		time.Sleep(5 * time.Second)

		isPlaySwipeCoins, err := c.playSwipeCoins(swipeCoins)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to play swipe coins: %v", c.account.username, err))
		}

		if isPlaySwipeCoins {
			tools.Logger("success", fmt.Sprintf("| %s | Playing Swipe Coins Completed | Award: %v", c.account.username, swipeCoins))
		}
	}

	isHoldCoins, err := c.checkHoldCoins()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check hold coins: %v", c.account.username, err))
	}

	if isHoldCoins {
		tools.Logger("success", fmt.Sprintf("| %s | Start Playing Hold Coins After 5s...", c.account.username))

		time.Sleep(5 * time.Second)

		isPlayHoldCoins, err := c.playHoldCoins(holdCoins)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to play hold coins: %v", c.account.username, err))
		}

		if isPlayHoldCoins {
			tools.Logger("success", fmt.Sprintf("| %s | Playing Swipe Coins Completed | Award: %v", c.account.username, holdCoins))
		}
	}

	solvePuzzle, err := c.getSolvePuzzle()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get solve durov puzzle: %v", c.account.username, err))
	}

	if solvePuzzle != nil {
		if answer, exits := solvePuzzle["tasks"].([]interface{}); exits && solvePuzzle["date"].(string) == time.Now().UTC().Format("2006-01-02") {
			for _, item := range answer {
				if taskMap, ok := item.(map[string]interface{}); ok {

					isDurovPuzzle, err := c.checkDurovPuzzle()
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Failed to check durov puzzle: %v", c.account.username, err))
					}

					if isDurovPuzzle {
						playDurovPuzzle, err := c.playDurovPuzzle(taskMap)
						if err != nil {
							tools.Logger("error", fmt.Sprintf("| %s | Failed to play durov puzzle: %v", c.account.username, err))
						}

						if playDurovPuzzle != nil {
							if _, exits := playDurovPuzzle["correct"].(map[string][]int); exits {
								tools.Logger("success", fmt.Sprintf("| %s | Play Solve Durov Puzzle Correct... %v", c.account.username))
							}
						}

					}
				}
			}
		}
	}

	isRoulette, err := c.checkRoulette()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check roulette: %v", c.account.username, err))
	}

	if isRoulette {
		tools.Logger("success", fmt.Sprintf("| %s | Start Playing Roulette Coins After 5s...", c.account.username))

		time.Sleep(5 * time.Second)

		playRoulette, err := c.playRoulette()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to play roulette: %v", c.account.username, err))
		}

		if playRoulette != nil {
			if award, exits := playRoulette["rating_award"].(float64); exits {
				tools.Logger("success", fmt.Sprintf("| %s | Play Roulette Completed | Award: %v", c.account.username, int(award)))
			}
		}
	}

	userData, err = c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return points
	}

	if userData != nil {
		points = int(userData["rating"].(float64))
		tools.Logger("success", fmt.Sprintf("| %s | Update Points: %v", c.account.username, points))
	}

	return points
}

func (c *Client) connectWallet() {
	token, err := c.getToken(c.account.queryData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get token: %v", c.account.username, err))
		return
	}

	if token != "" {
		c.accessToken = fmt.Sprintf("Bearer %s", token)
	} else {
		tools.Logger("error", "Token not found")
		return
	}

	err = c.bindWallet()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to connect wallet: %v", c.account.username, err))
	} else {
		tools.Logger("success", fmt.Sprintf("| %s | Successfully connect wallet address: %s", c.account.username, c.account.walletAddress))
	}
}
