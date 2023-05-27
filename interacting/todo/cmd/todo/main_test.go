package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {

	if len(os.Getenv("TODO_FILENAME")) > 0 {
		fileName = os.Getenv("TODO_FILENAME")
	}

	fmt.Println("Deleting", fileName, "...")
	os.Remove(fileName)

	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")

	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("AddMultipleFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, "test task number 4\ntest task number 5")
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task3 := "test task number 3"
	t.Run("DeleteTask", func(t *testing.T) {
		add := exec.Command(cmdPath, "-add", task3)
		if err := add.Run(); err != nil {
			t.Fatal(err)
		}
		del := exec.Command(cmdPath, "-del", "3")
		if err := del.Run(); err != nil {
			t.Fatal(err)
		}

	})

	task4 := "test task number 4"
	t.Run("DeleteTaskFromSTDIN", func(t *testing.T) {
		add := exec.Command(cmdPath, "-add", task4)
		if err := add.Run(); err != nil {
			t.Fatal(err)
		}
		del := exec.Command(cmdPath, "-del")
		cmdStdIn, err := del.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, "3")
		cmdStdIn.Close()

		if err := del.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("X	1: %s\n	2: %s\n	3: test task number 3\n	4: test task number 4\n", task, task2)

		if expected != string(out) {
			t.Errorf("expected %q, got %q instead\n", expected, string(out))
		}
	})
}
