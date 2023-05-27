package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"pragprog.com/rggo/interacting/todo"
)

var todoFileName = ".todo.json"

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "'-add' and '-del' accept input from STDIN")
	}

	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all the tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Bool("del", false, "Delete task from the ToDo list")
	verbose := flag.Bool("verbose", false, "List all the tasks with additional information")
	pending := flag.Bool("pending", false, "List only tasks that are not Done")

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {

	case *pending:
		formatted := l.Pending()
		fmt.Println(formatted)

	case *verbose:
		formatted := ""
		for _, t := range *l {
			prefix := "	"
			if t.Done {
				prefix = "X	"
			}

			formatted += fmt.Sprintf("%s[%s] %s\n", prefix, t.CreatedAt, t.Task)
		}

		fmt.Println(formatted)

	case *list:
		fmt.Print(l)

	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		tasks, err := getTasks(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, t := range tasks {
			l.Add(t)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *delete:
		i, err := getI(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.Delete(i)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTasks(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}

	tasks := []string{}
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return []string{}, err
		}

		if len(s.Text()) == 0 {
			return []string{}, fmt.Errorf("task cannot be blank")
		}

		tasks = append(tasks, s.Text())
	}

	return tasks, nil
}

func getI(r io.Reader, args ...string) (int, error) {
	if len(args) > 1 {
		return 0, fmt.Errorf("too many arguments")
	}

	if len(args) == 1 {
		parsed, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return 0, err
		}

		return int(parsed), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return 0, err
	}

	if len(s.Text()) == 0 {
		return 0, fmt.Errorf("index is required")
	}

	parsed, err := strconv.ParseInt(s.Text(), 10, 32)
	if err != nil {
		return 0, err
	}

	return int(parsed), nil

}
