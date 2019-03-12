package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

var change = make(chan string)
var done = make(chan string)
var watcher *fsnotify.Watcher

func main() {

	fmt.Printf("RgxMon %s\n", os.Args[1])
	//TODO: take regex arg
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	subscribe(os.Args[1])

	go watch()
	go act()

	exitReason := <-done
	fmt.Printf("rsgMon quit because %s", exitReason)
}

func act() {
	for {
		select {
		case c := <-change:
			fmt.Println("Reacting to %s", c)
			//TODO: regex the file.
		case <-done:
			return
		}
	}
}

func subscribe(filename string) {
	err := watcher.Add(filename)
	if err != nil {
		fmt.Printf("Failed to sub: %s", err)
		log.Fatal(err)
	}
}

func pauseTryResubscribe(filename string) bool {

	var err error
	for i := 0; i < 5; i++ {

		time.Sleep(50 * time.Millisecond)

		err = watcher.Add(filename)
		if err != nil {
			return true
		}
	}
	return false
}

func watch() {
	for {
		select {
		case event := <-watcher.Events:
			fmt.Println(event.Name)
			fmt.Println(event.Op)
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf(event.Name)
				change <- event.Name
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				success := pauseTryResubscribe(event.Name)
				if !success {
					done <- "File deleted."
				}
			}
		case err, ok := <-watcher.Errors:
			fmt.Println(err)
			if !ok {
				done <- "Errored."
				return
			}
		}
	}
}
