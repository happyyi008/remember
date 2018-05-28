package main

/*
 *
 * next features:
 * status for todos
 * ability to change the status of a todo mark done
 * show status by colors
 * filter ls by status
 * interactive mode
 *
 */

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	logging "github.com/op/go-logging"
	"os"
	"strings"
	"text/tabwriter"
)

const app = "remember"

var (
	log     = logging.MustGetLogger(app)
	RMBFILE = os.Getenv("HOME") + "/.remember"
)

func usage() {
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s usage:\n", app)
	fmt.Fprintf(w, "\thelp\tprint usage\n")
	fmt.Fprintf(w, "\tls\tprint your list of todos\n")
	fmt.Fprintf(w, "\trm [index ...]\tremoves the todo at <index>\n")
	fmt.Fprintf(w, "\tset <status> [index ...]\tset the status of <index>\n")
	fmt.Fprintf(w, "\tnew <todo>\tadds a new todo to your list\n")
	fmt.Fprintf(w, "\n")
	w.Flush()
}

func init() {
	help := flag.Bool("help", false, "print usage")
	flag.BoolVar(help, "h", false, "print usage")

	logLevel := flag.String("log-level", "INFO", "set log level")
	flag.Parse()

	// check if help is set
	if *help {
		usage()
		return
	}

	level, _ := logging.LogLevel(*logLevel)
	logging.SetLevel(level, app)
}

func main() {
	remember := NewRemember()
	remember.listTodo()

	for {
		action := prompt.Input("> ", completer)
		log.Debugf("action: %+v", action)
		if action == "" {
			break
		}
		dispatch(remember, action)
	}

	return
}

func dispatch(r *Remember, action string) {
	log.Debugf("Action: %+v", action)
	args := strings.Split(action, " ")

	switch args[0] {
	case "ls":
		r.listTodo()
	case "rm":
		if len(args) < 1 {
			fmt.Println("Error: missing delete index")
			break
		}
		r.deleteTodo(args[1:])
	case "set":
		if len(args) < 3 {
			fmt.Println("Error: too little arguments")
			break
		}
		r.setStatus(args[1], args[2:])
	case "new":
		trimmedMessage := strings.TrimLeft(action, "new ")
		r.addTodo(checkMessage(trimmedMessage))
	case "help":
		usage()
	default:
		r.addTodo(checkMessage(action))
	}
}

func checkMessage(userInput string) string {
	if userInput != "" {
		return userInput
	}
	return readStdin()
}

func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n") // trim newline from the reader
	return input
}
