package main

import (
	"time"
)

type Todo struct {
	Message   string    `json:message`
	Timestamp time.Time `json:timestamp`
}

func NewTodo(message string, timestamp time.Time) *Todo {
	return &Todo{Message: message, Timestamp: timestamp}
}

// The type def and the follwing three funcs are for sorting todos
// by timestamp. Latest todo at the top.
type TodoByTimestamp []Todo

func (t TodoByTimestamp) Len() int {
	return len(t)
}

func (t TodoByTimestamp) Less(i, j int) bool {
	return t[j].Timestamp.Before(t[i].Timestamp)
}

func (t TodoByTimestamp) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
