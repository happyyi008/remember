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
		fmt.Fprintf(w, " %d.\t%s\t\t%s\n", i+1,
			todo.Message,
			todo.Timestamp.Format("(15:04, Mon, Jan 2, 2006)"))
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
}

func (r *Remember) deleteTodo(args []string) {
	if len(args) < 1 {
		fmt.Println("Invalid command: missing delete index")
		return
	}

	// sort the indice so we will delete entries in order
	sort.Strings(args)
	log.Debugf("sorted args %+v", args)

	for offset, index := range args {
		indexToDelete, err := strconv.Atoi(index)
		indexToDelete -= offset

		if err != nil || indexToDelete > len(r.Todos) {
			fmt.Println("Error: Not a valid index to delete.")
			return
		}
		r.Todos = append(r.Todos[:indexToDelete-1], r.Todos[indexToDelete:]...)
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
