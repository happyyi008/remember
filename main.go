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
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	logging "github.com/op/go-logging"
	"os"
	"strings"
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
	fmt.Fprintf(w, "$ rmb rm <index>...\tremoves the todo at <index> from your list\n")
	fmt.Fprintf(w, "$ rmb <todo>\tadds a new todo to your list\n")
	w.Flush()
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "ls", Description: "Print your todo list"},
		{Text: "rm", Description: "Delete from todo list"},
		{Text: "set", Description: "Set status of Todo"},
		{Text: "help", Description: "Get help"},
	}

	currentLine := strings.Split(d.CurrentLineBeforeCursor(), " ")
	log.Debugf("current line %+v, len: %d, 0: %v", currentLine, len(currentLine), currentLine[0])

	// get suggestion for set instead
	if currentLine[0] == "set" && len(currentLine) == 2 {
		s = []prompt.Suggest{
			{Text: "start", Description: "Start the task"},
			{Text: "done", Description: "Finish task"},
			{Text: "restart", Description: "Reset status to new"},
		}
	} else if len(currentLine) == 2 {
		s = []prompt.Suggest{}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func dispatch(args []string, r *Remember) {
	log.Debugf("Arguments: %+v", args)
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
		r.setStatus(args[2:], args[1])
	default:
		r.addTodo(args)
	}

}

func main() {
	help := flag.Bool("help", false, "print usage")
	flag.BoolVar(help, "h", false, "print usage")

	logLevel := flag.String("log-level", "INFO", "set log level")
	flag.Parse()

	// check if help is set
	if *help {
		Usage()
		return
	}

	level, _ := logging.LogLevel(*logLevel)
	logging.SetLevel(level, app)

	cliArgs := flag.Args()

	remember := NewRemember()

	if len(cliArgs) == 0 { // run interactive mode
		remember.listTodo()
		for {
			action := prompt.Input("> ", completer)
			log.Debugf("action: %+v", action)
			if action == "" {
				break
			}
			args := strings.Split(action, " ")
			dispatch(args, remember)
		}
		return
	}
}
