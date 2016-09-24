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

func TweetsFileCollector(cfg streamer.Config, out *chan streamer.Message) {
	lines, _ := streamer.ReadLines(cfg.GetString("source.file"))

	for _, line := range lines {
		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)
		log.Printf("Read message from file: %s\n", out_message)
		*out <- out_message
	}
}

func TweetsSocketCollector(cfg streamer.Config, out *chan streamer.Message) {
	listener, _ := net.Listen("tcp", ":" + cfg.GetString("source.port"))
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

func HashTagExtractor(cfg streamer.Config, input streamer.Message, out *chan streamer.Message) {
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

func HashTagCountPublisher(cfg streamer.Config, input streamer.Message, out *chan streamer.Message) {
	hashtag, _ := input.Get("hashtag").(string)
	count, _ := input.Get("count").(int)

	log.Printf("Publishing %s/%d\n", hashtag, count)
}

func RunPipeline(cfg streamer.Config) {
	var tweetSource = defineTweetSource(cfg)
	var extractorParallelismHint = cfg.GetInt("parallelism.extractor")
	var counterParallelismHint = cfg.GetInt("parallelism.counter")
	var publisherParallelismHint = cfg.GetInt("parallelism.publisher")

	sequence := streamer.SCollector(cfg, "collector", tweetSource)

	extracted := streamer.SProcessor(cfg, "extractor", streamer.NewRandomDemux(extractorParallelismHint), sequence, HashTagExtractor)

	counter := streamer.NewCounter()

	go func() {
		// start a routine that periodically prints the report
		for {
			log.Printf("count report: %s\n", counter.Count)
			time.Sleep(10 * time.Second)
		}
	}()

	counted := streamer.SProcessor(cfg, "counter", streamer.NewGroupDemux(counterParallelismHint, "hashtag"), extracted,
		func(cfg streamer.Config, input streamer.Message, out*chan streamer.Message) {
			hashtag := input.Get("hashtag").(string)

			count := counter.Increment(hashtag)

			out_message := streamer.NewMessage()
			out_message.Put("hashtag", hashtag)
			out_message.Put("count", count)

			*out <- out_message
		})

	<-streamer.SProcessor(cfg, "publisher", streamer.NewGroupDemux(publisherParallelismHint, "hashtag"), counted, HashTagCountPublisher)

	log.Printf("final count report: %s\n", counter.Count)
}

func defineTweetSource(cfg streamer.Config) streamer.CollectorFunction {
	switch cfg.GetString("source.mode") {
	case "file": return TweetsFileCollector
	case "socket": return TweetsSocketCollector
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	if (len(os.Args) < 2) {
		fmt.Println("Usage: " + os.Args[0] + " [<path/to/config/file>]")
		os.Exit(2)
	}

	var cfg, err = streamer.LoadProperties(os.Args[1])

	if (err != nil) {
		fmt.Printf("An error ocurred reading the properties file %s [cause: %s]\n", os.Args[1], err)
		os.Exit(2)
	}

	log.Printf("Loaded configuration: %s\n", cfg.ToString())

	RunPipeline(cfg)
}
