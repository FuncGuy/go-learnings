package main

import (
	"fmt"
	"os"

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

func newHnStories() []Story {
	var stories []Story
	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, id := range changes.Items {

		story, err := hackerNewsClient.GetStory(id)

		if err != nil {
			continue
		}

		newStory := Story{
			title:  story.Title,
			url:    story.URL,
			author: story.By,
			source: "HackerNews",
		}
		stories = append(stories, newStory)
	}
	return stories

}

func main() {
	hnStories := newHnStories()
	hnStories2 := newHnStories()
	var stories []Story
	if hnStories != nil {
		stories = append(stories, hnStories...)
		stories = append(stories, hnStories2...)
	}

	file, err := os.Create("stories.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	for _, s := range stories {
		fmt.Fprintf(file, "%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}

	for _, s := range stories {
		fmt.Printf("%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}
