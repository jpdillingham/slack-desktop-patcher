package util

import (
	"bufio"
	"fmt"
	"os"
)

func PromptForInput(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(prompt)

	input := ""

	for scanner.Scan() {
		input = scanner.Text()
		break
	}

	return input
}
