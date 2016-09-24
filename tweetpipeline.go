package main

import (
	"log"
	"strings"
	"github.com/picadoh/gostreamer/streamer"
	"os"
	"fmt"
)

func TweetsFileCollector(out *chan streamer.Message) {
	lines, _ := streamer.ReadLines(os.Args[1])

	rand.Intn(50);
	for _, line := range lines {
		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)
		log.Printf("Generated message %s\n", out_message)
		*out <- out_message
	}
}

func HashTagExtractor(input streamer.Message, out *chan streamer.Message) {
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

func HashTagCountPublisher(input streamer.Message, out *chan streamer.Message) {
	hashtag, _ := input.Get("hashtag").(string)
	count, _ := input.Get("count").(int)

	log.Printf("Publishing %s/%d\n", hashtag, count)
}

func RunPipeline() {
	sequence := streamer.SCollector("collector", TweetsFileCollector)

	extracted := streamer.SProcessor("extractor", streamer.NewRandomDemux(5), sequence, HashTagExtractor)

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

	<-streamer.SProcessor("publisher", streamer.NewGroupDemux(5, "hashtag"), counted, HashTagCountPublisher)

	log.Printf("report: %s\n", counter.Count)
}

func ValidateArguments() {
	if (len(os.Args) < 2) {
		fmt.Println("Usage: " + os.Args[0] + " <path/to/tweets/file>")
		os.Exit(2)
	}
}

func main() {
	ValidateArguments()
	RunPipeline()
}
