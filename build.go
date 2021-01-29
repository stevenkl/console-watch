package main

import (
	"os"

	"github.com/pellared/taskflow"
)

func main() {
	flow := taskflow.New()
	flow.MustRegister(taskBuild())

	flow.Main()
}

// taskBuild returns a new `flowtask.Task` for the primary build step
func taskBuild() taskflow.Task {

	// create `./build` folder if its not exist
	if ok, _ := exists("./build"); !ok {
		os.Mkdir("./build", 0666)
	}

	// Creating new `flowtask.Task`, it will execute the `go build ...` command
	return taskflow.Task{
		Name:        "build",
		Description: "Build the project",
		Command:     taskflow.Exec("go", "build", "-o", "./build/console-watch.exe", "./cmd/console-watch"),
	}
}

// exists checks if the given path (file or folder) already exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
