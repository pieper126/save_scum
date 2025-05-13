package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type SaveWatcher struct {
	config           *Config
	done             chan (bool)
	internal_watcher *fsnotify.Watcher
}

type EventType int

const (
	CREATE EventType = iota
	UPDATE
	DELETE
)

type Event struct {
	EventType EventType
	FileName  string
}

func BuildWatcher(config Config) *SaveWatcher {
	folderToWatch := config.ReadFrom
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(folderToWatch); os.IsNotExist(err) {
		fmt.Println("Creating folder:", folderToWatch)
		os.Mkdir(folderToWatch, os.ModePerm)
		fmt.Print("err: ", err)
	}

	err = watcher.Add(folderToWatch)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	return &SaveWatcher{
		config:           &config,
		done:             done,
		internal_watcher: watcher,
	}
}

func (saveWatcher *SaveWatcher) Watch() <-chan Event {
	res := make(chan Event, 100)
	go func() {
		for {
			select {
			case event, ok := <-saveWatcher.internal_watcher.Events:
				if !ok {
					return
				}
				fmt.Printf("Event: %s on %s\n", event.Op, event.Name)
				parsed := saveWatcher.mapEvent(event)
				if parsed != nil {
					res <- *parsed
				}
			case err, ok := <-saveWatcher.internal_watcher.Errors:
				if !ok {
					return
				}
				res <- Event{}
				fmt.Println("Error:", err)
			case <-saveWatcher.done:
				fmt.Println("Exiting!")
				close(res)
				return
			}
		}
	}()
	return res
}

func (saveWatcher *SaveWatcher) Done() {
	saveWatcher.done <- true
}

func (saveWatcher *SaveWatcher) mapEvent(event fsnotify.Event) *Event {
	switch event.Op {
	case fsnotify.Create:
		return &Event{EventType: CREATE, FileName: event.Name}
	case fsnotify.Rename:
		return &Event{EventType: UPDATE, FileName: event.Name}
	case fsnotify.Write:
		return &Event{EventType: UPDATE, FileName: event.Name}
	default:
		return nil
	}
}
