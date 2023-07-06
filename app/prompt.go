package app

import (
	"bufio"
	"fmt"
)

func prompt(reader *bufio.Reader, label string) (string, error) {
	fmt.Printf("%s: ", label)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if len(input) > 0 {
		return input, nil
	}

	return "", fmt.Errorf("Invalid input received for prompt %q", label)
}
