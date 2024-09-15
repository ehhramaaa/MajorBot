package core

import (
	"MajorBot/helper"
	"fmt"
	"net/http"
	"time"
)

func launchBot(account *Account, swipeCoins int, holdCoins int, isBindWallet bool, walletAddress string) {
	client := &Client{
		username:   account.Username,
		httpClient: &http.Client{},
	}

	token := client.getToken(account)

	if len(token["access_token"].(string)) > 0 {
		client.authToken = fmt.Sprintf("Bearer %s", token["access_token"].(string))
	} else {
		helper.PrettyLog("error", "Failed To Get Token")
		return
	}

	if isBindWallet {
		connectWallet := client.bindWallet(walletAddress)
		fmt.Println(connectWallet)
		return
	}

	userData := client.getUserInfo(account)

	if _, exits := userData["username"].(string); exits && userData["username"].(string) == account.Username {
		helper.PrettyLog("success", fmt.Sprintf("%s | Points: %v", client.username, int(userData["rating"].(float64))))
	}

	_, ok := userData["squad_id"].(float64)
	if !ok {
		client.joinSquad()
	} else if int(userData["squad_id"].(float64)) != 2414599412 {
		client.leaveSquad()
		client.joinSquad()
	}

	getSquad := client.getSquad()

	if _, exits := getSquad["name"].(string); exits {
		helper.PrettyLog("success", fmt.Sprintf("%s | Squad: %s | Points: %v | Member: %v", client.username, getSquad["name"].(string), int(getSquad["rating"].(float64)), int(getSquad["members_count"].(float64))))
	}

	// dailyStreak := client.streak()
	// fmt.Println(dailyStreak)

	dailyVisit := client.visit()
	if dailyVisit["is_increased"].(bool) {
		helper.PrettyLog("success", fmt.Sprintf("%s | Daily Streak: %v", client.username, int(dailyVisit["streak"].(float64))))
	}

	var allTask []map[string]interface{}

	dailyTask := client.getDailyTask()

	anotherTask := client.getAnotherTask()

	allTask = append(allTask, dailyTask...)
	allTask = append(allTask, anotherTask...)

	for _, task := range allTask {
		if !task["is_completed"].(bool) {
			completingTask := client.completingTask(int(task["id"].(float64)), task["title"].(string))
			if completingTask["is_completed"].(bool) {
				helper.PrettyLog("success", fmt.Sprintf("%s | Claim Task: %s Completed | Award: %v | Sleep 15s Before Completing Next Task...", client.username, completingTask["title"].(string), int(completingTask["award"].(float64))))
			} else {
				helper.PrettyLog("error", fmt.Sprintf("%s | Claim Task: %v Failed | Sleep 15s Before Completing Next Task...", client.username, task["title"].(string)))
			}

			time.Sleep(15 * time.Second)
		}
	}

	isSwipeCoins := client.checkSwipeCoins()
	if _, exits := isSwipeCoins["success"].(bool); exits && isSwipeCoins["success"].(bool) {
		helper.PrettyLog("success", fmt.Sprintf("%s | Start Playing Swipe Coins After 5s...", client.username))
		time.Sleep(5 * time.Second)
		playSwipeCoins := client.playSwipeCoins(swipeCoins)

		if _, exits := playSwipeCoins["success"].(bool); exits && playSwipeCoins["success"].(bool) {
			helper.PrettyLog("success", fmt.Sprintf("%s | Playing Swipe Coins Completed | Award: %v", client.username, swipeCoins))
		}
	}

	isHoldCoins := client.checkHoldCoins()
	if _, exits := isHoldCoins["success"].(bool); exits && isHoldCoins["success"].(bool) {
		helper.PrettyLog("success", fmt.Sprintf("%s | Start Playing Hold Coins After 5s...", client.username))
		time.Sleep(5 * time.Second)
		playHoldCoins := client.playHoldCoins(holdCoins)

		if _, exits := playHoldCoins["success"].(bool); exits && playHoldCoins["success"].(bool) {
			helper.PrettyLog("success", fmt.Sprintf("%s | Playing Swipe Coins Completed | Award: %v", client.username, holdCoins))
		}
	}

	isGetSolvePuzzle := client.getSolvePuzzle()
	if answer, exits := isGetSolvePuzzle["tasks"].([]interface{}); exits && isGetSolvePuzzle["date"].(string) == time.Now().UTC().Format("2006-01-02") {
		for _, item := range answer {
			if taskMap, ok := item.(map[string]interface{}); ok {
				checkDurovPuzzle := client.checkDurovPuzzle()
				if _, exits := checkDurovPuzzle["success"].(bool); exits && checkDurovPuzzle["success"].(bool) {
					playDurovPuzzle := client.playDurovPuzzle(taskMap)
					if _, exits := playDurovPuzzle["correct"].(map[string][]int); exits {
						helper.PrettyLog("success", fmt.Sprintf("%s | Play Solve Durov Puzzle Correct... %v", client.username))
					}
				}
			}
		}
	}
}
