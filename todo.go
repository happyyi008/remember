package main

import (
	"fmt"
	"text/tabwriter"
	"time"
)

type Status int

const (
	New Status = iota
	InProgress
	Done
)

func (s Status) String() string {
	names := []string{
		"New",
		"In Progress",
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

func (todo *Todo) Print(index int, w *tabwriter.Writer) {
	fmt.Fprintf(w, " %d.\t%s\t\t%s\t%s\n", index,
		todo.Message,
		todo.Timestamp.Format("(15:04, Mon. Jan 2, 2006)"),
		todo.TodoStatus.String(),
	)
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
