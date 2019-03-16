package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/t3hmun/t3hstr"
	"github.com/t3hmun/t3hterm"
)

var change = make(chan string)
var done = make(chan string)
var watcher *fsnotify.Watcher
var target []byte
var targetFilename string

func main() {

	fmt.Printf("RgxMon %s %s\n", os.Args[1], os.Args[2])
	//TODO: take regex arg
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	targetFilename = os.Args[2]
	target = readFile(targetFilename)
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
			if rgxText == "" {
				// Some editor behaviour causes the file to be read blank before being poperly read.
				// Also actual blank reges are pointless.
				fmt.Printf("Blank regex, skipped.\n")
			} else {
				// passing the target because i want to extend this to multiple files later.
				runRegex(rgxText, targetFilename, target)
			}
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

func runRegex(regexText string, filename string, fileText []byte) {
	rgx, err := regexp.Compile(regexText)
	if err == nil {
		matches := rgx.FindAllSubmatch(fileText, -1)
		width := getTerminalWidth()
		fmt.Printf("\nFile=%s, Regex=%s\n", filename, regexText)
		for i, match := range matches {
			fmt.Printf("Match %d\n", i)
			for j, group := range match {
				line := fmt.Sprintf("(%d) %s", j, group)
				g := t3hstr.Gestalt([]rune(line), width, []rune("..."))
				fmt.Printf("%s\n", g)
			}
		}
	} else {
		fmt.Printf("Err: %s\n", err)
	}
}

func getTerminalWidth() int {
	termSize, err := t3hterm.GetSizeAndPosition()
	if err != nil {
		return 80
	} else {
		return termSize.Width
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
