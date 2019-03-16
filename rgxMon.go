package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/t3hmun/t3hterm"

	"github.com/fsnotify/fsnotify"
	"github.com/t3hmun/t3hstr"
)

var change = make(chan string)
var done = make(chan string)
var watcher *fsnotify.Watcher
var target []byte

func main() {

	fmt.Printf("RgxMon %s %s\n", os.Args[1], os.Args[2])
	//TODO: take regex arg
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	target = ReadFile(as.Args[2])
	subscribe(os.Args[1])

	go watch()
	go act()
	change <- os.Args[1]

	exitReason := <-done
	fmt.Printf("rgxMon quit because %s\n", exitReason)
}

func act() {
	for {
		select {
		case c := <-change:
			rgxText := string(readFile(c))
			runRegex(rgxText)
			//TODO: regex the file.
		case <-done:
			return
		}
	}
}

func readFile(filename string) []byte {
	result, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read %s", filename)
		log.Fatal(err)
	}
	return result
}

func runRegex(regexText string) {
	rgx, err := regexp.Compile(regexText)
	if err == nil {
		matches := rgx.FindAllSubmatch(target, -1)
		termSize, err := t3hterm.GetSizeAndPosition()
		if err != nil {
			panic(err)
		}
		for i, match := range matches {
			for j, group := range match {
				g := t3hstr.Gestalt([]rune(string(group)), termSize.Width, []rune("..."))
				fmt.Printf("%d: %s\n", j, g)
			}
		}
	} else {
		fmt.Printf("Err: %s\n", err)
	}
}

func subscribe(filename string) {
	err := watcher.Add(filename)
	if err != nil {
		fmt.Printf("Failed to sub: %s\n", err)
		log.Fatal(err)
	}
}

func pauseTryResubscribe(filename string) bool {

	var err error
	for i := 0; i < 5; i++ {

		time.Sleep(50 * time.Millisecond)

		err = watcher.Add(filename)
		if err == nil {
			return true
		}
	}
	return false
}

func watch() {
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf(event.Name)
				change <- event.Name
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				success := pauseTryResubscribe(event.Name)
				if !success {
					done <- "File deleted."
				}
				change <- event.Name
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
