package tools

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func PrintLogo() {
	levelColor := color.New(color.FgCyan)
	levelColor.Println(`
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

	levelColor.Println("ρσωєяє∂ ву: ѕкιвι∂ι ѕιgмα ¢σ∂є")
}

func Logger(level, message string) {
	message = strings.ReplaceAll(message, "\n", "")
	message = strings.ReplaceAll(message, "\r", "")

	level = strings.ToLower(level)
	var levelColor *color.Color

	switch level {
	case "info":
		levelColor = color.New(color.FgWhite)
	case "error":
		levelColor = color.New(color.FgRed)
	case "success":
		levelColor = color.New(color.FgGreen)
	case "warning":
		levelColor = color.New(color.FgYellow)
	default:
		levelColor = color.New(color.FgWhite)
	}

	if level == "input" {
		levelColor.Printf("[+] %s", message)
	} else if level == "error" || level == "warning" {
		levelColor.Println(fmt.Sprintf("[!] %s", message))
	} else if _, err := strconv.Atoi(level); err == nil {
		levelColor.Println(fmt.Sprintf("[%s] %s", level, message))
	} else {
		levelColor.Println(fmt.Sprintf("[*] %s", message))
	}
}
