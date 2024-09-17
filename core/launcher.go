package core

import (
	"MajorBot/helper"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

type Account struct {
	QueryId        string
	UserId         int
	Username       string
	FirstName      string
	LastName       string
	AuthDate       string
	Hash           string
	AllowWriteToPm bool
	LanguageCode   string
	QueryData      string
}

func getAccountFromQuery(account *Account) {
	// Parsing Query To Get Username
	value, err := url.ParseQuery(account.QueryData)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to parse query: %v", err.Error()))
		return
	}

	if len(value.Get("query_id")) > 0 {
		account.QueryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.AuthDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.Hash = value.Get("hash")
	}

	userParam := value.Get("user")

	// Mendekode string JSON
	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		panic(err)
	}

	// Mengambil ID dan username dari hasil decode
	userIDFloat, ok := userData["id"].(float64)
	if !ok {
		helper.PrettyLog("error", "Failed to convert ID to float64")
		return
	}

	account.UserId = int(userIDFloat)

	// Ambil username
	username, ok := userData["username"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get username")
		return
	}
	account.Username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get first_name")
		return
	}
	account.FirstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get last_name")
		return
	}
	account.LastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get language_code")
		return
	}
	account.LanguageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		helper.PrettyLog("error", "Failed to get allows_write_to_pm")
		return
	}
	account.AllowWriteToPm = allowWriteToPm
}

func processBotForAccount(account *Account, walletAddress string, swipeCoins int, holdCoins int, isBindWallet bool) {
	helper.PrettyLog("info", fmt.Sprintf("%s |Starting Bot...", account.Username))

	launchBot(account, swipeCoins, holdCoins, isBindWallet, walletAddress)

	helper.PrettyLog("info", fmt.Sprintf("%s | Launch Bot Finished...", account.Username))

	if !isBindWallet {
		randomSleep := helper.RandomNumber(config.Int("RANDOM_SLEEP.MIN"), config.Int("RANDOM_SLEEP.MAX"))
		helper.PrettyLog("info", fmt.Sprintf("%s | Sleep %vs..", account.Username, randomSleep))
		time.Sleep(time.Duration(randomSleep) * time.Second)
	}
}

func ProcessBot(config *config.Config) {
	queryPath := "./query.txt"
	maxThread := config.Int("MAX_THREAD")
	swipeCoins := helper.RandomNumber(config.Int("SWIPE_COINS.MIN"), config.Int("SWIPE_COINS.MAX"))
	holdCoins := helper.RandomNumber(config.Int("HOLD_COINS.MIN"), config.Int("HOLD_COINS.MAX"))

	queryData := helper.ReadFileTxt(queryPath)
	if queryData == nil {
		helper.PrettyLog("error", "Query data not found")
		return
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Query Data Detected", len(queryData)))

	var choice int
	flagArg := flag.Int("c", 0, "Input Choice With Flag -c, 1 = Auto Completing All Task (Unlimited Loop), 2 = Connect Wallet")

	// Parse flag dari command line
	flag.Parse()

	if *flagArg > 2 {
		helper.PrettyLog("error", "Invalid Flag Choice")
	} else if *flagArg != 0 {
		choice = *flagArg
	}

	if choice == 0 {
		helper.PrettyLog("1", "Auto Completing All Task (Unlimited Loop)")
		helper.PrettyLog("2", "Connect Wallet")

		helper.PrettyLog("input", "Select Your Choice: ")

		_, err := fmt.Scan(&choice)
		if err != nil {
			helper.PrettyLog("error", "Selection Invalid")
			return
		}
	}

	var walletAddress []string

	if choice == 2 {
		walletAddress = helper.ReadFileTxt("./wallet_address.txt")
		if walletAddress == nil {
			helper.PrettyLog("error", "Wallet Address Data Not Found")
			return
		}

		helper.PrettyLog("info", fmt.Sprintf("%v Wallet Address Detected", len(walletAddress)))

		if len(walletAddress) != len(queryData) {
			helper.PrettyLog("error", fmt.Sprintf("Wallet Address Count (%v) Must Match With Query Data Count (%v)", len(walletAddress), len(queryData)))
			return
		}

	}

	helper.PrettyLog("info", "Start Processing Account...")

	time.Sleep(3 * time.Second)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxThread)

	processAccount := func(index int, query string) {
		defer wg.Done()
		semaphore <- struct{}{}

		account := &Account{QueryData: query}
		getAccountFromQuery(account)

		isBindWallet := (choice == 2)
		wallet := ""
		if isBindWallet {
			wallet = walletAddress[index]
		}

		processBotForAccount(account, wallet, swipeCoins, holdCoins, isBindWallet)

		<-semaphore
	}

	switch choice {
	case 1:
		for {
			for j, query := range queryData {
				wg.Add(1)
				go processAccount(j, query)
			}
			wg.Wait() // Tunggu semua goroutine selesai
		}
	case 2:
		for j, query := range queryData {
			wg.Add(1)
			go processAccount(j, query)
		}
		wg.Wait() // Tunggu semua goroutine selesai, lalu program selesai
	}
}
