package main

import (
	"log"
	"strings"
	"github.com/picadoh/gostreamer/streamer"
	"os"
	"fmt"
	"net"
	"bufio"
)

type TextFileCollector struct {
	streamer.Collector
}

type TextSocketCollector struct {
	streamer.Collector
}

type HashTagExtractor struct {
	streamer.Processor
}

type TextPublisher struct {
	streamer.Processor
}

type HashTagCounter struct {
	streamer.Processor
	State streamer.Counter
}

type HashTagCountPublisher struct {
	streamer.Processor
}

func (collector*TextFileCollector) Execute(name string, cfg streamer.Config, out*chan streamer.Message) {
	lines, _ := streamer.ReadLines(cfg.GetString("source.file"))

	for _, line := range lines {
		out_message := streamer.NewMessage()
		out_message.Put("tweet", line)
		log.Printf("Read message from file: %s\n", out_message)
		*out <- out_message
	}
}

func (collector*TextSocketCollector) Execute(name string, cfg streamer.Config, out*chan streamer.Message) {
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

func (processor*TextPublisher) Execute(name string, cfg streamer.Config, input streamer.Message, out *chan streamer.Message) {
	tweet, _ := input.Get("tweet").(string)

	words := strings.Split(tweet, " ")

	for _, word := range words {
		fmt.Println(word)
	}
}

func (processor*HashTagExtractor) Execute(name string, cfg streamer.Config, input streamer.Message, out *chan streamer.Message) {
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

func (processor*HashTagCounter) Execute(name string, cfg streamer.Config, input streamer.Message, out *chan streamer.Message) {
	hashtag := input.Get("hashtag").(string)

	count := processor.State.Increment(hashtag)

	out_message := streamer.NewMessage()
	out_message.Put("hashtag", hashtag)
	out_message.Put("count", count)

	*out <- out_message
}

func (processor*HashTagCountPublisher) Execute(name string, cfg streamer.Config, input streamer.Message, out *chan streamer.Message) {
	hashtag, _ := input.Get("hashtag").(string)
	count, _ := input.Get("count").(int)

	log.Printf("Publishing %s/%d\n", hashtag, count)
}

func defineTweetSource(cfg streamer.Config) streamer.Collector {
	switch cfg.GetString("source.mode") {
	case "file": return &TextFileCollector{}
	case "socket": return &TextSocketCollector{}
	}
	return nil
}

func RunPipeline(cfg streamer.Config) {
	// read config
	var tweetSource = defineTweetSource(cfg)
	var extractorParallelismHint = cfg.GetInt("parallelism.extractor")
	var counterParallelismHint = cfg.GetInt("parallelism.counter")
	var publisherParallelismHint = cfg.GetInt("parallelism.publisher")

	// define state
	var counterState = streamer.NewCounter()

	// build pipeline
	collector := &streamer.BaseCollector{Delegate:tweetSource}
	extractor := &streamer.BaseProcessor{Delegate:&HashTagExtractor{}, Balancer:streamer.NewRandomDemux(extractorParallelismHint)}
	counter := &streamer.BaseProcessor{Delegate:&HashTagCounter{State:*counterState}, Balancer:streamer.NewGroupDemux(counterParallelismHint, "hashtag")}
	publisher := &streamer.BaseProcessor{Delegate:&HashTagCountPublisher{}, Balancer:streamer.NewGroupDemux(publisherParallelismHint, "hashtag")}

	// execute pipeline
	sequence := collector.Execute("collector", cfg)
	extracted := extractor.Execute("extractor", cfg, sequence)
	counted := counter.Execute("counter", cfg, extracted)
	<-publisher.Execute("publisher", cfg, counted)

	// print final report
	log.Printf("final count report: %s\n", counterState.Count)
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
