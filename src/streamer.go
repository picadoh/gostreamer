package main

import (
	"log"
	"strings"
)

func collector(out *chan Message) {
	lines, _ := readLines("tweets.txt")

	for _, line := range lines {
		out_message := NewMessage()
		out_message.Put("tweet", line)
		log.Printf("Generated message %s\n", out_message)
		*out <- out_message
	}
}

func extractor(input Message, out *chan Message) {
	tweet, _ := input.Get("tweet").(string)

	words := strings.Split(tweet, " ")

	for _, word := range words {
		if (strings.HasPrefix(word, "#")) {
			out_message := NewMessage()
			out_message.Put("hashtag", word)
			log.Printf("Extracted hashtag %s\n", word)
			*out <- out_message
		}
	}
}

func publisher(input Message, out *chan Message) {
	hashtag, _ := input.Get("hashtag").(string)
	count, _ := input.Get("count").(int)

	log.Printf("Publishing %s/%d\n", hashtag, count)
}

func pipeline() {
	scollector("collector", collector)

	sequence := scollector("collector", collector)

	extracted := sprocessor("extractor", NewRandomDemux(5), sequence, extractor)

	counter := NewCounter()
	counted := sprocessor("counter", NewGroupDemux(5, "hashtag"), extracted, func(input Message, out*chan Message) {
		hashtag := input.Get("hashtag").(string)

		count := counter.Increment(hashtag)

		out_message := NewMessage()
		out_message.Put("hashtag", hashtag)
		out_message.Put("count", count)

		log.Printf("Counted %s/%d\n", hashtag, count)

		*out <- out_message
	})

	<-sprocessor("publisher", NewGroupDemux(5, "hashtag"), counted, publisher)

	log.Printf("report: %s\n", counter.count)
}

func main() {
	pipeline()
}