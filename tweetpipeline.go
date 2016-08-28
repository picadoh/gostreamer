package main

import (
	"log"
	"strings"
	"github.com/picadoh/gostreamer/streamer"
)

func collector(out *chan streamer.Message) {
	lines, _ := streamer.ReadLines("tweets.txt")

	for _, line := range lines {
		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)
		log.Printf("Generated message %s\n", out_message)
		*out <- out_message
	}
}

func extractor(input streamer.Message, out *chan streamer.Message) {
	tweet, _ := input.Get("tweet").(string)

	words := strings.Split(tweet, " ")

	for _, word := range words {
		if (strings.HasPrefix(word, "#")) {
			out_message := streamer.NewMessage()
			out_message.Put("hashtag", word)
			log.Printf("Extracted hashtag %s\n", word)
			*out <- out_message
		}
	}
}

func publisher(input streamer.Message, out *chan streamer.Message) {
	hashtag, _ := input.Get("hashtag").(string)
	count, _ := input.Get("count").(int)

	log.Printf("Publishing %s/%d\n", hashtag, count)
}

func pipeline() {
	sequence := streamer.SCollector("collector", collector)

	extracted := streamer.SProcessor("extractor", streamer.NewRandomDemux(5), sequence, extractor)

	counter := streamer.NewCounter()
	counted := streamer.SProcessor("counter", streamer.NewGroupDemux(5, "hashtag"), extracted,
		func(input streamer.Message, out*chan streamer.Message) {
			hashtag := input.Get("hashtag").(string)

			count := counter.Increment(hashtag)

			out_message := streamer.NewMessage()
			out_message.Put("hashtag", hashtag)
			out_message.Put("count", count)

			log.Printf("Counted %s/%d\n", hashtag, count)

			*out <- out_message
		})

	<-streamer.SProcessor("publisher", streamer.NewGroupDemux(5, "hashtag"), counted, publisher)

	log.Printf("report: %s\n", counter.Count)
}

func main() {
	pipeline()
}
