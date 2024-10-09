package core

import (
	"MajorBot/tools"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

func (account *Account) parsingQueryData() {
	value, err := url.ParseQuery(account.queryData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse query data: %s", err))
	}

	if len(value.Get("query_id")) > 0 {
		account.queryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.authDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.hash = value.Get("hash")
	}

	userParam := value.Get("user")

	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse user data: %s", err))
	}

	userId, ok := userData["id"].(float64)
	if !ok {
		tools.Logger("error", "Failed to convert ID to float64")
	}

	account.userId = int(userId)

	username, ok := userData["username"].(string)
	if !ok {
		tools.Logger("error", "Failed to get username from query")
		return
	}

	account.username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get first name from query")
	}

	account.firstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get last name from query")
	}
	account.lastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		tools.Logger("error", "Failed to get language code from query")
	}
	account.languageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		tools.Logger("error", "Failed to get allows write to pm from query")
	}

	account.allowWriteToPm = allowWriteToPm
}

func (account *Account) worker(wg *sync.WaitGroup, semaphore *chan struct{}, totalPointsChan *chan int, index int, query string, proxyList []string, walletList []string) {
	defer wg.Done()
	*semaphore <- struct{}{}

	var points int
	var proxy string

	if len(proxyList) > 0 {
		proxy = proxyList[index%len(proxyList)]
	}

	tools.Logger("info", fmt.Sprintf("| %s | Starting Bot...", account.username))

	setDns(&net.Dialer{})

	client := Client{
		account: *account,
		proxy:   proxy,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	if len(client.proxy) > 0 {
		err := client.setProxy()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to set proxy: %v", account.username, err))
		} else {
			tools.Logger("success", fmt.Sprintf("| %s | Proxy Successfully Set...", account.username))
		}
	}

	infoIp, err := client.checkIp()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to check ip: %v", err))
	}

	if infoIp != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Ip: %s | City: %s | Country: %s | Provider: %s", account.username, infoIp["ip"].(string), infoIp["city"].(string), infoIp["country"].(string), infoIp["org"].(string)))
	}

	if len(walletList) > 0 {
		account.walletAddress = walletList[index]
		client.connectWallet()
	} else {
		points = client.autoCompleteTask()
		*totalPointsChan <- points
	}

	<-*semaphore
}

func LaunchBot(selectedTools int) {
	defer tools.HandleRecover()
	queryPath := "configs/query.txt"
	proxyPath := "configs/proxy.txt"
	walletAddressPath := "configs/wallet_address.txt"
	maxThread := config.Int("MAX_THREAD")
	isUseProxy := config.Bool("USE_PROXY")

	queryData, err := tools.ReadFileTxt(queryPath)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Query Data Not Found: %s", err))
		return
	}

	tools.Logger("info", fmt.Sprintf("%v Query Data Detected", len(queryData)))

	var wg sync.WaitGroup
	var semaphore chan struct{}
	var proxyList, walletList []string

	if isUseProxy {
		proxyList, err = tools.ReadFileTxt(proxyPath)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("Proxy Data Not Found: %s", err))
		}

		tools.Logger("info", fmt.Sprintf("%v Proxy Detected", len(proxyList)))
	}

	if selectedTools == 2 {
		walletList, err = tools.ReadFileTxt(walletAddressPath)
		if err != nil {
			tools.Logger("error", "Wallet Address Data Not Found")
			return
		}

		tools.Logger("info", fmt.Sprintf("%v Wallet Address Detected", len(walletList)))

		if len(walletList) != len(queryData) {
			tools.Logger("error", fmt.Sprintf("Wallet Address Count (%v) Must Match With Query Data Count (%v)", len(walletList), len(queryData)))
			return
		}
	}

	if maxThread > len(queryData) {
		semaphore = make(chan struct{}, len(queryData))
	} else {
		semaphore = make(chan struct{}, maxThread)
	}

	switch selectedTools {
	case 1:
		for {
			totalPointsChan := make(chan int, len(queryData))

			for index, query := range queryData {
				wg.Add(1)
				account := &Account{
					queryData: query,
				}

				account.parsingQueryData()

				go account.worker(&wg, &semaphore, &totalPointsChan, index, query, proxyList, walletList)
			}

			go func() {
				wg.Wait()
				close(totalPointsChan)
			}()

			var totalPoints int

			for points := range totalPointsChan {
				totalPoints += points
			}

			tools.Logger("success", fmt.Sprintf("Total Points All Account: %v", totalPoints))

			randomSleep := tools.RandomNumber(config.Int("RANDOM_SLEEP.MIN"), config.Int("RANDOM_SLEEP.MAX"))

			tools.Logger("info", fmt.Sprintf("Launch Bot Finished | Sleep %vs Before Next Lap...", randomSleep))

			time.Sleep(time.Duration(randomSleep) * time.Second)
		}
	case 2:
		for index, query := range queryData {
			wg.Add(1)
			account := &Account{
				queryData: query,
			}

			account.parsingQueryData()

			go account.worker(&wg, &semaphore, nil, index, query, proxyList, walletList)
		}
		wg.Wait()
	}
}
