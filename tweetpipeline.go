package main

import (
	"log"
	"strings"
	"github.com/picadoh/gostreamer/streamer"
	"os"
	"fmt"
	"net"
	"bufio"
	"time"
)

func TweetsFileCollector(out *chan streamer.Message) {
	lines, _ := streamer.ReadLines(os.Args[2])

	for _, line := range lines {
		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)
		log.Printf("Read message from file: %s\n", out_message)
		*out <- out_message
	}
}

func TweetsSocketCollector(out *chan streamer.Message) {
	listener, _ := net.Listen("tcp", ":" + os.Args[2])
	conn, _ := listener.Accept()

	for {
		line, _ := bufio.NewReader(conn).ReadString('\n')
		line = strings.TrimSuffix(line, "\n")

		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)

		log.Printf("Received raw message from socket: %s\n", out_message)

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

func RunPipeline(tweetSource streamer.CollectorFunction) {
	sequence := streamer.SCollector("collector", tweetSource)

	extracted := streamer.SProcessor("extractor", streamer.NewRandomDemux(5), sequence, HashTagExtractor)

	counter := streamer.NewCounter()

	go func() {
		// start a routine that periodically prints the report
		for {
			log.Printf("count report: %s\n", counter.Count)
			time.Sleep(10 * time.Second)
		}
	}()

	counted := streamer.SProcessor("counter", streamer.NewGroupDemux(5, "hashtag"), extracted,
		func(input streamer.Message, out*chan streamer.Message) {
			hashtag := input.Get("hashtag").(string)

			count := counter.Increment(hashtag)

			out_message := streamer.NewMessage()
			out_message.Put("hashtag", hashtag)
			out_message.Put("count", count)

			*out <- out_message
		})

	<-streamer.SProcessor("publisher", streamer.NewGroupDemux(5, "hashtag"), counted, HashTagCountPublisher)

	log.Printf("final count report: %s\n", counter.Count)
}

func DefineSourceCollectorFromArgs() streamer.CollectorFunction {
	if (len(os.Args) < 3 || (os.Args[1] != "-f" && os.Args[1] != "-l")) {
		return nil
	}

	if (os.Args[1] == "-f") {
		return TweetsFileCollector
	} else {
		return TweetsSocketCollector
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var tweetSource = DefineSourceCollectorFromArgs()

	if (tweetSource == nil) {
		fmt.Println("Usage: " + os.Args[0] + " [-f <path/to/tweets/file> | -l <port>]")
		os.Exit(2)
	}

	RunPipeline(tweetSource)
}
