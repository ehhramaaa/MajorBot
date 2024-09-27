package core

import (
	"MajorBot/tools"

	"github.com/mileusna/useragent"
)

func generateRandomUserAgent() (string, string) {
	userAgents, err := tools.ReadFileTxt("./configs/useragent.txt")
	if err != nil {
		return "", ""
	}

	userAgent := userAgents[tools.RandomNumber(0, len(userAgents))]

	os := useragent.Parse(userAgent).OS

	return userAgent, os
}
