package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/fsnotify/fsnotify"
)

var args struct {
	Command string   `arg:"-c,--command" help:"Command to execute when file(s) changes."`
	Input   []string `arg:"positional"`
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

func execCommand(cmd string, file string) error {
	executable := "cmd.exe"
	args := []string{"/c"}

	parsed := cmd
	if strings.Contains(cmd, "%s") {
		parsed = fmt.Sprintf(cmd, file)
	}

	args = append(args, parsed)

	command := exec.Command(executable, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		log.Println("error:", err)
		return err
	}
	return nil
}

func main() {
	arg.MustParse(&args)
	fmt.Println("Command:", args.Command)
	fmt.Println("File(s):", args.Input)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					execCommand(args.Command, event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
					execCommand(args.Command, event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("removed file:", event.Name)

					watcher.Remove(event.Name)
					time.Sleep(3 * time.Second)
					if ok, err := exists(event.Name); ok {
						watcher.Add(event.Name)
						execCommand(args.Command, event.Name)
					} else {
						log.Println("error:", err)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, file := range args.Input {
		err = watcher.Add(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done
}
