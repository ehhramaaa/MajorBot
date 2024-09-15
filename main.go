package main

import (
	"MajorBot/core"
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {
	fmt.Println(`
 /$$      /$$                                   /$$$$$$$              /$$    
| $$$    /$$$                                  | $$__  $$            | $$    
| $$$$  /$$$$  /$$$$$$  /$$  /$$$$$$   /$$$$$$ | $$  \ $$  /$$$$$$  /$$$$$$  
| $$ $$/$$ $$ |____  $$|__/ /$$__  $$ /$$__  $$| $$$$$$$  /$$__  $$|_  $$_/  
| $$  $$$| $$  /$$$$$$$ /$$| $$  \ $$| $$  \__/| $$__  $$| $$  \ $$  | $$    
| $$\  $ | $$ /$$__  $$| $$| $$  | $$| $$      | $$  \ $$| $$  | $$  | $$ /$$
| $$ \/  | $$|  $$$$$$$| $$|  $$$$$$/| $$      | $$$$$$$/|  $$$$$$/  |  $$$$/
|__/     |__/ \_______/| $$ \______/ |__/      |_______/  \______/    \___/  
                  /$$  | $$                                                  
                 |  $$$$$$/                                                  
                  \______/                                                   
`)
	fmt.Println(`ρσωєяє∂ ву : нσℓу¢αη`)

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yml")
	if err != nil {
		panic(err)
	}

	core.ProcessBot(config.Default())
}
