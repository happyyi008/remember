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

const app = "remember"

var (
	log     = logging.MustGetLogger(app)
	RMBFILE = os.Getenv("HOME") + "/.remember"
)

var Usage = func() {
	w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s usage:\n", app)
	fmt.Fprintf(w, "\thelp\tprint usage\n")
	fmt.Fprintf(w, "\tls\tprint your list of todos\n")
	fmt.Fprintf(w, "\trm [index...]\tremoves the todo at <index>\n")
	fmt.Fprintf(w, "\tset <status> [index...]\tset the status of <index>\n")
	fmt.Fprintf(w, "\tnew <todo>\tadds a new todo to your list\n")
	fmt.Fprintf(w, "\n")
	w.Flush()
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
	case "new":
		r.addTodo(args[1:])
	case "help":
		Usage()
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
