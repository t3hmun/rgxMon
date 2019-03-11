package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var change = make(chan string)
var done = make(chan bool)
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

	err = watcher.Add(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	go watch()
	go act()

	<-done
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

func watch() {
	for {
		select {
		case event, ok := <-watcher.Events:
			fmt.Println(event.Name)
			fmt.Println(event.Op)
			if !ok {
				done <- true
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf(event.Name)
				change <- event.Name
			}
		case err, ok := <-watcher.Errors:
			fmt.Println(err)
			if !ok {
				done <- true
				return
			}
		}
	}
}
