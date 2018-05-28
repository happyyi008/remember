package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
	"time"
)

type Remember struct {
	Todos []*Todo `json:"todoList"`
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
	fmt.Fprintf(w, "\n")
	for i, todo := range r.Todos {
		fmt.Fprintf(w, "\t%d. %s\n", i+1, todo)
	}
	fmt.Fprintf(w, "\n")
	w.Flush()
}

func (r *Remember) addTodo(message string) {
	todo := NewTodo(message, time.Now())
	r.Todos = append([]*Todo{todo}, r.Todos...)
	log.Debug("added new todo")

	r.writeToFile()
	fmt.Println("New task created")
}

func (r *Remember) deleteTodo(args []string) {
	// sort the indices so we will delete entries in order
	toDelete, err := sliceAtoi(args)
	checkErr(err)

	sort.Ints(toDelete)
	log.Debugf("sorted args %+v", toDelete)

	numDeleted := 0
	for _, index := range toDelete {
		index -= numDeleted
		err := r.deleteAtIndex(index)
		if err != nil {
			fmt.Println(err)
			return
		}
		numDeleted += 1
	}

	r.writeToFile()
	fmt.Printf("Deleted %d tasks\n", numDeleted)
}

func (r *Remember) deleteAtIndex(index int) error {
	if index > len(r.Todos) {
		log.Debugf("can't delete index: %d", index)
		return errors.New("Not a valid index to delete")
	}
	r.Todos = append(r.Todos[:index-1], r.Todos[index:]...)
	return nil
}

func (r *Remember) setStatus(action string, indices []string) {
	toSet, err := sliceAtoi(indices)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, index := range toSet {
		index -= 1
		switch action {
		case "start":
			r.Todos[index].TodoStatus = InProgress
		case "done":
			r.Todos[index].TodoStatus = Done
		case "restart":
			r.Todos[index].TodoStatus = New
		}
	}

	r.writeToFile()
	fmt.Printf("Updated %d tasks\n", len(toSet))
}

func (r *Remember) writeToFile() {
	jsonTodo, err := json.Marshal(r)
	checkErr(err)
	ioutil.WriteFile(RMBFILE, jsonTodo, 0644)
}

// transforms a slice of strings to a slice of ints and remove empty strings
func sliceAtoi(strSlice []string) ([]int, error) {
	elementCount := 0
	intSlice := make([]int, len(strSlice))
	for _, strIdx := range strSlice {
		if strIdx == "" {
			continue
		}
		idx, err := strconv.Atoi(strIdx)
		if err != nil {
			return []int{}, err
		}

		intSlice[elementCount] = idx
		elementCount += 1
	}
	log.Debugf("sliceAtoi: %v", intSlice)
	return intSlice[:elementCount], nil
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
		usage()
		return
	}
}
