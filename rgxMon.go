package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

var op = make(chan fsnotify.Event)
var watcher Watcher

func main() {

	var err Error
	watcher, err = wathcher.NewWatcher
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

}

func watch() {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			fmt.Printf(event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf(event)
				op <- event
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.PrintF(err)
		}
	}
}
