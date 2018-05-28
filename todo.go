package main

import (
	"fmt"
	"time"
)

type Status int

const (
	InProgress Status = iota
	New
	Done
)

func (s Status) String() string {
	names := []string{
		"In Progress",
		"New",
		"Finished",
	}
	return names[s]
}

type Todo struct {
	Message    string    `json:message`
	Timestamp  time.Time `json:timestamp`
	TodoStatus Status    `json:status`
}

func NewTodo(message string, timestamp time.Time) *Todo {
	return &Todo{Message: message, Timestamp: timestamp, TodoStatus: New}
}

func (todo Todo) String() string {
	return fmt.Sprintf("\t%s\t%s",
		todo.Message,
		todo.TodoStatus,
	)
}

// The type def and the follwing three funcs are for sorting todos
// by timestamp. Latest todo at the top.
type TodoByTimestamp []*Todo

func (t TodoByTimestamp) Len() int {
	return len(t)
}

func (t TodoByTimestamp) Less(i, j int) bool {
	if t[j].TodoStatus > t[i].TodoStatus {
		return true
	}
	if t[j].TodoStatus < t[i].TodoStatus {
		return false
	}
	return t[j].Timestamp.Before(t[i].Timestamp)
}

func (t TodoByTimestamp) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
