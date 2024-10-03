package main

import (
	"MajorBot/core"
	"MajorBot/tools"
	"flag"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("configs/config.yml")
	if err != nil {
		panic(err)
	}

	tools.PrintLogo()

	var selectedTools int

	flagArg := flag.Int("c", 0, "Input Choice With Flag -c, 1 = Auto Completing All Task (Unlimited Loop), 2 = Connect Wallet")
	tools.Logger("1", "Auto Completing Task (Unlimited Loop)")
	tools.Logger("2", "Connect Wallet")

	flag.Parse()

	if *flagArg > 2 {
		tools.Logger("error", "Invalid Flag Choice")
	} else if *flagArg != 0 {
		selectedTools = *flagArg
	} else {
		selectedTools = tools.InputChoice(2)
		if selectedTools == 0 {
			tools.Logger("error", "Invalid Choice")
			return
		}
	}

	core.LaunchBot(selectedTools)
}
