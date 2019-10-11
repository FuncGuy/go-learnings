package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

func init() {
	hackerNewsClient = gophernews.NewClient()

}

type Story struct {
	title  string
	url    string
	author string
	source string
}

func newHnStories(c chan<- Story) {
	defer close(c)

	changes, err := hackerNewsClient.GetChanges()
	if err != nil {

	}

	var wg sync.WaitGroup

	for _, id := range changes.Items {
		wg.Add(1)
		go getHnStroyDetails(id, c, &wg)

	}
	wg.Wait()

}

func getHnStroyDetails(id int, c chan<- Story, wg *sync.WaitGroup) {
	defer wg.Done()
	story, err := hackerNewsClient.GetStory(id)

	if err != nil {

	}

	newStory := Story{
		title:  story.Title,
		url:    story.URL,
		author: story.By,
		source: "HackerNews",
	}

	c <- newStory

}

func outputToConsole(c <-chan Story) {
	for {
		s := <-c
		fmt.Printf("%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

func outputToFile(c <-chan Story, file *os.File) {
	for {
		s := <-c
		fmt.Fprintf(file, "%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

func main() {

	fromHn := make(chan Story, 8) // input channel

	toConsole := make(chan Story, 8)

	toFile := make(chan Story, 8)

	go newHnStories(fromHn)

	go newHnStories(fromHn)

	file, err := os.Create("stories.txt")

	if err != nil {

	}

	go outputToConsole(toConsole)
	go outputToFile(toFile, file)

	hnOpen := true

	for hnOpen {
		select {
		case story, open := <-fromHn:
			if open {
				toFile <- story
				toConsole <- story
			} else {
				hnOpen = false
			}
		}
	}

}
