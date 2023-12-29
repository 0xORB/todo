package commands

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"todo/app"
)

/*
Commands to Implement:
- Add
- Remove
- Load (file)
- List
- Store (file)
- Mark Complete

Parse Args
*/

const filename = ".tasks.json"

func App() {
	add := flag.Bool(
		"add",
		false,
		"Add a new task.",
	)

	del := flag.Int(
		"del",
		0,
		"Remove an existing task via taskID.",
	)

	complete := flag.Int(
		"done",
		0,
		"Mark an existing task as complete via taskID.",
	)

	list := flag.Bool(
		"ls",
		false,
		"List all tasks.",
	)

	flag.Parse()

	tasks := &app.Tasks{}

	if err := tasks.LoadTasks(filename); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:

		t, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		tasks.CreateTask(t)
		err = tasks.StoreTasks(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		tasks.PrintTasks()

	case *del > 0:

		err := tasks.RemoveTask(*del)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = tasks.StoreTasks(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		tasks.PrintTasks()

	case *complete > 0:

		err := tasks.MarkTaskComplete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = tasks.StoreTasks(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		tasks.PrintTasks()

	case *list:
		tasks.PrintTasks()

	default:
		fmt.Fprintln(os.Stdout, "Invalid command.")
		os.Exit(0)

	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()
	if len(text) == 0 {
		return "", errors.New("No value specified.")
	}

	return text, nil
}
