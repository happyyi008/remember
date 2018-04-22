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
	logging "github.com/op/go-logging"
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
