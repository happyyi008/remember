package main

import (
	"github.com/c-bata/go-prompt"
	"strings"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "ls", Description: "Print your todo list"},
		{Text: "rm", Description: "Delete from todo list"},
		{Text: "set", Description: "Set status of Todo"},
		{Text: "new", Description: "Add a new task"},
		{Text: "help", Description: "Get detailed description of each command"},
	}

	currentLine := strings.Split(d.CurrentLineBeforeCursor(), " ")

	// get suggestion for set instead
	if currentLine[0] == "set" && len(currentLine) == 2 {
		s = []prompt.Suggest{
			{Text: "start", Description: "Start the task"},
			{Text: "done", Description: "Finish task"},
			{Text: "restart", Description: "Reset status to New"},
		}
	} else if len(currentLine) >= 2 {
		s = []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
