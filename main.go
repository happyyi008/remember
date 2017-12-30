package main

/*
 *
 * next features:
 * status for todos
 * ability to change the status of a todo mark done
 * show status by colors
 * filter ls by status
 *
 */

import (
	"encoding/json"
	"flag"
	"fmt"
	logging "github.com/op/go-logging"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

const app = "rmb"

var (
	log     = logging.MustGetLogger(app)
	RMBFILE = os.Getenv("HOME") + "/.remember"
)

var Usage = func() {
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Usage of %s:\n", app)
	fmt.Fprintf(w, "$ rmb -help | -h\tprint your list of todos\n")
	fmt.Fprintf(w, "$ rmb ls\tprint your list of todos\n")
	fmt.Fprintf(w, "$ rmb rm <index>\tremoves the todo at <index> from your list\n")
	fmt.Fprintf(w, "$ rmb <todo>\tadds a new todo to your list\n")
	w.Flush()
}

func initializeFile() []byte {
	log.Debug("initialize Remember")
	empty := []byte(`{"todoList":[]}`)
	write(empty)
	log.Debug("created init file")
	return empty
}

func main() {
	help := flag.Bool("help", false, "print usage")
	flag.BoolVar(help, "h", false, "print usage")
	logLevel := flag.String("log-level", "INFO", "set log level")
	flag.Parse()
	if *help {
		Usage()
		return
	}

	// set log level
	level, _ := logging.LogLevel(*logLevel)
	logging.SetLevel(level, app)
	content, err := ioutil.ReadFile(RMBFILE)
	if err != nil {
		content = initializeFile()
	}
	remember := &Remember{}
	json.Unmarshal(content, remember)
	log.Debugf("Done init: %+v", remember)

	cliArgs := flag.Args()

	if len(cliArgs) == 0 {
		remember.addTodo(cliArgs)
		return
	}

	log.Debugf("Arguments: %+v", cliArgs)
	switch cliArgs[0] {
	case "ls":
		remember.listTodo()
	case "rm":
		remember.deleteTodo(cliArgs)
	default:
		remember.addTodo(cliArgs)
	}
	remember.writeToFile()
}

func write(payload []byte) {
	ioutil.WriteFile(RMBFILE, payload, 0644)
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
		return
	}
}
