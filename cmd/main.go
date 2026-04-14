package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fusu2323/cli-todo/internal/task"
)

// printGlobalHelp prints the global help message when no subcommand or help is given.
func printGlobalHelp() {
	fmt.Println("Usage: todo <command> [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  add <title> [-c category]    Add a new task")
	fmt.Println("  list [-c category]           List all tasks")
	fmt.Println("  done <id>                    Mark a task as complete")
	fmt.Println("  delete <id>                  Delete a task")
	fmt.Println("  help                         Show this help message")
}

// handleAdd handles the "add" subcommand.
func handleAdd(store *task.JSONFileStore) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	category := fs.String("c", "", "category")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: todo add <title> [-c category]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}
	t, err := task.NewTask(fs.Arg(0), *category)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := store.Add(t); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	// Silent success per D-01
}

// handleList handles the "list" subcommand.
func handleList(store *task.JSONFileStore) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	category := fs.String("c", "", "category")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: todo list [-c category]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])
	tasks, err := store.List(*category)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, t := range tasks {
		check := " "
		if t.Completed {
			check = "x"
		}
		if t.Category != "" {
			fmt.Printf("[%s] %s @%s\n", check, t.Title, t.Category)
		} else {
			fmt.Printf("[%s] %s\n", check, t.Title)
		}
	}
}

// handleDone handles the "done" subcommand.
func handleDone(store *task.JSONFileStore) {
	fs := flag.NewFlagSet("done", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: todo done <id>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}
	err := store.MarkDone(fs.Arg(0))
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			fmt.Fprintln(os.Stderr, "task not found")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	// Silent success
}

// handleDelete handles the "delete" subcommand.
func handleDelete(store *task.JSONFileStore) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: todo delete <id>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])
	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}
	err := store.Delete(fs.Arg(0))
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			fmt.Fprintln(os.Stderr, "task not found")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	// Silent success
}

func main() {
	if len(os.Args) < 2 {
		printGlobalHelp()
		os.Exit(0)
	}

	store, err := task.NewJSONFileStore("")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		handleAdd(store)
	case "list":
		handleList(store)
	case "done":
		handleDone(store)
	case "delete":
		handleDelete(store)
	case "help":
		printGlobalHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		printGlobalHelp()
		os.Exit(1)
	}
}
