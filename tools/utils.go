package tools

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func ExitRecover() {
	if r := recover(); r != nil {
		Logger("error", fmt.Sprintf("%v", r))
		Logger("info", "Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func HandleRecover() {
	if r := recover(); r != nil {
		Logger("error", fmt.Sprintf("%v", r))
	}
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func InputChoice(length int) int {
	var choice int

	Logger("input", "Select Choice: ")

	fmt.Scan(&choice)
	if choice <= 0 || choice > length {
		Logger("error", "Invalid selection")
		return 0
	}

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	return choice
}
