package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
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

type Remember struct {
	Todos []Todo `json:"todoList"`
}

func (r *Remember) addTodo(message []string) {
	messageStr := strings.Join(message, " ")
	if messageStr == "" {
		reader := bufio.NewReader(os.Stdin)
		messageStr, _ = reader.ReadString('\n')
		messageStr = strings.TrimSuffix(messageStr, "\n") // trim newline from the reader
	}
	todo := *NewTodo(messageStr, time.Now())
	r.Todos = append([]Todo{todo}, r.Todos...) // hack to prepend new todo
	log.Debug("added new todo")
	r.writeToFile()
}

// TODO have different ways of printing
func (r *Remember) listTodo() {
	if len(r.Todos) == 0 {
		fmt.Println("You have no todos in your list.")
		return
	}

	sort.Sort(TodoByTimestamp(r.Todos))
	if sort.IsSorted(TodoByTimestamp(r.Todos)) {
		fmt.Println("this guy is sorted")
	}

	fmt.Println("Your list of Todos:")
	// TODO have max line length based on terminal width and break line
	// so long todos don't break printing
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for i, todo := range r.Todos {
		fmt.Fprintf(w, " %d.\t%s\t\t%s\n", i+1,
			todo.Message,
			todo.Timestamp.Format("(15:04, Mon, Jan 2, 2006)"))
	}
	w.Flush()
}

func (r *Remember) deleteTodo(args []string) {
	if len(args) != 2 {
		fmt.Println("Invalid command: missing arguments")
		return
	}
	indexToDelete, err := strconv.Atoi(args[1])
	if err != nil || indexToDelete > len(r.Todos) {
		fmt.Println("Error: Not a valid index to delete.")
		return
	}
	r.Todos = append(r.Todos[:indexToDelete-1], r.Todos[indexToDelete:]...)
	r.writeToFile()
}

func (r *Remember) writeToFile() {
	res, err := json.Marshal(r)
	checkErr(err)
	write(res)
}
