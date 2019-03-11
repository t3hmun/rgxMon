package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var op = make(chan fsnotify.Event)
var watcher *fsnotify.Watcher

func main() {

	fmt.Printf("RgxMon %s\n", os.Args[1])
	var err error
	watcher, err = fsnotify.NewWatcher()
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
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf(event.Name)
				op <- event
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println(err)
		}
	}
}
