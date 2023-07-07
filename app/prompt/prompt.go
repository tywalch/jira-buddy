package prompt

import (
	"bufio"
	"fmt"
	MultipleChoice "github.com/thewolfnl/go-multiplechoice"
	"strings"
)

type PromptReader struct {
	*bufio.Reader
}

type Prompter interface {
	GetString(label string) (string, error)
	PickString(label string, options []StringPickerOption) (string, error)
	PickNumeric(label string, options []NumericPickerOption) (int, error)
}

func (reader PromptReader) GetString(label string) (string, error) {
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

type StringPickerOption struct {
	Name  string
	Value string
}

func (reader PromptReader) PickString(label string, options []StringPickerOption) (string, error) {
	text := fmt.Sprintf("%s: ", label)
	names := make([]string, len(options))
	for i, option := range options {
		names[i] = option.Name
	}

	selectedName := MultipleChoice.Selection(text, names)
	selection := ""
	for _, option := range options {
		if option.Name == selectedName {
			selection = option.Value
		}
	}

	return strings.TrimSpace(selection), nil
}

type NumericPickerOption struct {
	Name  string
	Value int
}

func (reader PromptReader) PickNumeric(label string, options []NumericPickerOption) (int, error) {
	text := fmt.Sprintf("%s: ", label)
	names := make([]string, len(options))
	for i, option := range options {
		names[i] = option.Name
	}

	selectedName := MultipleChoice.Selection(text, names)
	var selection int
	for _, option := range options {
		if option.Name == selectedName {
			selection = option.Value
		}
	}

	return selection, nil
}
