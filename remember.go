package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type Remember struct {
	Todos []Todo `json:"todoList"`
}

// initializes remember struct from file or create a new one
func NewRemember() *Remember {
	content, err := ioutil.ReadFile(RMBFILE) // read previously saved todo
	if err != nil {
		log.Debug("initialize Remember")
		content = []byte(`{"todoList":[]}`)
		log.Debug("created init file")
	}

	remember := &Remember{}
	json.Unmarshal(content, remember)
	log.Debugf("Done init: %+v", remember)
	return remember
}

// TODO have different ways of printing
// print order: in progress -> new -> finished
// sorted by latest timestamp first
func (r *Remember) listTodo() {
	if len(r.Todos) == 0 {
		fmt.Println("You have no todos in your list.")
		return
	}

	sort.Sort(TodoByTimestamp(r.Todos))

	fmt.Println("Your list of Todos:")
	// TODO have max line length based on terminal width and break line
	// so long todos don't break printing
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for i, todo := range r.Todos {
		todo.Print(i+1, w)
	}
	w.Flush()
}

func (r *Remember) addTodo(message []string) {
	messageStr := strings.Join(message, " ")
	if messageStr == "" {
		reader := bufio.NewReader(os.Stdin)
		messageStr, _ = reader.ReadString('\n')
		messageStr = strings.TrimSuffix(messageStr, "\n") // trim newline from the reader
	}
	todo := *NewTodo(messageStr, time.Now())
	r.Todos = append([]Todo{todo}, r.Todos...)
	log.Debug("added new todo")

	r.writeToFile()

	fmt.Println("New task created")
}

func (r *Remember) deleteTodo(args []string) {
	// sort the indices so we will delete entries in order
	sort.Strings(args)
	log.Debugf("sorted args %+v", args)

	numDeleted := 0
	for _, index := range args {
		if index == "" {
			continue
		}
		indexToDelete, err := strconv.Atoi(index)
		indexToDelete -= numDeleted

		if err != nil || indexToDelete > len(r.Todos) {
			fmt.Println("Error: Not a valid index to delete.")
			return
		}
		r.Todos = append(r.Todos[:indexToDelete-1], r.Todos[indexToDelete:]...)
		numDeleted += 1
	}

	r.writeToFile()

	fmt.Printf("Deleted %d tasks\n", numDeleted)
}

func (r *Remember) setStatus(index []string, action string) {
	for _, sidx := range index {
		if sidx == "" {
			continue
		}
		idx, err := strconv.Atoi(sidx)
		checkErr(err)
		idx -= 1

		switch action {
		case "start":
			r.Todos[idx].TodoStatus = InProgress
		case "done":
			r.Todos[idx].TodoStatus = Done
		case "restart":
			r.Todos[idx].TodoStatus = New
		}
	}

	r.writeToFile()
}

func (r *Remember) writeToFile() {
	jsonTodo, err := json.Marshal(r)
	checkErr(err)
	ioutil.WriteFile(RMBFILE, jsonTodo, 0644)
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
		Usage()
		return
	}
}
