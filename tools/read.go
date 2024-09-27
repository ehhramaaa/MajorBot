package tools

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFileTxt(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s: %v", filepath, err)
	}

	defer file.Close()

	var value []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		value = append(value, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v", filepath, err)
	}

	return value, nil
}
